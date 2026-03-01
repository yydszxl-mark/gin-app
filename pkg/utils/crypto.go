package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

// GenerateSalt 生成指定长度的盐值
func GenerateSalt(length int) (string, error) {
	if length <= 0 {
		length = 16
	}
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword 使用盐值对密码进行哈希
func HashPassword(password, salt string) string {
	hash := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码是否正确
func VerifyPassword(password, salt, hashedPassword string) bool {
	return HashPassword(password, salt) == hashedPassword
}

// ValidatePassword 验证密码（别名，为了兼容性）
func ValidatePassword(password, salt, hashedPassword string) bool {
	return VerifyPassword(password, salt, hashedPassword)
}
