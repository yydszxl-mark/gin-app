# Swagger 文档生成指南

## 📚 目录
- [快速开始](#快速开始)
- [Swagger 注释规范](#swagger-注释规范)
- [常用命令](#常用命令)
- [注释示例](#注释示例)
- [常见问题](#常见问题)

## 🚀 快速开始

### 1. 安装 swag 工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

或使用 Makefile：

```bash
make install-tools
```

### 2. 生成文档

```bash
make swagger
```

或直接使用命令：

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

### 3. 访问文档

启动服务后访问：

```
http://localhost:9060/swagger/index.html
```

## 📝 Swagger 注释规范

### 全局 API 信息（main.go）

```go
//	@title			Gin App API
//	@version		1.0
//	@description	This is a RESTful API server built with Gin framework.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:9060
//	@BasePath	/

//	@schemes					http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
```

### Controller 方法注释

#### 基本格式

```go
// MethodName godoc
//
//	@Summary		简短描述
//	@Description	详细描述
//	@Tags			标签名
//	@Accept			json
//	@Produce		json
//	@Param			参数名	参数位置	参数类型	是否必需	"参数说明"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Router			/api/v1/path [method]
func (ctrl *Controller) MethodName(c *gin.Context) {
    // ...
}
```

#### 需要认证的接口

添加 `@Security BearerAuth`：

```go
// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get user information by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Router			/api/v1/users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
    // ...
}
```

## 📖 注释示例

### 1. POST 请求（创建资源）

```go
// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with username, email, phone and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateUserRequest	true	"User information"
//	@Success		200		{object}	response.Response{data=model.User}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
    // ...
}
```

### 2. GET 请求（路径参数）

```go
// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get user information by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response{data=model.User}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/api/v1/users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
    // ...
}
```

### 3. GET 请求（查询参数）

```go
// ListUsers godoc
//
//	@Summary		List users
//	@Description	Get paginated list of users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Page size"			default(10)
//	@Param			username	query		string	false	"Filter by username"
//	@Success		200			{object}	response.Response{data=response.PageResponse}
//	@Failure		401			{object}	response.Response
//	@Failure		500			{object}	response.Response
//	@Router			/api/v1/users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
    // ...
}
```

### 4. PUT 请求（更新资源）

```go
// UpdateUser godoc
//
//	@Summary		Update user information
//	@Description	Update user information by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"User ID"
//	@Param			request	body		dto.UpdateUserRequest	true	"User information to update"
//	@Success		200		{object}	response.Response{data=model.User}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/api/v1/users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
    // ...
}
```

### 5. DELETE 请求

```go
// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete user by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/api/v1/users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
    // ...
}
```

### 6. 登录接口（无需认证）

```go
// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
    // ...
}
```

## 🔧 常用命令

### 生成文档

```bash
# 使用 Makefile
make swagger

# 直接使用 swag
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

### 格式化注释

```bash
# 使用 Makefile
make swagger-fmt

# 直接使用 swag
swag fmt
```

### 验证注释

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal --parseVendor
```

## 📋 参数位置说明

| 位置 | 说明 | 示例 |
|------|------|------|
| `path` | 路径参数 | `/users/{id}` 中的 `id` |
| `query` | 查询参数 | `/users?page=1` 中的 `page` |
| `body` | 请求体 | JSON 请求体 |
| `header` | 请求头 | `Authorization` |
| `formData` | 表单数据 | 文件上传 |

## 🎯 响应类型说明

### 基本响应

```go
@Success 200 {object} response.Response
```

### 带数据类型的响应

```go
@Success 200 {object} response.Response{data=model.User}
```

### 数组响应

```go
@Success 200 {object} response.Response{data=[]model.User}
```

### 分页响应

```go
@Success 200 {object} response.Response{data=response.PageResponse}
```

## ❓ 常见问题

### 1. 生成的文档没有更新？

删除 `docs` 目录后重新生成：

```bash
rm -rf docs
make swagger
```

### 2. 找不到 internal 包的类型？

确保使用了 `--parseInternal` 参数：

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

### 3. JWT 认证如何配置？

在 Swagger UI 中：
1. 点击右上角 "Authorize" 按钮
2. 输入：`Bearer <your_token>`
3. 点击 "Authorize"

### 4. 如何隐藏某些接口？

在方法注释中不添加 Swagger 注释即可。

### 5. 如何自定义响应示例？

在 DTO 结构体中使用 `example` 标签：

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required" example:"john_doe"`
    Email    string `json:"email" binding:"required" example:"john@example.com"`
    Password string `json:"password" binding:"required" example:"password123"`
}
```

## 📚 更多资源

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 规范](https://swagger.io/specification/)
- [OpenAPI 3.0](https://spec.openapis.org/oas/v3.0.0)

## 🔄 工作流程

1. 编写 Controller 方法
2. 添加 Swagger 注释
3. 运行 `make swagger` 生成文档
4. 启动服务查看文档
5. 在 Swagger UI 中测试 API

## ✅ 最佳实践

1. **保持注释更新**：每次修改 API 时同步更新注释
2. **使用有意义的标签**：按功能模块分组 API
3. **详细的描述**：提供清晰的 Summary 和 Description
4. **完整的错误码**：列出所有可能的错误响应
5. **示例数据**：在 DTO 中提供 example 标签
6. **版本控制**：提交生成的文档到版本控制系统
