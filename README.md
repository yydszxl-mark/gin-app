# Gin App Start

🚀 基于 Go 1.24 和 Gin 框架的企业级 Web 应用脚手架，采用清晰的分层架构设计，开箱即用。

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 核心特性

- 🏗️ **标准项目布局** - 遵循 Go 标准项目结构（cmd/internal/pkg）
- 📦 **三层架构** - Controller → Service → Repository 清晰分层
- 🗄️ **数据库支持** - PostgreSQL 17 + Redis 7 + GORM ORM
- 📝 **结构化日志** - 基于 uber/zap 的高性能日志系统
- 🛡️ **统一错误处理** - 自定义业务错误码和错误包装
- 🔄 **统一响应格式** - 标准化的 JSON API 响应
- 🔌 **丰富中间件** - 日志、恢复、CORS、限流等开箱即用
- ⚙️ **多环境配置** - 基于 Viper 的灵活配置管理
- 🐳 **容器化支持** - 完整的 Docker 和 Docker Compose 配置
- 📚 **Swagger 文档** - 自动生成 API 文档
- 🔄 **优雅关闭** - 支持 Graceful Shutdown
- 🔧 **泛型仓储** - 使用 Go 泛型实现通用数据访问层

## 📚 文档导航

| 文档 | 说明 |
|------|------|
| [项目使用指南](docs/PROJECT_GUIDE.md) | 详细的开发文档和最佳实践 |
| [API 接口文档](docs/API_REFERENCE.md) | 完整的 API 参考和示例 |
| [架构设计文档](docs/ARCHITECTURE.md) | 技术架构和设计理念 |

## 🚀 快速开始

### 前置要求

- Go 1.24+
- PostgreSQL 12+ (推荐 17)
- Redis 6.0+ (推荐 7)
- Docker & Docker Compose (可选)

### 方式一：使用 Docker Compose（推荐）

一键启动所有服务（应用 + PostgreSQL + Redis）：

```bash
# 克隆项目
git clone <your-repo-url>
cd gin-app-start

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 健康检查
curl http://localhost:9060/health
```

### 方式二：本地开发

```bash
# 1. 安装依赖
go mod download

# 2. 启动数据库（使用 Docker）
docker run -d --name postgres -e POSTGRES_DB=gin_app \
  -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:17-alpine

docker run -d --name redis -p 6379:6379 redis:7-alpine

# 3. 配置环境（编辑 configs/config.local.yaml）
# 根据实际情况修改数据库连接信息

# 4. 运行应用
make run
# 或
SERVER_ENV=local go run cmd/server/main.go

# 5. 验证运行
curl http://localhost:9060/health
```

### 方式三：热重载开发

```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 启动热重载
air
```

## 📁 项目结构

```
gin-app-start/
├── cmd/                      # 应用程序入口
│   └── server/main.go       # 主程序
├── internal/                 # 私有应用代码
│   ├── config/              # 配置管理
│   ├── controller/          # HTTP 控制器（处理请求/响应）
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据访问层（含泛型基类）
│   ├── model/               # 数据模型定义
│   ├── dto/                 # 数据传输对象
│   ├── middleware/          # Gin 中间件
│   └── router/              # 路由配置
├── pkg/                      # 可复用公共库
│   ├── database/            # 数据库连接（PostgreSQL/Redis）
│   ├── logger/              # 日志封装
│   ├── errors/              # 错误定义
│   ├── response/            # 统一响应
│   └── utils/               # 工具函数
├── configs/                  # 配置文件
│   ├── config.local.yaml    # 本地开发
│   ├── config.dev.yaml      # 开发环境
│   └── config.prod.yaml     # 生产环境
├── docs/                     # 文档和 Swagger
├── docker-compose.yml        # Docker Compose 配置
├── Dockerfile                # Docker 镜像构建
├── Makefile                  # 常用命令
└── go.mod                    # Go 模块依赖
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
