# 权限管理模块

## 模块介绍

权限管理模块负责系统的访问控制，实现基于角色的权限管理（RBAC）。

## 功能清单

### 角色管理
- [ ] 角色 CRUD
- [ ] 角色权限分配

### 用户角色
- [ ] 用户角色分配
- [ ] 用户角色查询

### 权限控制
- [ ] 接口权限校验
- [ ] 按钮权限控制
- [ ] 数据权限隔离

## 核心知识点

- RBAC 权限模型
- 中间件权限校验
- 菜单树结构设计
- 分布式下的权限管理

## 数据表

```sql
-- 角色表
CREATE TABLE `roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  `code` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255),
  `status` TINYINT DEFAULT 1,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
);

-- 权限表
CREATE TABLE `permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  `code` VARCHAR(100) NOT NULL,
  `type` TINYINT COMMENT '1:菜单 2:按钮 3:接口',
  `parent_id` BIGINT DEFAULT 0,
  `path` VARCHAR(255),
  `component` VARCHAR(255),
  `sort` INT DEFAULT 0,
  `status` TINYINT DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
);

-- 用户角色关联表
CREATE TABLE `user_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `role_id` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
);

-- 角色权限关联表
CREATE TABLE `role_permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL,
  `permission_id` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
);
```

## 目录结构

```
auth/
├── cmd/
├── config/
├── internal/
│   ├── domain/
│   │   ├── role.go
│   │   └── permission.go
│   ├── repo/
│   ├── service/
│   └── middleware/
└── README.md
```
