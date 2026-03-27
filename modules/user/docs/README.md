# 用户模块

## 模块简介

用户模块是整个项目的基础模块，提供用户注册、登录、会话管理等核心功能。通过本模块可以学习到：

- JWT 认证机制
- 双 Token 刷新机制
- 验证码设计与存储
- 密码加密存储
- 接口限流实现

## 快速开始

### 1. 环境准备

```bash
# 确保已安装 Go 1.21+
go version

# 启动 MySQL (端口 3306)
# 启动 Redis (端口 6379)
```

### 2. 配置数据库

编辑 `config/config.yaml`，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: code_scene_user  # 用户模块专用数据库
```

### 3. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS code_scene_user;"

# 执行建表脚本
mysql -u root -p code_scene_user < sql/init.sql
```

### 4. 启动服务

```bash
cd modules/user
go run cmd/main.go
```

服务启动后访问 `http://localhost:8080`

---

## 功能清单与代码对照

### 1. 用户注册

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| 邮箱注册 | 邮箱+验证码注册 | `internal/handler/user.go` - `RegisterByEmail` |
| 手机注册 | 手机号+验证码注册 | `internal/handler/user.go` - `RegisterByPhone` |
| 频率限制 | 防止恶意注册 | `shared/middleware/ratelimit.go` |

**实现要点：**
- 验证码存入 Redis，设置过期时间（5分钟）
- 密码使用 bcrypt 加密存储
- 注册前校验用户名/邮箱/手机号唯一性

### 2. 用户登录

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| 账号密码登录 | 用户名+密码登录 | `internal/handler/user.go` - `LoginByPassword` |
| 手机验证码登录 | 手机号+验证码登录 | `internal/handler/user.go` - `LoginByPhone` |
| 登录限流 | 防止暴力破解 | `shared/middleware/ratelimit.go` |

**实现要点：**
- 密码错误次数限制（5次后锁定15分钟）
- 登录成功后生成 Access Token 和 Refresh Token

### 3. 会话管理

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| Access Token | 短期访问凭证(2小时) | `shared/jwt/jwt.go` - `GenerateToken` |
| Refresh Token | 长期刷新凭证(7天) | `shared/jwt/jwt.go` - `GenerateRefreshToken` |
| Token 刷新 | 刷新过期 Token | `internal/handler/user.go` - `RefreshToken` |
| Token 失效 | 退出登录 | `internal/handler/user.go` - `Logout` |
| Token 验证 | 请求拦截验证 | `shared/middleware/middleware.go` - `JWT()` |

**实现要点：**
- Access Token 短期有效，用于接口访问
- Refresh Token 长期有效，用于刷新 Access Token
- Redis 存储 Token 黑名单，支持强制登出

### 4. 找回密码

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| 邮箱找回 | 邮箱+验证码重置密码 | `internal/handler/user.go` - `ResetPasswordByEmail` |
| 手机找回 | 手机+验证码重置密码 | `internal/handler/user.go` - `ResetPasswordByPhone` |

### 5. 用户信息管理

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| 获取信息 | 获取当前用户信息 | `internal/handler/user.go` - `GetUserInfo` |
| 更新信息 | 修改昵称/头像等 | `internal/handler/user.go` - `UpdateUserInfo` |
| 修改密码 | 修改登录密码 | `internal/handler/user.go` - `ChangePassword` |

### 6. 验证码

| 功能 | 说明 | 代码位置 |
|------|------|----------|
| 发送邮箱验证码 | 发送注册/登录验证码 | `internal/handler/user.go` - `SendEmailCode` |
| 发送手机验证码 | 发送注册/登录验证码 | `internal/handler/user.go` - `SendPhoneCode` |

**验证码存储：**
- Key: `code:{type}:{email|phone}`
- Value: 6位数字验证码
- 过期时间: 300秒（5分钟）
- 发送间隔: 60秒

---

## API 接口文档

### 基础信息

- 基础路径: `/api/v1/user`
- 认证方式: JWT (Header: `Authorization: Bearer <token>`)

### 接口列表

#### 1. 发送邮箱验证码

```
POST /api/v1/user/send/email/code
Content-Type: application/json

Request:
{
    "email": "user@example.com",
    "scene": "register"  // register/login/reset
}

Response:
{
    "code": 0,
    "message": "success"
}
```

#### 2. 发送手机验证码

```
POST /api/v1/user/send/phone/code
Content-Type: application/json

Request:
{
    "phone": "13800138000",
    "scene": "register"
}

Response:
{
    "code": 0,
    "message": "success"
}
```

#### 3. 邮箱注册

```
POST /api/v1/user/register/email
Content-Type: application/json

Request:
{
    "email": "user@example.com",
    "password": "123456",
    "code": "123456"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": 1,
        "access_token": "xxx",
        "refresh_token": "xxx"
    }
}
```

#### 4. 手机号注册

```
POST /api/v1/user/register/phone
Content-Type: application/json

Request:
{
    "phone": "13800138000",
    "password": "123456",
    "code": "123456"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": 1,
        "access_token": "xxx",
        "refresh_token": "xxx"
    }
}
```

#### 5. 账号密码登录

```
POST /api/v1/user/login
Content-Type: application/json

Request:
{
    "username": "user@example.com",
    "password": "123456"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": 1,
        "access_token": "xxx",
        "refresh_token": "xxx"
    }
}
```

#### 6. 手机验证码登录

```
POST /api/v1/user/login/phone
Content-Type: application/json

Request:
{
    "phone": "13800138000",
    "code": "123456"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": 1,
        "access_token": "xxx",
        "refresh_token": "xxx"
    }
}
```

#### 7. 刷新 Token

```
POST /api/v1/user/refresh/token
Content-Type: application/json

Request:
{
    "refresh_token": "xxx"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "xxx",
        "refresh_token": "xxx"
    }
}
```

#### 8. 获取用户信息

```
GET /api/v1/user/info
Authorization: Bearer <access_token>

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "user@example.com",
        "email": "user@example.com",
        "phone": "13800138000",
        "nickname": "用户昵称",
        "avatar": "https://...",
        "status": 1,
        "created_at": "2024-01-01 00:00:00"
    }
}
```

#### 9. 更新用户信息

```
PUT /api/v1/user/info
Authorization: Bearer <access_token>
Content-Type: application/json

Request:
{
    "nickname": "新昵称",
    "avatar": "https://..."
}

Response:
{
    "code": 0,
    "message": "success"
}
```

#### 10. 修改密码

```
PUT /api/v1/user/password
Authorization: Bearer <access_token>
Content-Type: application/json

Request:
{
    "old_password": "123456",
    "new_password": "abcdef"
}

Response:
{
    "code": 0,
    "message": "success"
}
```

#### 11. 找回密码-邮箱

```
POST /api/v1/user/reset/password/email
Content-Type: application/json

Request:
{
    "email": "user@example.com",
    "code": "123456",
    "new_password": "abcdef"
}

Response:
{
    "code": 0,
    "message": "success"
}
```

#### 12. 退出登录

```
POST /api/v1/user/logout
Authorization: Bearer <access_token>

Response:
{
    "code": 0,
    "message": "success"
}
```

---

## 错误码

### 通用错误 (1-99)

| 错误码 | 说明 |
|--------|------|
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 用户模块错误 (1001-1999)

| 错误码 | 说明 |
|--------|------|
| 1001 | 用户不存在 |
| 1002 | 用户已存在 |
| 1003 | 密码错误 |
| 1004 | 验证码错误 |
| 1005 | 验证码已过期 |
| 1006 | 验证码发送过快 |
| 1007 | Token无效 |
| 1008 | Token已过期 |
| 1009 | Refresh Token无效 |
| 1010 | 登录失败次数过多 |
| 1011 | 用户已被禁用 |

---

## 目录结构

```
user/
├── cmd/                         # 程序入口
│   └── main.go                 # 启动文件
│
├── config/                      # 配置
│   ├── config.go               # 配置结构定义
│   └── config.yaml             # 配置文件
│
├── sql/                        # 数据库脚本
│   └── init.sql                # 建表脚本
│
├── internal/                   # 业务代码
│   ├── domain/                 # 领域层
│   │   └── user.go             # 用户实体定义
│   │
│   ├── repo/                   # 数据访问层
│   │   ├── mysql.go            # MySQL 操作
│   │   └── redis.go            # Redis 操作
│   │
│   ├── service/                # 业务逻辑层
│   │   └── user.go             # 用户服务
│   │
│   └── handler/                # 接口处理层
│       └── user.go             # 用户接口处理
│
├── global/                     # 模块全局变量
│   ├── global.go               # 数据库/Redis 初始化
│   └── response.go             # 统一响应封装
│
├── docs/                       # 模块设计文档
│   └── README.md               # 本文档
│
├── go.mod                      # 模块依赖
└── README.md                   # 模块说明
```

---

## 核心技术点

### 1. JWT 双 Token 机制

```
┌─────────────────────────────────────────────────────────┐
│                      Access Token                        │
│  - 有效期: 2小时                                         │
│  - 用途: 接口访问                                        │
│  - 存储: 前端内存/LocalStorage                           │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                     Refresh Token                        │
│  - 有效期: 7天                                          │
│  - 用途: 刷新 Access Token                               │
│  - 存储: 前端 HttpOnly Cookie (更安全)                   │
└─────────────────────────────────────────────────────────┘
```

**刷新流程：**
```
1. Access Token 过期
2. 使用 Refresh Token 调用刷新接口
3. 返回新的 Access Token 和 Refresh Token
4. 继续正常访问
```

### 2. 验证码设计

```
Redis Key 格式:
  - 邮箱验证码: code:email:{email}
  - 手机验证码: code:phone:{phone}

存储内容:
  - 6位数字验证码
  - 过期时间: 300秒
  - 发送间隔: 60秒
```

### 3. 密码安全

- 密码使用 bcrypt 加密存储
- 登录错误5次后锁定15分钟
- 错误次数存入 Redis: `login:fail:{username}`

### 4. 接口限流

- 基于 Redis 实现滑动窗口限流
- 验证码发送: 每IP每手机号60秒限流
- 登录: 每IP每账号5次/15分钟

---

## 常见问题

### Q: 如何在生产环境部署？

1. 修改 `config/config.yaml` 中的数据库密码
2. 修改 JWT Secret 为复杂的随机字符串
3. 使用 `go build -o user cmd/main.go` 编译
4. 配置 Nginx 反向代理

### Q: 如何连接其他模块？

模块间通过 HTTP API 调用，例如：
```go
// 调用订单服务
resp, err := http.Post("http://localhost:8082/api/v1/order/list", ...)
```

### Q: 如何添加新功能？

1. 在 `internal/domain/` 定义实体
2. 在 `internal/repo/` 实现数据访问
3. 在 `internal/service/` 实现业务逻辑
4. 在 `internal/handler/` 实现接口
5. 在 `internal/router/` 注册路由

---

## 扩展学习

### 进阶话题

- [ ] SSO 单点登录
- [ ] OAuth2 第三方登录
- [ ] 短信验证码服务集成
- [ ] 邮件发送服务集成
- [ ] 用户画像与推荐系统

### 相关模块

- [权限模块](../auth/README.md) - RBAC 权限控制
- [缓存模块](../cache/README.md) - Redis 高级应用
