# RBAC 权限系统文档

## 概述

本项目实现了完整的 RBAC (Role-Based Access Control) 权限管理系统，支持用户、角色、权限、菜单的管理，并集成了基于 RSA 的 JWT 认证。

## 数据库设计

### 核心表结构

#### 1. users (用户表)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE,
    phone VARCHAR(32) UNIQUE,
    password VARCHAR(128) NOT NULL,
    salt VARCHAR(32) NOT NULL,
    avatar VARCHAR(256),
    status SMALLINT DEFAULT 1
);
```

#### 2. roles (角色表)
```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) UNIQUE NOT NULL,
    description VARCHAR(255),
    sort INT DEFAULT 0,
    status INT DEFAULT 1
);
```

#### 3. permissions (权限表)
```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(128) UNIQUE NOT NULL,
    type VARCHAR(32) NOT NULL,  -- menu, button, api
    method VARCHAR(16),          -- GET, POST, PUT, DELETE
    path VARCHAR(255),
    description VARCHAR(255),
    sort INT DEFAULT 0,
    status INT DEFAULT 1
);
```

#### 4. menus (菜单表)
```sql
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    parent_id INT DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    title VARCHAR(64) NOT NULL,
    icon VARCHAR(128),
    path VARCHAR(255),
    component VARCHAR(255),
    redirect VARCHAR(255),
    type VARCHAR(32) DEFAULT 'menu',  -- menu, button
    hidden BOOLEAN DEFAULT false,
    sort INT DEFAULT 0,
    status INT DEFAULT 1
);
```

### 关联表

#### 5. user_roles (用户-角色关联表)
```sql
CREATE TABLE user_roles (
    user_id INT,
    role_id INT,
    PRIMARY KEY (user_id, role_id)
);
```

#### 6. role_permissions (角色-权限关联表)
```sql
CREATE TABLE role_permissions (
    role_id INT,
    permission_id INT,
    PRIMARY KEY (role_id, permission_id)
);
```

#### 7. role_menus (角色-菜单关联表)
```sql
CREATE TABLE role_menus (
    role_id INT,
    menu_id INT,
    PRIMARY KEY (role_id, menu_id)
);
```

## API 接口文档

### 认证相关

#### 1. 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1735689600,
    "user_info": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "phone": "13800138000",
      "avatar": "",
      "status": 1,
      "roles": ["管理员", "开发者"]
    }
  }
}
```

#### 2. 刷新 Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3. 获取当前用户信息
```http
GET /api/v1/user/info
Authorization: Bearer {token}
```

#### 4. 修改密码
```http
POST /api/v1/user/change-password
Authorization: Bearer {token}
Content-Type: application/json

{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

#### 5. 登出
```http
POST /api/v1/user/logout
Authorization: Bearer {token}
```

### 角色管理

#### 1. 创建角色
```http
POST /api/v1/roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "开发者",
  "code": "developer",
  "description": "开发人员角色",
  "sort": 1,
  "status": 1
}
```

#### 2. 获取角色详情
```http
GET /api/v1/roles/:id
Authorization: Bearer {token}
```

#### 3. 更新角色
```http
PUT /api/v1/roles/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "高级开发者",
  "description": "高级开发人员",
  "sort": 2,
  "status": 1
}
```

#### 4. 删除角色
```http
DELETE /api/v1/roles/:id
Authorization: Bearer {token}
```

#### 5. 角色列表
```http
GET /api/v1/roles?page=1&page_size=10
Authorization: Bearer {token}
```

#### 6. 分配权限给角色
```http
POST /api/v1/roles/:id/permissions
Authorization: Bearer {token}
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 4]
}
```

#### 7. 分配菜单给角色
```http
POST /api/v1/roles/:id/menus
Authorization: Bearer {token}
Content-Type: application/json

{
  "menu_ids": [1, 2, 3]
}
```

#### 8. 获取角色的权限
```http
GET /api/v1/roles/:id/permissions
Authorization: Bearer {token}
```

#### 9. 获取角色的菜单
```http
GET /api/v1/roles/:id/menus
Authorization: Bearer {token}
```

### 权限管理

#### 1. 创建权限
```http
POST /api/v1/permissions
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "创建用户",
  "code": "user:create",
  "type": "api",
  "method": "POST",
  "path": "/api/v1/users",
  "description": "创建用户权限",
  "sort": 1,
  "status": 1
}
```

#### 2. 获取权限详情
```http
GET /api/v1/permissions/:id
Authorization: Bearer {token}
```

#### 3. 更新权限
```http
PUT /api/v1/permissions/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "创建用户（更新）",
  "description": "创建用户权限（已更新）",
  "status": 1
}
```

#### 4. 删除权限
```http
DELETE /api/v1/permissions/:id
Authorization: Bearer {token}
```

#### 5. 权限列表
```http
GET /api/v1/permissions?page=1&page_size=10
Authorization: Bearer {token}
```

#### 6. 获取当前用户权限
```http
GET /api/v1/permissions/user
Authorization: Bearer {token}
```

### 菜单管理

#### 1. 创建菜单
```http
POST /api/v1/menus
Authorization: Bearer {token}
Content-Type: application/json

{
  "parent_id": 0,
  "name": "system",
  "title": "系统管理",
  "icon": "setting",
  "path": "/system",
  "component": "Layout",
  "redirect": "/system/user",
  "type": "menu",
  "hidden": false,
  "sort": 1,
  "status": 1
}
```

#### 2. 获取菜单详情
```http
GET /api/v1/menus/:id
Authorization: Bearer {token}
```

#### 3. 更新菜单
```http
PUT /api/v1/menus/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "系统设置",
  "icon": "settings",
  "sort": 2
}
```

#### 4. 删除菜单
```http
DELETE /api/v1/menus/:id
Authorization: Bearer {token}
```

#### 5. 菜单列表
```http
GET /api/v1/menus
Authorization: Bearer {token}
```

#### 6. 获取菜单树
```http
GET /api/v1/menus/tree
Authorization: Bearer {token}
```

#### 7. 获取当前用户菜单树
```http
GET /api/v1/menus/user/tree
Authorization: Bearer {token}
```

### 用户管理

#### 分配角色给用户
```http
POST /api/v1/users/:id/roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "role_ids": [1, 2]
}
```

## JWT 认证机制

### RSA 密钥对

系统使用 RSA 2048 位密钥对进行 JWT 签名和验证：

- **私钥**：用于签名 Token（存储在 `configs/private_key.pem`）
- **公钥**：用于验证 Token（存储在 `configs/public_key.pem`）

首次启动时，如果密钥文件不存在，系统会自动生成。

### Token 结构

```go
type MyClaims struct {
    UserID   uint     `json:"user_id"`
    Username string   `json:"username"`
    RoleIDs  []uint   `json:"role_ids"`
    jwt.RegisteredClaims
}
```

### Token 使用

1. **获取 Token**：用户登录成功后获取
2. **使用 Token**：在请求头中添加 `Authorization: Bearer {token}`
3. **Token 有效期**：默认 7 天
4. **刷新 Token**：在 Token 过期前可以刷新获取新 Token

## 权限检查流程

### 1. JWT 认证中间件

```go
middleware.JWTAuth()
```

- 从请求头获取 Token
- 验证 Token 有效性
- 解析用户信息并存入上下文

### 2. 权限检查中间件（可选）

```go
middleware.PermissionCheck(permRepo)
```

- 获取当前用户 ID
- 获取请求的 Method 和 Path
- 查询用户是否有对应权限

### 3. 角色检查中间件（可选）

```go
middleware.RoleCheck("admin", "developer")
```

- 检查用户是否拥有指定角色

## 使用示例

### 1. 初始化系统

```bash
# 启动应用
make run

# 系统会自动：
# 1. 生成 JWT 密钥对
# 2. 创建数据库表
# 3. 启动服务
```

### 2. 创建管理员用户

```bash
curl -X POST http://localhost:9060/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### 3. 创建角色

```bash
curl -X POST http://localhost:9060/api/v1/roles \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "管理员",
    "code": "admin",
    "description": "系统管理员",
    "status": 1
  }'
```

### 4. 创建权限

```bash
curl -X POST http://localhost:9060/api/v1/permissions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "用户管理",
    "code": "user:manage",
    "type": "api",
    "method": "GET",
    "path": "/api/v1/users",
    "status": 1
  }'
```

### 5. 分配权限给角色

```bash
curl -X POST http://localhost:9060/api/v1/roles/1/permissions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3]
  }'
```

### 6. 分配角色给用户

```bash
curl -X POST http://localhost:9060/api/v1/users/1/roles \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": [1]
  }'
```

### 7. 用户登录

```bash
curl -X POST http://localhost:9060/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

## 最佳实践

### 1. 权限设计

- **API 权限**：`resource:action` 格式，如 `user:create`, `user:update`
- **菜单权限**：对应前端路由，如 `/system/user`
- **按钮权限**：细粒度控制，如 `user:delete:button`

### 2. 角色设计

- **超级管理员**：拥有所有权限
- **部门管理员**：拥有部门内权限
- **普通用户**：基础权限

### 3. 安全建议

- 定期更换 JWT 密钥
- Token 有效期不宜过长
- 敏感操作需要二次验证
- 实施 API 访问频率限制
- 记录权限变更日志

### 4. 性能优化

- 使用 Redis 缓存用户权限
- 权限检查结果缓存
- 批量查询优化

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 10001 | 参数错误 |
| 10002 | 资源不存在 |
| 10003 | 未授权访问 |
| 10004 | 资源已存在 |
| 10005 | 数据库操作失败 |
| 50000 | 系统内部错误 |

## 常见问题

### 1. Token 过期怎么办？

使用刷新 Token 接口获取新的 Token。

### 2. 如何实现单点登录？

可以在 Redis 中维护 Token 黑名单，用户登出时将 Token 加入黑名单。

### 3. 如何实现数据权限？

在 Service 层根据用户角色过滤数据，例如只能查看自己部门的数据。

### 4. 密钥丢失怎么办？

删除旧密钥文件，重启应用会自动生成新密钥，但所有旧 Token 将失效。

## 总结

本 RBAC 系统提供了完整的权限管理功能，支持：

- ✅ 用户、角色、权限、菜单的 CRUD 操作
- ✅ 灵活的权限分配机制
- ✅ 基于 RSA 的 JWT 认证
- ✅ 中间件级别的权限控制
- ✅ 树形菜单结构
- ✅ 多对多关系管理

系统架构清晰，易于扩展和维护。
