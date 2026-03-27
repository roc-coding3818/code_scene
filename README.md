# CodeScene - 编程场景学习项目

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个聚焦于真实业务场景的 Go 后端学习项目，通过模块化方式学习各种技术实现。

</div>

## 项目简介

CodeScene 是一个专注于**业务场景化**的 Go 后端学习项目。每个模块对应一个真实的业务场景，让学习者能够：
- 掌握常见业务场景的技术实现
- 应对面试中的各种技术提问
- 学习企业级项目的最佳实践

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL** - 数据库
- **Redis** - 缓存

## 模块列表

### 第一阶段：核心基础

| 模块 | 描述 | 状态 |
|------|------|------|
| [用户系统](./modules/user/README.md) | 注册登录、JWT认证、验证码、限流 | 🚧 开发中 |
| [权限管理](./modules/auth/README.md) | RBAC、接口权限、按钮权限 | 📋 规划中 |
| [缓存系统](./modules/cache/README.md) | Redis数据类型、缓存策略、分布式锁 | 📋 规划中 |
| [接口设计](./modules/api-design/README.md) | RESTful、统一响应、参数校验 | 📋 规划中 |

### 第二阶段：业务核心

| 模块 | 描述 | 状态 |
|------|------|------|
| [订单系统](./modules/order/README.md) | 订单状态机、幂等设计、分布式事务 | 📋 规划中 |
| [支付系统](./modules/payment/README.md) | 微信/支付宝接入、回调验签 | 📋 规划中 |
| [商品模块](./modules/product/README.md) | 商品CRUD、SKU/SPU、搜索 | 📋 规划中 |
| [文件服务](./modules/file/README.md) | 本地存储、OSS接入、秒传 | 📋 规划中 |

### 第三阶段：高阶模块

| 模块 | 描述 | 状态 |
|------|------|------|
| [消息队列](./modules/mq/README.md) | RabbitMQ/Kafka、异步任务 | 📋 规划中 |
| [搜索系统](./modules/search/README.md) | ElasticSearch、聚合查询 | 📋 规划中 |
| [分布式系统](./modules/distributed/README.md) | 分布式ID、分布式锁 | 📋 规划中 |

## 项目结构

```
code_scene/
├── shared/                      # 通用代码（所有模块共享）
│   ├── config/                  # 配置加载
│   ├── db/                     # 数据库连接
│   ├── redis/                  # Redis连接
│   ├── response/               # 统一响应
│   ├── jwt/                    # JWT工具
│   ├── middleware/             # 中间件
│   └── utils/                  # 工具函数
│
├── modules/                     # 业务模块目录（每个模块独立运行）
│   ├── user/ (8080)             # 用户模块
│   │   ├── cmd/                 # 入口
│   │   ├── config/              # 配置
│   │   ├── sql/                 # 数据库脚本
│   │   ├── internal/            # 业务代码
│   │   │   ├── domain/          # 领域层
│   │   │   ├── repo/            # 数据访问层
│   │   │   └── service/         # 业务逻辑层
│   │   ├── global/              # 模块全局变量
│   │   ├── docs/                # 模块设计文档
│   │   └── go.mod               # 模块依赖
│   │
│   ├── auth/ (8081)             # 权限模块
│   ├── order/ (8082)            # 订单模块
│   ├── payment/ (8083)          # 支付模块
│   ├── product/ (8084)          # 商品模块
│   ├── cache/ (8085)            # 缓存模块
│   ├── file/ (8086)             # 文件模块
│   ├── mq/ (8087)              # 消息队列模块
│   ├── search/ (8088)           # 搜索模块
│   ├── distributed/ (8089)     # 分布式模块
│   └── api-design/ (8090)       # API设计模块
│
├── docs/                        # 项目通用设计文档
│   ├── architecture.md          # 架构设计
│   ├── database.md              # 数据库设计
│   └── api.md                   # 接口规范
│
├── go.work                      # Go工作区配置（根目录）
└── README.md                    # 项目总览
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 运行单个模块

以用户模块为例：

```bash
# 进入用户模块目录
cd modules/user

# 配置数据库
# 修改 config/config.yaml

# 初始化数据库
mysql -u root -p < sql/init.sql

# 运行
go run cmd/main.go
```

服务启动后访问 `http://localhost:8080`

## 学习路线

```
新手入门路线：
用户系统 → 权限管理 → 接口设计 → 缓存基础
     ↓                    ↓              ↓
进阶路线：
订单系统 → 支付系统 → 商品模块 → 文件服务
     ↓                    ↓
高阶路线：
消息队列 → 搜索系统 → 分布式进阶
```

## 核心特性

### 1. 模块独立
- 每个业务场景独立成模块，可单独学习、测试、运行
- 模块间通过清晰的接口通信

### 2. 贴近生产
- 遵循企业级项目规范
- 包含完整的错误处理、日志记录
- 适配真实业务场景

### 3. 面试导向
- 每个模块涵盖高频面试知识点
- 代码注释详细，便于理解

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## License

MIT License - 详见 [LICENSE](LICENSE)
