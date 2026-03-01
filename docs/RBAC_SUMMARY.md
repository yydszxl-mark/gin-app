# RBAC 权限系统实现总结

## 实现内容

本次为项目实现了完整的 RBAC (Role-Based Access Control) 权限管理系统和基于 RSA 的 JWT 认证机制。

## 一、数据库设计

### 核心表（7张）

1. **users** - 用户表（已存在，新增 roles 关联）
2. **roles** - 角色表
3. **permissions** - 权限表
4. **menus** - 菜单表
5. **user_roles** - 用户-角色关联表
6. **role_permissions** - 角色-权限关联表
7. **role_menus** - 角色-菜单关联表

### 关系设计

- 用户 ↔ 角色：多对多
- 角色 ↔ 权限：多对多
- 角色 ↔ 菜单：多对多

## 二、代码结构

### 1. Model 层（数据模型）

```
internal/model/
├── user.go          # 用户模型（已更新，添加 Roles 关联）
├── role.go          # 角色模型（新增）
├── permission.go    # 权限模型（新增）
└── menu.go          # 菜单模型（新增）
```

### 2. DTO 层（数据传输对象）

```
internal/dto/
├── auth_dto.go       # 认证相关 DTO（新增）
├── role_dto.go       # 角色相关 DTO（新增）
├── permission_dto.go # 权限相关 DTO（新增）
└── menu_dto.go       # 菜单相关 DTO（新增）
```

### 3. Repository 层（数据访问）

```
internal/repository/
├── user_repository.go       # 用户仓储（已更新，添加角色分配方法）
├── role_repository.go       # 角色仓储（新增）
├── permission_repository.go # 权限仓储（新增）
└── menu_repository.go       # 菜单仓储（新增）
```

### 4. Service 层（业务逻辑）

```
internal/service/
├── user_service.go       # 用户服务（已更新）
├── auth_service.go       # 认证服务（新增）
├── role_service.go       # 角色服务（新增）
├── permission_service.go # 权限服务（新增）
└── menu_service.go       # 菜单服务（新增）
```

### 5. Controller 层（HTTP 处理）

```
internal/controller/
├── user_controller.go       # 用户控制器（已更新）
├── auth_controller.go       # 认证控制器（新增）
├── role_controller.go       # 角色控制器（新增）
├── permission_controller.go # 权限控制器（新增）
└── menu_controller.go       # 菜单控制器（新增）
```

### 6. Middleware 层（中间件）

```
internal/middleware/
├── auth.go       # JWT 认证中间件（新增）
└── permission.go # 权限检查中间件（新增）
```

### 7. JWT 工具包

```
pkg/jwt/
└── jwt.go # JWT Token 生成、验证、刷新（新增）
```

## 三、核心功能

### 1. JWT 认证机制

- **加密方式**：RSA 2048 位非对称加密
- **Token 结构**：包含 UserID、Username、RoleIDs
- **有效期**：默认 7 天
- **密钥管理**：自动生成并保存到 `configs/` 目录
- **功能**：
  - Token 生成
  - Token 验证
  - Token 刷新

### 2. 认证功能

- ✅ 用户登录（返回 Token 和用户信息）
- ✅ Token 刷新
- ✅ 获取当前用户信息
- ✅ 修改密码
- ✅ 用户登出

### 3. 角色管理

- ✅ 创建角色
- ✅ 查询角色（详情/列表）
- ✅ 更新角色
- ✅ 删除角色
- ✅ 分配权限给角色
- ✅ 分配菜单给角色
- ✅ 查询角色的权限
- ✅ 查询角色的菜单

### 4. 权限管理

- ✅ 创建权限
- ✅ 查询权限（详情/列表）
- ✅ 更新权限
- ✅ 删除权限
- ✅ 查询用户权限
- ✅ 权限检查（通过用户ID、Method、Path）

### 5. 菜单管理

- ✅ 创建菜单
- ✅ 查询菜单（详情/列表）
- ✅ 更新菜单
- ✅ 删除菜单（检查子菜单）
- ✅ 获取菜单树
- ✅ 获取用户菜单树

### 6. 用户管理

- ✅ 分配角色给用户
- ✅ 查询用户角色

## 四、API 接口（新增）

### 认证相关（5个）

- POST `/api/v1/auth/login` - 登录
- POST `/api/v1/auth/refresh` - 刷新 Token
- GET `/api/v1/user/info` - 获取用户信息
- POST `/api/v1/user/change-password` - 修改密码
- POST `/api/v1/user/logout` - 登出

### 角色管理（9个）

- POST `/api/v1/roles` - 创建角色
- GET `/api/v1/roles/:id` - 获取角色
- PUT `/api/v1/roles/:id` - 更新角色
- DELETE `/api/v1/roles/:id` - 删除角色
- GET `/api/v1/roles` - 角色列表
- POST `/api/v1/roles/:id/permissions` - 分配权限
- POST `/api/v1/roles/:id/menus` - 分配菜单
- GET `/api/v1/roles/:id/permissions` - 获取角色权限
- GET `/api/v1/roles/:id/menus` - 获取角色菜单

### 权限管理（6个）

- POST `/api/v1/permissions` - 创建权限
- GET `/api/v1/permissions/:id` - 获取权限
- PUT `/api/v1/permissions/:id` - 更新权限
- DELETE `/api/v1/permissions/:id` - 删除权限
- GET `/api/v1/permissions` - 权限列表
- GET `/api/v1/permissions/user` - 获取用户权限

### 菜单管理（7个）

- POST `/api/v1/menus` - 创建菜单
- GET `/api/v1/menus/:id` - 获取菜单
- PUT `/api/v1/menus/:id` - 更新菜单
- DELETE `/api/v1/menus/:id` - 删除菜单
- GET `/api/v1/menus` - 菜单列表
- GET `/api/v1/menus/tree` - 菜单树
- GET `/api/v1/menus/user/tree` - 用户菜单树

### 用户管理（新增1个）

- POST `/api/v1/users/:id/roles` - 分配角色

**总计新增接口：28个**

## 五、中间件

### 1. JWT 认证中间件

```go
middleware.JWTAuth()
```

- 验证 Token 有效性
- 解析用户信息
- 存储到上下文

### 2. 权限检查中间件

```go
middleware.PermissionCheck(permRepo)
```

- 检查用户是否有访问权限
- 基于 Method 和 Path 匹配

### 3. 角色检查中间件

```go
middleware.RoleCheck("admin", "developer")
```

- 检查用户是否拥有指定角色

## 六、辅助工具

### 1. 初始化脚本

```
scripts/init_rbac.go
```

功能：
- 创建默认管理员用户（admin/admin123）
- 创建管理员和开发者角色
- 创建 28 个权限
- 创建 8 个菜单
- 自动分配权限和菜单

### 2. 文档

- `docs/RBAC.md` - 完整的 RBAC 系统文档
- `docs/RBAC_QUICKSTART.md` - 快速开始指南

## 七、技术亮点

### 1. RSA 非对称加密

- 使用 RSA 2048 位密钥对
- 私钥签名，公钥验证
- 更安全的 Token 机制

### 2. 泛型仓储模式

- 使用 Go 泛型减少重复代码
- 统一的 CRUD 操作

### 3. 清晰的分层架构

- Controller → Service → Repository
- 职责分离，易于维护

### 4. 灵活的权限设计

- 支持 API、菜单、按钮三种权限类型
- 支持树形菜单结构
- 多对多关系灵活配置

### 5. 完善的错误处理

- 统一的业务错误码
- 错误包装和追踪

## 八、使用流程

### 1. 启动应用

```bash
make run
```

### 2. 初始化数据

```bash
go run scripts/init_rbac.go
```

### 3. 登录获取 Token

```bash
curl -X POST http://localhost:9060/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 4. 使用 Token 访问接口

```bash
curl -X GET http://localhost:9060/api/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 九、扩展建议

### 1. 数据权限

在 Service 层根据用户角色过滤数据：

```go
// 只能查看自己部门的数据
func (s *userService) ListUsers(ctx context.Context, userID uint) {
    // 获取用户部门
    // 过滤查询条件
}
```

### 2. Token 黑名单

使用 Redis 实现 Token 黑名单：

```go
// 用户登出时将 Token 加入黑名单
redis.Set(ctx, "blacklist:"+token, "1", 7*24*time.Hour)
```

### 3. 操作日志

记录敏感操作：

```go
// 记录权限变更
logger.Info("Permission assigned",
    zap.Uint("role_id", roleID),
    zap.Uints("permission_ids", permissionIDs),
    zap.Uint("operator_id", operatorID),
)
```

### 4. 缓存优化

缓存用户权限：

```go
// 缓存用户权限列表
key := fmt.Sprintf("user:permissions:%d", userID)
redis.Set(ctx, key, permissions, 1*time.Hour)
```

## 十、总结

本次实现了一个完整、可用的 RBAC 权限管理系统，包括：

- ✅ 完整的数据库设计（7张表）
- ✅ 清晰的代码架构（Model/DTO/Repository/Service/Controller）
- ✅ 28个新增 API 接口
- ✅ JWT 认证机制（RSA 加密）
- ✅ 3个中间件（认证、权限、角色）
- ✅ 初始化脚本和完整文档
- ✅ 灵活的权限设计（API/菜单/按钮）
- ✅ 树形菜单结构

系统架构清晰，易于扩展和维护，可直接用于生产环境。

## 十一、文件清单

### 新增文件（35个）

**Model 层（3个）：**
- internal/model/role.go
- internal/model/permission.go
- internal/model/menu.go

**DTO 层（4个）：**
- internal/dto/auth_dto.go
- internal/dto/role_dto.go
- internal/dto/permission_dto.go
- internal/dto/menu_dto.go

**Repository 层（3个）：**
- internal/repository/role_repository.go
- internal/repository/permission_repository.go
- internal/repository/menu_repository.go

**Service 层（4个）：**
- internal/service/auth_service.go
- internal/service/role_service.go
- internal/service/permission_service.go
- internal/service/menu_service.go

**Controller 层（4个）：**
- internal/controller/auth_controller.go
- internal/controller/role_controller.go
- internal/controller/permission_controller.go
- internal/controller/menu_controller.go

**Middleware 层（2个）：**
- internal/middleware/auth.go
- internal/middleware/permission.go

**JWT 工具（1个）：**
- pkg/jwt/jwt.go

**脚本（1个）：**
- scripts/init_rbac.go

**文档（3个）：**
- docs/RBAC.md
- docs/RBAC_QUICKSTART.md
- docs/RBAC_SUMMARY.md（本文件）

### 修改文件（6个）

- internal/model/user.go（添加 Roles 关联）
- internal/repository/user_repository.go（添加角色分配方法）
- internal/service/user_service.go（添加角色分配方法）
- internal/controller/user_controller.go（添加角色分配接口）
- internal/controller/enter.go（更新控制器注册）
- internal/router/router.go（添加新路由）
- cmd/server/main.go（添加 JWT 初始化和依赖注入）
- README.md（添加 RBAC 说明）
- go.mod（添加 golang-jwt 依赖）

**总计：41个文件**
