package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 自定义 JWT Claims
type MyClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleIDs  []uint `json:"role_ids"`
	jwt.RegisteredClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// InitJWT 初始化 JWT 密钥对
func InitJWT(privateKeyPath, publicKeyPath string) error {
	// 如果密钥文件不存在，则生成新的密钥对
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		if err := GenerateKeyPair(privateKeyPath, publicKeyPath); err != nil {
			return err
		}
	}

	// 加载私钥
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return errors.New("failed to decode private key")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// 加载公钥
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	block, _ = pem.Decode(publicKeyData)
	if block == nil {
		return errors.New("failed to decode public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	var ok bool
	publicKey, ok = pubKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	return nil
}

// GenerateKeyPair 生成 RSA 密钥对并保存到文件
func GenerateKeyPair(privateKeyPath, publicKeyPath string) error {
	// 生成 2048 位的 RSA 密钥对
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 保存私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyBlock); err != nil {
		return err
	}

	// 保存公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, publicKeyBlock); err != nil {
		return err
	}

	return nil
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string, roleIDs []uint, expireDuration time.Duration) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key not initialized")
	}

	claims := MyClaims{
		UserID:   userID,
		Username: username,
		RoleIDs:  roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*MyClaims, error) {
	if publicKey == nil {
		return nil, errors.New("public key not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string, expireDuration time.Duration) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的 Token
	return GenerateToken(claims.UserID, claims.Username, claims.RoleIDs, expireDuration)
}
