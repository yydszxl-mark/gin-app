# RBAC 系统快速开始指南

## 1. 启动应用

```bash
# 确保 PostgreSQL 和 Redis 已启动
docker-compose up -d

# 或者使用本地开发
make run
```

应用启动后会自动：
- 生成 JWT RSA 密钥对（`configs/private_key.pem` 和 `configs/public_key.pem`）
- 创建数据库表（users, roles, permissions, menus 及关联表）

## 2. 初始化 RBAC 数据

运行初始化脚本创建默认数据：

```bash
go run scripts/init_rbac.go
```

这将创建：
- **管理员角色**（admin）：拥有所有权限
- **开发者角色**（developer）：拥有查看和设备管理权限
- **管理员用户**：
  - 用户名：`admin`
  - 密码：`admin123`
- **28个权限**：涵盖用户、角色、权限、菜单、设备管理
- **8个菜单**：仪表盘、系统管理、设备管理等

## 3. 测试 API

### 3.1 用户登录

```bash
curl -X POST http://localhost:9060/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**响应示例：**
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
      "roles": ["管理员"]
    }
  }
}
```

保存返回的 `token`，后续请求需要使用。

### 3.2 获取当前用户信息

```bash
curl -X GET http://localhost:9060/api/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3.3 获取用户菜单树

```bash
curl -X GET http://localhost:9060/api/v1/menus/user/tree \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3.4 获取用户权限列表

```bash
curl -X GET http://localhost:9060/api/v1/permissions/user \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3.5 创建新用户

```bash
curl -X POST http://localhost:9060/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13900139000",
    "password": "test123"
  }'
```

### 3.6 分配角色给用户

```bash
curl -X POST http://localhost:9060/api/v1/users/2/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": [2]
  }'
```

### 3.7 创建新角色

```bash
curl -X POST http://localhost:9060/api/v1/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试角色",
    "code": "tester",
    "description": "测试人员角色",
    "sort": 3,
    "status": 1
  }'
```

### 3.8 分配权限给角色

```bash
curl -X POST http://localhost:9060/api/v1/roles/3/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 24, 25]
  }'
```

### 3.9 分配菜单给角色

```bash
curl -X POST http://localhost:9060/api/v1/roles/3/menus \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "menu_ids": [1, 3, 8]
  }'
```

### 3.10 修改密码

```bash
curl -X POST http://localhost:9060/api/v1/user/change-password \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "admin123",
    "new_password": "newpass123"
  }'
```

## 4. 使用 Postman 测试

### 4.1 导入环境变量

创建环境变量：
- `base_url`: `http://localhost:9060`
- `token`: 登录后获取的 token

### 4.2 设置 Authorization

在请求的 Headers 中添加：
```
Authorization: Bearer {{token}}
```

### 4.3 常用请求集合

**认证相关：**
- POST `{{base_url}}/api/v1/auth/login` - 登录
- POST `{{base_url}}/api/v1/auth/refresh` - 刷新 Token
- GET `{{base_url}}/api/v1/user/info` - 获取用户信息
- POST `{{base_url}}/api/v1/user/logout` - 登出

**用户管理：**
- GET `{{base_url}}/api/v1/users?page=1&page_size=10` - 用户列表
- POST `{{base_url}}/api/v1/users` - 创建用户
- GET `{{base_url}}/api/v1/users/:id` - 获取用户
- PUT `{{base_url}}/api/v1/users/:id` - 更新用户
- DELETE `{{base_url}}/api/v1/users/:id` - 删除用户
- POST `{{base_url}}/api/v1/users/:id/roles` - 分配角色

**角色管理：**
- GET `{{base_url}}/api/v1/roles?page=1&page_size=10` - 角色列表
- POST `{{base_url}}/api/v1/roles` - 创建角色
- GET `{{base_url}}/api/v1/roles/:id` - 获取角色
- PUT `{{base_url}}/api/v1/roles/:id` - 更新角色
- DELETE `{{base_url}}/api/v1/roles/:id` - 删除角色
- POST `{{base_url}}/api/v1/roles/:id/permissions` - 分配权限
- POST `{{base_url}}/api/v1/roles/:id/menus` - 分配菜单

**权限管理：**
- GET `{{base_url}}/api/v1/permissions?page=1&page_size=10` - 权限列表
- POST `{{base_url}}/api/v1/permissions` - 创建权限
- GET `{{base_url}}/api/v1/permissions/:id` - 获取权限
- PUT `{{base_url}}/api/v1/permissions/:id` - 更新权限
- DELETE `{{base_url}}/api/v1/permissions/:id` - 删除权限
- GET `{{base_url}}/api/v1/permissions/user` - 获取当前用户权限

**菜单管理：**
- GET `{{base_url}}/api/v1/menus` - 菜单列表
- POST `{{base_url}}/api/v1/menus` - 创建菜单
- GET `{{base_url}}/api/v1/menus/:id` - 获取菜单
- PUT `{{base_url}}/api/v1/menus/:id` - 更新菜单
- DELETE `{{base_url}}/api/v1/menus/:id` - 删除菜单
- GET `{{base_url}}/api/v1/menus/tree` - 获取菜单树
- GET `{{base_url}}/api/v1/menus/user/tree` - 获取用户菜单树

## 5. 数据库查看

### 5.1 查看用户及其角色

```sql
SELECT u.id, u.username, u.email, r.name as role_name, r.code as role_code
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.deleted_at IS NULL;
```

### 5.2 查看角色及其权限

```sql
SELECT r.name as role_name, p.name as permission_name, p.code, p.method, p.path
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE r.deleted_at IS NULL
ORDER BY r.id, p.sort;
```

### 5.3 查看角色及其菜单

```sql
SELECT r.name as role_name, m.title as menu_title, m.path, m.parent_id
FROM roles r
LEFT JOIN role_menus rm ON r.id = rm.role_id
LEFT JOIN menus m ON rm.menu_id = m.id
WHERE r.deleted_at IS NULL
ORDER BY r.id, m.sort;
```

## 6. 常见问题

### Q1: Token 无效或已过期？
**A:** Token 默认有效期为 7 天，过期后需要重新登录或使用刷新接口。

### Q2: 如何重置管理员密码？
**A:** 直接在数据库中更新：
```sql
-- 密码将被重置为 admin123
UPDATE users 
SET password = 'your_hashed_password', 
    salt = 'your_salt'
WHERE username = 'admin';
```
或者重新运行初始化脚本。

### Q3: 如何启用权限检查中间件？
**A:** 在 `internal/router/router.go` 中取消注释：
```go
apiV1.Use(middleware.PermissionCheck(permRepo))
```

### Q4: JWT 密钥丢失怎么办？
**A:** 删除 `configs/private_key.pem` 和 `configs/public_key.pem`，重启应用会自动生成新密钥。注意：所有旧 Token 将失效。

### Q5: 如何添加新的权限？
**A:** 
1. 通过 API 创建权限
2. 将权限分配给相应角色
3. 如需接口级权限控制，启用权限检查中间件

## 7. 下一步

- 查看完整文档：[RBAC.md](./RBAC.md)
- 查看 API 文档：http://localhost:9060/swagger/index.html
- 查看架构设计：[ARCHITECTURE.md](./ARCHITECTURE.md)
- 查看项目指南：[PROJECT_GUIDE.md](./PROJECT_GUIDE.md)

## 8. 技术支持

如有问题，请查看：
1. 应用日志：`logs/app.log`
2. 数据库日志
3. 项目文档

祝使用愉快！🎉
