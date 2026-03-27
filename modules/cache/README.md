# 缓存系统模块

## 模块介绍

缓存系统模块专注于 Redis 的各种应用场景，包括缓存策略、分布式锁等。

## 功能清单

### 缓存应用
- [ ] 热点数据缓存
- [ ] 分布式 Session
- [ ] 缓存预热/预加载

### 分布式锁
- [ ] 互斥锁（SetNX）
- [ ] 可重入锁（Redisson）
- [ ] 公平锁

### 特殊数据结构
- [ ] 延迟队列
- [ ] 排行榜（ZSet）
- [ ] 消息队列（List）
- [ ] 计数器

### 缓存策略
- [ ] Cache Aside
- [ ] Read/Write Through
- [ ] 延迟双删
- [ ] 缓存穿透/击穿/雪崩

## 核心知识点

- Redis 数据类型与底层实现
- 缓存一致性方案
- 分布式锁原理与实现
- Redis 集群与哨兵

## 目录结构

```
cache/
├── cmd/
├── config/
├── internal/
│   ├── domain/
│   ├── repo/
│   ├── service/
│   └── handler.go
└── README.md
```
