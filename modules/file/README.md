# 文件服务模块

## 模块介绍

文件服务模块负责处理文件上传、存储、CDN分发等功能。

## 功能清单

### 文件上传
- [ ] 本地存储
- [ ] 对象存储（OSS/MinIO）
- [ ] 分片上传
- [ ] 秒传功能

### 文件管理
- [ ] 文件查询
- [ ] 文件删除
- [ ] 缩略图生成

### CDN 加速
- [ ] CDN 域名配置
- [ ] URL 签名

## 核心知识点

- 分片上传原理
- 断点续传
- 文件 Hash 校验
- CDN 缓存策略

## 目录结构

```
file/
├── cmd/
├── config/
├── internal/
│   ├── domain/
│   ├── repo/
│   ├── service/
│   │   ├── local.go
│   │   └── oss.go
│   └── handler.go
└── README.md
```
