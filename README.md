# Gin App Start

<div align="center">

🚀 **企业级 Go Web 应用脚手架**

基于 Go 1.24 + Gin 框架，采用清晰的三层架构设计，内置 RBAC 权限系统，开箱即用

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-1.25.12-00ADD8?style=flat)](https://gorm.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[快速开始](#-快速开始) • [功能特性](#-功能特性) • [项目结构](#-项目结构) • [文档](#-文档导航) • [部署](#-部署)

</div>

---

## ✨ 功能特性

### 🏗️ 架构设计
- **标准项目布局** - 遵循 Go 官方推荐的项目结构（cmd/internal/pkg）
- **三层架构** - Controller → Service → Repository 职责清晰，易于维护
- **依赖注入** - 手动依赖注入，代码简洁，无需额外框架
- **泛型仓储** - 基于 Go 1.18+ 泛型的通用数据访问层，减少重复代码

### 🔐 安全与认证
- **RBAC 权限系统** - 完整的用户-角色-权限-菜单四级权限控制
- **JWT 认证** - 基于 RSA 非对称加密的 Token 认证机制
- **权限中间件** - 灵活的路由级权限控制
- **密码加密** - 使用 bcrypt 安全存储用户密码

### 🗄️ 数据存储
- **PostgreSQL 17** - 主数据库，支持自动迁移和连接池管理
- **Redis 7** - 缓存和会话存储，支持连接池配置
- **GORM ORM** - 强大的 ORM 框架，支持关联查询和事务

### 📝 日志与监控
- **结构化日志** - 基于 uber/zap 的高性能日志系统
- **请求日志** - 自动记录所有 HTTP 请求的详细信息
- **错误追踪** - 完整的错误堆栈和上下文信息

### 🔧 开发体验
- **统一错误处理** - 自定义业务错误码，错误信息标准化
- **统一响应格式** - 标准化的 JSON API 响应结构
- **Swagger 文档** - 自动生成交互式 API 文档
- **热重载开发** - 支持 Air 热重载，提升开发效率
- **多环境配置** - 基于 Viper 的灵活配置管理（local/dev/prod）

### 🔌 中间件生态
- **日志中间件** - 记录请求响应详情
- **恢复中间件** - 捕获 panic 并优雅处理
- **CORS 中间件** - 跨域资源共享配置
- **限流中间件** - 基于令牌桶的 API 限流
- **认证中间件** - JWT Token 验证
- **权限中间件** - 基于 RBAC 的权限校验

### 🐳 部署运维
- **Docker 支持** - 多阶段构建，镜像体积小
- **Docker Compose** - 一键启动完整开发环境
- **优雅关闭** - 支持 Graceful Shutdown，确保请求完整处理
- **健康检查** - 内置健康检查接口

## 📚 文档导航

| 文档 | 说明 |
|------|------|
| [项目使用指南](docs/PROJECT_GUIDE.md) | 详细的开发文档和最佳实践 |
| [API 接口文档](docs/API_REFERENCE.md) | 完整的 API 参考和示例 |
| [架构设计文档](docs/ARCHITECTURE.md) | 技术架构和设计理念 |
| [RBAC 权限系统](docs/RBAC.md) | RBAC 权限系统完整文档 |
| [RBAC 快速开始](docs/RBAC_QUICKSTART.md) | RBAC 系统快速上手指南 |

## 🚀 快速开始

### 📋 前置要求

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.24+ | [下载安装](https://go.dev/dl/) |
| PostgreSQL | 12+ (推荐 17) | 主数据库 |
| Redis | 6.0+ (推荐 7) | 缓存服务 |
| Docker | 20.10+ | 可选，用于容器化部署 |
| Make | - | 可选，用于快捷命令 |

### 🎯 方式一：Docker Compose（推荐新手）

一键启动完整开发环境，包含应用、PostgreSQL 和 Redis：

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd gin-app-start

# 2. 启动所有服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f app

# 4. 初始化 RBAC 数据（可选）
docker-compose exec app go run scripts/init_rbac.go

# 5. 验证服务
curl http://localhost:9060/health
```

访问 Swagger 文档：http://localhost:9060/swagger/index.html

### 💻 方式二：本地开发

适合需要调试和开发的场景：

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd gin-app-start

# 2. 安装依赖
go mod download

# 3. 启动数据库服务（使用 Docker）
docker run -d --name postgres \
  -e POSTGRES_DB=gin_app \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:17-alpine

docker run -d --name redis \
  -p 6379:6379 \
  redis:7-alpine

# 4. 配置环境变量（可选）
export SERVER_ENV=local  # 使用 configs/config.local.yaml

# 5. 运行应用
make run
# 或直接运行
go run cmd/server/main.go

# 6. 初始化 RBAC 数据（首次运行）
go run scripts/init_rbac.go

# 7. 验证服务
curl http://localhost:9060/health
```

### 🔥 方式三：热重载开发

使用 Air 实现代码修改后自动重启：

```bash
# 1. 安装 Air
go install github.com/cosmtrek/air@latest

# 2. 启动热重载
air

# 代码修改后会自动重新编译和运行
```

### 🎬 快速体验 RBAC 功能

```bash
# 1. 初始化 RBAC 数据（创建管理员账号和权限）
go run scripts/init_rbac.go

# 2. 登录获取 Token
curl -X POST http://localhost:9060/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 响应示例：
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
#     "user": { ... }
#   }
# }

# 3. 使用 Token 访问受保护接口
curl -X GET http://localhost:9060/api/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. 获取用户菜单树
curl -X GET http://localhost:9060/api/v1/menus/user/tree \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. 查看用户权限
curl -X GET http://localhost:9060/api/v1/permissions/user \
  -H "Authorization: Bearer YOUR_TOKEN"
```

更多 RBAC 使用说明请查看 [RBAC 快速开始](docs/RBAC_QUICKSTART.md)。

## 📁 项目结构

```
gin-app-start/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口，初始化所有组件
│
├── internal/                     # 私有应用代码（不可被外部导入）
│   ├── config/                  # 配置管理
│   │   └── config.go           # 配置结构定义和加载逻辑
│   │
│   ├── controller/              # HTTP 控制器层（处理请求/响应）
│   │   ├── enter.go            # 控制器组注册
│   │   ├── health_controller.go    # 健康检查
│   │   ├── auth_controller.go      # 认证登录
│   │   ├── user_controller.go      # 用户管理
│   │   ├── device_controller.go    # 设备管理
│   │   ├── role_controller.go      # 角色管理
│   │   ├── permission_controller.go # 权限管理
│   │   └── menu_controller.go      # 菜单管理
│   │
│   ├── service/                 # 业务逻辑层
│   │   ├── enter.go            # 服务组注册
│   │   ├── auth_service.go     # 认证业务逻辑
│   │   ├── user_service.go     # 用户业务逻辑
│   │   ├── device_service.go   # 设备业务逻辑
│   │   ├── role_service.go     # 角色业务逻辑
│   │   ├── permission_service.go # 权限业务逻辑
│   │   └── menu_service.go     # 菜单业务逻辑
│   │
│   ├── repository/              # 数据访问层（与数据库交互）
│   │   ├── enter.go            # 仓储组注册
│   │   ├── base_repository.go  # 泛型基础仓储
│   │   ├── user_repository.go  # 用户数据访问
│   │   ├── device_repository.go # 设备数据访问
│   │   ├── role_repository.go  # 角色数据访问
│   │   ├── permission_repository.go # 权限数据访问
│   │   └── menu_repository.go  # 菜单数据访问
│   │
│   ├── model/                   # 数据模型定义
│   │   ├── base.go             # 基础模型（ID、时间戳等）
│   │   ├── user.go             # 用户模型
│   │   ├── device.go           # 设备模型
│   │   ├── role.go             # 角色模型
│   │   ├── permission.go       # 权限模型
│   │   └── menu.go             # 菜单模型
│   │
│   ├── dto/                     # 数据传输对象（请求/响应结构）
│   │   ├── auth_dto.go         # 认证相关 DTO
│   │   ├── user_dto.go         # 用户相关 DTO
│   │   ├── device_dto.go       # 设备相关 DTO
│   │   ├── role_dto.go         # 角色相关 DTO
│   │   ├── permission_dto.go   # 权限相关 DTO
│   │   └── menu_dto.go         # 菜单相关 DTO
│   │
│   ├── middleware/              # Gin 中间件
│   │   ├── logger.go           # 日志中间件
│   │   ├── recovery.go         # 恢复中间件（捕获 panic）
│   │   ├── cors.go             # CORS 跨域中间件
│   │   ├── rate_limit.go       # 限流中间件
│   │   ├── auth.go             # JWT 认证中间件
│   │   └── permission.go       # 权限校验中间件
│   │
│   └── router/                  # 路由配置
│       └── router.go           # 路由注册和分组
│
├── pkg/                          # 可复用的公共库（可被外部导入）
│   ├── database/                # 数据库连接
│   │   ├── postgres.go         # PostgreSQL 连接池
│   │   └── redis.go            # Redis 连接池
│   │
│   ├── logger/                  # 日志封装
│   │   └── logger.go           # 基于 zap 的日志工具
│   │
│   ├── errors/                  # 错误定义
│   │   └── errors.go           # 业务错误码和错误类型
│   │
│   ├── response/                # 统一响应
│   │   └── response.go         # 标准化 JSON 响应格式
│   │
│   ├── jwt/                     # JWT 工具
│   │   └── jwt.go              # Token 生成和验证
│   │
│   └── utils/                   # 工具函数
│       ├── crypto.go           # 加密解密工具
│       └── utils.go            # 通用工具函数
│
├── configs/                      # 配置文件
│   ├── config.local.yaml        # 本地开发配置
│   ├── config.dev.yaml          # 开发环境配置
│   ├── config.prod.yaml         # 生产环境配置
│   ├── private_key.pem          # JWT 私钥
│   └── public_key.pem           # JWT 公钥
│
├── docs/                         # 项目文档
│   ├── PROJECT_GUIDE.md         # 项目使用指南
│   ├── API_REFERENCE.md         # API 接口文档
│   ├── ARCHITECTURE.md          # 架构设计文档
│   ├── RBAC.md                  # RBAC 权限系统文档
│   ├── RBAC_QUICKSTART.md       # RBAC 快速开始
│   ├── docs.go                  # Swagger 生成文件
│   ├── swagger.json             # Swagger JSON
│   └── swagger.yaml             # Swagger YAML
│
├── scripts/                      # 脚本工具
│   └── init_rbac.go             # RBAC 数据初始化脚本
│
├── .gitignore                    # Git 忽略文件
├── docker-compose.yml            # Docker Compose 配置
├── Dockerfile                    # Docker 镜像构建文件
├── Makefile                      # 常用命令快捷方式
├── go.mod                        # Go 模块依赖
├── go.sum                        # 依赖校验和
└── README.md                     # 项目说明文档
```

### 架构分层说明

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Request                        │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│  Controller Layer (控制器层)                             │
│  - 处理 HTTP 请求和响应                                   │
│  - 参数验证和绑定                                         │
│  - 调用 Service 层                                       │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│  Service Layer (业务逻辑层)                              │
│  - 核心业务逻辑处理                                       │
│  - 事务管理                                              │
│  - 调用 Repository 层                                    │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│  Repository Layer (数据访问层)                           │
│  - 数据库 CRUD 操作                                      │
│  - 查询构建                                              │
│  - 数据持久化                                            │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│              Database (PostgreSQL/Redis)                │
└─────────────────────────────────────────────────────────┘
```

## 🔧 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.24+ |
| Web 框架 | Gin | 1.10.0 |
| ORM | GORM | 1.25.12 |
| 数据库 | PostgreSQL | 17 |
| 缓存 | Redis | 7 |
| 日志 | Zap | 1.27.0 |
| 配置 | Viper | 1.19.0 |
| 文档 | Swagger | 1.16.3 |

## 📖 API 示例

### 健康检查

```bash
GET /health
```

### 用户管理

```bash
# 创建用户
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "phone": "13800138000",
  "password": "password123"
}

# 获取用户
GET /api/v1/users/:id

# 更新用户
PUT /api/v1/users/:id

# 删除用户
DELETE /api/v1/users/:id

# 用户列表（分页）
GET /api/v1/users?page=1&page_size=10
```

### 设备管理

```bash
# 创建设备
POST /api/v1/devices

# 获取设备
GET /api/v1/devices/:id

# 更新设备
PUT /api/v1/devices/:id

# 删除设备
DELETE /api/v1/devices/:id

# 设备列表
GET /api/v1/devices?page=1&page_size=10
```

### 统一响应格式

成功响应：
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

错误响应：
```json
{
  "code": 10001,
  "message": "参数错误",
  "data": null
}
```

分页响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

## ⚙️ 配置说明

配置文件位于 `configs/` 目录，通过环境变量 `SERVER_ENV` 选择：

```bash
SERVER_ENV=local    # configs/config.local.yaml
SERVER_ENV=dev      # configs/config.dev.yaml
SERVER_ENV=prod     # configs/config.prod.yaml
```

主要配置项：

```yaml
server:
  port: 9060              # 服务端口
  mode: debug             # 运行模式: debug/release/test
  read_timeout: 60        # 读超时（秒）
  write_timeout: 60       # 写超时（秒）
  limit_num: 100          # 限流：每秒请求数（0=不限流）

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: gin_app
  sslmode: disable
  max_idle_conns: 10      # 最大空闲连接
  max_open_conns: 100     # 最大打开连接
  max_lifetime: 3600      # 连接生命周期（秒）
  log_level: info         # 日志级别
  auto_migrate: true      # 自动迁移表结构

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3

log:
  level: debug            # debug/info/warn/error
  file_path: logs/app.log
  max_size: 100           # MB
  max_age: 7              # 天
```

## 🛠️ 开发指南

### Makefile 命令

```bash
make run            # 运行应用
make build          # 编译应用
make test           # 运行测试
make fmt            # 格式化代码
make lint           # 代码检查
make clean          # 清理编译文件
make deps           # 下载依赖
make docker-build   # 构建 Docker 镜像
make docker-run     # 运行 Docker 容器
make swagger        # 生成 Swagger 文档
make install-tools  # 安装开发工具
```

### 添加新功能

1. **定义模型** - `internal/model/your_model.go`
2. **创建 Repository** - `internal/repository/your_repository.go`
3. **实现 Service** - `internal/service/your_service.go`
4. **添加 Controller** - `internal/controller/your_controller.go`
5. **注册路由** - `internal/router/router.go`
6. **更新 main.go** - 初始化依赖注入

详细开发指南请参考 [项目使用指南](docs/PROJECT_GUIDE.md)。

### 错误处理

```go
import "gin-app-start/pkg/errors"

// 使用预定义错误
return errors.ErrUserNotFound

// 创建自定义错误
return errors.NewBusinessError(10001, "自定义错误")

// 包装错误
return errors.WrapBusinessError(10001, "操作失败", err)
```

### 日志记录

```go
import (
    "gin-app-start/pkg/logger"
    "go.uber.org/zap"
)

logger.Info("操作成功", 
    zap.String("username", username),
    zap.Uint("user_id", userID),
)

logger.Error("操作失败", zap.Error(err))
```

## 🐳 Docker 部署

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据
docker-compose down -v
```

### 单独构建和运行

```bash
# 构建镜像
docker build -t gin-app-start:latest .

# 运行容器
docker run -d -p 9060:9060 \
  -e SERVER_ENV=prod \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_ADDR=your-redis:6379 \
  --name gin-app \
  gin-app-start:latest
```

## 📊 性能优化建议

- 数据库连接池：根据 QPS 调整 `max_idle_conns` 和 `max_open_conns`
- Redis 缓存：合理使用缓存减少数据库查询
- 日志级别：生产环境使用 `warn` 或 `error` 级别
- 限流配置：根据服务器性能设置 `limit_num`
- 索引优化：为常用查询字段添加数据库索引

## 🔒 安全建议

- 生产环境使用 HTTPS
- 实现 JWT 或 Session 认证
- 密码使用 bcrypt 加密
- 启用 CORS 白名单
- 实施 API 限流
- 定期更新依赖包

## 📝 常见问题

| 问题 | 解决方案 |
|------|----------|
| 数据库连接失败 | 检查数据库是否启动，验证连接信息 |
| 端口被占用 | 修改配置文件中的端口号 |
| 依赖下载慢 | 设置 GOPROXY：`go env -w GOPROXY=https://goproxy.cn,direct` |
| 热重载不工作 | 安装 Air：`go install github.com/cosmtrek/air@latest` |

更多问题请查看 [项目使用指南](docs/PROJECT_GUIDE.md#常见问题)。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支：`git checkout -b feature/AmazingFeature`
3. 提交更改：`git commit -m 'Add some AmazingFeature'`
4. 推送到分支：`git push origin feature/AmazingFeature`
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 Web 框架
- [GORM](https://gorm.io) - 优秀的 ORM 库
- [Zap](https://github.com/uber-go/zap) - 高性能日志库
- [Viper](https://github.com/spf13/viper) - 配置管理库

---

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！
