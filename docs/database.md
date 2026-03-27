# 数据库设计规范

## 命名规范

### 表命名
- 使用下划线命名法
- 模块名前缀：如 `user_`, `order_`, `product_`
- 示例：`users`, `user_profiles`, `order_items`

### 字段命名
- 使用下划线命名法
- 主键：`id`（自增）或 `{table}_id`（业务主键）
- 创建时间：`created_at`
- 更新时间：`updated_at`
- 软删除：`deleted_at`

### 索引命名
- 主键索引：`pk_{table}_{column}`
- 唯一索引：`uk_{table}_{column}`
- 普通索引：`idx_{table}_{column}`

## 表设计原则

### 1. 主键设计
- 使用 BIGINT 自增作为主键
- 业务主键需添加唯一索引

### 2. 软删除
- 使用 `deleted_at` 字段
- 查询时自动过滤已删除记录

### 3. 时间戳
- 所有表包含 `created_at` 和 `updated_at`
- 使用 DATETIME 类型

### 4. 状态字段
- 使用 TINYINT
- 添加注释说明各值含义

## 常用字段

```sql
-- 通用字段
`id` BIGINT NOT NULL AUTO_INCREMENT
`created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
`updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
`deleted_at` DATETIME DEFAULT NULL

-- 状态
`status` TINYINT DEFAULT 1 COMMENT '状态: 1-正常 2-禁用'

-- 用户
`user_id` BIGINT NOT NULL

-- 金额（使用 DECIMAL）
`amount` DECIMAL(10,2) NOT NULL
```

## 模块数据库

### 用户模块
- users: 用户表
- user_profiles: 用户详情表

### 订单模块
- orders: 订单表
- order_items: 订单商品表

### 支付模块
- payments: 支付记录表
- refunds: 退款记录表

### 商品模块
- categories: 分类表
- products: 商品SPU表
- product_skus: 商品SKU表
- product_attributes: 商品属性表
