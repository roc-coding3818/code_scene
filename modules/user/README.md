# 用户模块

## 模块介绍

用户模块是整个项目的基础模块，负责用户相关的所有业务功能。

## 功能清单

### 用户注册
- [ ] 邮箱注册 + 验证码
- [ ] 手机号注册 + 短信验证码
- [ ] 注册频率限制（防刷）

### 用户登录
- [ ] 账号密码登录
- [ ] 短信验证码登录
- [ ] 登录限流（防暴力破解）

### 会话管理
- [ ] JWT Access Token + Refresh Token
- [ ] Token 黑名单
- [ ] 多端登录控制

### 找回密码
- [ ] 邮箱找回 + 验证码
- [ ] 手机找回 + 验证码

### 用户信息
- [ ] 个人信息 CRUD
- [ ] 修改密码
- [ ] 头像上传

## 核心知识点

- JWT 原理与实现
- Access Token 与 Refresh Token 双 Token 机制
- 验证码设计与存储（Redis）
- 登录限流（防暴力破解）
- 密码加密存储（bcrypt）
- Session 与 Token 对比

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/user/register | 邮箱注册 |
| POST | /api/user/phone/register | 手机号注册 |
| POST | /api/user/login | 账号密码登录 |
| POST | /api/user/phone/login | 手机号登录 |
| POST | /api/user/send/email/code | 发送邮箱验证码 |
| POST | /api/user/send/phone/code | 发送手机验证码 |
| POST | /api/user/refresh/token | 刷新 Token |
| GET | /api/user/info | 获取用户信息 |
| PUT | /api/user/info | 更新用户信息 |
| PUT | /api/user/password | 修改密码 |
| POST | /api/user/logout | 登出 |

## 数据表

```sql
CREATE TABLE `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100) DEFAULT NULL,
  `phone` VARCHAR(20) DEFAULT NULL,
  `nickname` VARCHAR(50) DEFAULT NULL,
  `avatar` VARCHAR(255) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 目录结构

```
user/
├── cmd/
│   └── main.go                  # 入口文件
├── config/
│   └── config.yaml              # 配置文件
├── internal/
│   ├── domain/                  # 领域层
│   │   └── user.go              # 用户实体
│   ├── repo/                    # 数据访问层
│   │   ├── mysql.go             # MySQL 仓库
│   │   └── redis.go             # Redis 仓库
│   ├── service/                 # 业务逻辑层
│   │   └── user.go              # 用户服务
│   └── handler.go               # 接口处理层
├── pkg/                         # 模块内公共包
├── sql/
│   └── init.sql                 # 数据库初始化脚本
└── README.md
```
