# Swagger 快速开始 🚀

## 一键生成文档

```bash
make swagger
```

## 访问文档

启动服务后访问：
```
http://localhost:9060/swagger/index.html
```

## 项目已完成配置 ✅

### 1. 依赖已安装
- ✅ `github.com/swaggo/swag`
- ✅ `github.com/swaggo/gin-swagger`
- ✅ `github.com/swaggo/files`

### 2. 主入口已配置 (cmd/server/main.go)
```go
import _ "gin-app-start/docs"

//	@title			Gin App API
//	@version		1.0
//	@description	This is a RESTful API server built with Gin framework.
//	@host			localhost:9060
//	@BasePath		/
//	@schemes		http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
```

### 3. 路由已配置 (internal/router/router.go)
```go
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

### 4. Controller 注释已添加
- ✅ AuthController (登录、刷新Token、用户信息等)
- ✅ UserController (用户CRUD)
- ✅ DeviceController (设备CRUD)
- ✅ HealthController (健康检查)

## 使用 JWT 认证测试 API

### 步骤 1: 登录获取 Token
1. 访问 Swagger UI
2. 找到 `POST /api/v1/auth/login`
3. 点击 "Try it out"
4. 输入用户名和密码
5. 点击 "Execute"
6. 复制返回的 `token`

### 步骤 2: 配置认证
1. 点击页面右上角的 🔒 "Authorize" 按钮
2. 在弹出框中输入：`Bearer <your_token>`
   - 注意：`Bearer` 和 token 之间有一个空格
3. 点击 "Authorize"
4. 点击 "Close"

### 步骤 3: 测试需要认证的 API
现在可以测试所有需要认证的接口了，例如：
- `GET /api/v1/user/info` - 获取当前用户信息
- `GET /api/v1/users` - 获取用户列表
- `POST /api/v1/users` - 创建用户

## 添加新接口的 Swagger 注释

### 模板

```go
// MethodName godoc
//
//	@Summary		简短描述（一句话）
//	@Description	详细描述
//	@Tags			标签名（用于分组）
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth（如果需要认证）
//	@Param			参数名	位置	类型	必需	"说明"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/api/v1/path [method]
func (ctrl *Controller) MethodName(c *gin.Context) {
    // 实现代码
}
```

### 参数位置说明

| 位置 | 说明 | 示例 |
|------|------|------|
| `path` | 路径参数 | `@Param id path int true "User ID"` |
| `query` | 查询参数 | `@Param page query int false "Page number"` |
| `body` | 请求体 | `@Param request body dto.CreateUserRequest true "User info"` |
| `header` | 请求头 | `@Param Authorization header string true "Bearer token"` |

### 常用示例

#### 1. 创建资源 (POST)
```go
//	@Summary		Create user
//	@Tags			users
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateUserRequest	true	"User info"
//	@Success		200		{object}	response.Response{data=model.User}
//	@Router			/api/v1/users [post]
```

#### 2. 获取单个资源 (GET with path param)
```go
//	@Summary		Get user by ID
//	@Tags			users
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response{data=model.User}
//	@Router			/api/v1/users/{id} [get]
```

#### 3. 获取列表 (GET with query params)
```go
//	@Summary		List users
//	@Tags			users
//	@Security		BearerAuth
//	@Param			page		query	int	false	"Page number"	default(1)
//	@Param			page_size	query	int	false	"Page size"		default(10)
//	@Success		200			{object}	response.Response{data=response.PageResponse}
//	@Router			/api/v1/users [get]
```

#### 4. 更新资源 (PUT)
```go
//	@Summary		Update user
//	@Tags			users
//	@Security		BearerAuth
//	@Param			id		path	int						true	"User ID"
//	@Param			request	body	dto.UpdateUserRequest	true	"User info"
//	@Success		200		{object}	response.Response{data=model.User}
//	@Router			/api/v1/users/{id} [put]
```

#### 5. 删除资源 (DELETE)
```go
//	@Summary		Delete user
//	@Tags			users
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Router			/api/v1/users/{id} [delete]
```

## 常用命令

```bash
# 生成文档
make swagger

# 格式化注释
make swagger-fmt

# 安装工具
make install-tools

# 运行服务
make run
```

## 文档更新流程

1. 修改或添加 Controller 方法
2. 添加/更新 Swagger 注释
3. 运行 `make swagger` 重新生成文档
4. 启动服务验证文档
5. 在 Swagger UI 中测试 API

## 常见问题

### Q: 文档没有更新？
A: 删除 docs 目录后重新生成
```bash
rm -rf docs
make swagger
```

### Q: 提示找不到类型？
A: 确保使用了正确的包路径，例如：
- ✅ `dto.CreateUserRequest`
- ✅ `model.User`
- ✅ `response.Response`

### Q: 如何测试需要认证的接口？
A: 
1. 先调用登录接口获取 token
2. 点击右上角 "Authorize" 按钮
3. 输入 `Bearer <token>`
4. 测试其他接口

### Q: 如何隐藏某些接口？
A: 不添加 Swagger 注释即可，swag 只会生成有注释的接口

## 更多信息

详细文档请查看：[SWAGGER_GUIDE.md](./SWAGGER_GUIDE.md)

## 项目结构

```
gin-app-start/
├── cmd/server/main.go          # 主入口（包含全局 API 信息）
├── internal/
│   ├── controller/             # Controller 层（添加 Swagger 注释）
│   │   ├── auth_controller.go
│   │   ├── user_controller.go
│   │   └── ...
│   ├── dto/                    # DTO 定义（Swagger 会自动解析）
│   └── router/router.go        # 路由配置（包含 Swagger 路由）
├── docs/                       # 自动生成的文档（不要手动修改）
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── Makefile                    # 包含 swagger 命令
```

## 🎉 开始使用

```bash
# 1. 生成文档
make swagger

# 2. 启动服务
make run

# 3. 访问 Swagger UI
# 浏览器打开: http://localhost:9060/swagger/index.html
```
