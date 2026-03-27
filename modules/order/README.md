# 订单系统模块

## 模块介绍

订单系统是电商业务的核心模块，涉及订单创建、状态管理、支付对接等核心流程。

## 功能清单

### 订单管理
- [ ] 订单创建
- [ ] 订单查询
- [ ] 订单取消
- [ ] 订单状态流转

### 订单详情
- [ ] 订单商品信息
- [ ] 收货地址
- [ ] 优惠信息
- [ ] 物流信息

### 库存控制
- [ ] 乐观锁防超卖
- [ ] 库存预占/释放
- [ ] 库存回滚

### 分布式事务
- [ ] 订单创建与库存扣减
- [ ] 订单状态同步

## 核心知识点

- 订单状态机设计
- 幂等性设计
- 乐观锁与悲观锁
- 分布式事务（Seata/TCC）
- 库存扣减方案

## 数据表

```sql
-- 订单表
CREATE TABLE `orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单号',
  `user_id` BIGINT NOT NULL,
  `total_amount` DECIMAL(10,2) NOT NULL,
  `pay_amount` DECIMAL(10,2) NOT NULL,
  `status` TINYINT NOT NULL COMMENT '订单状态',
  `remark` VARCHAR(500),
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`)
);

-- 订单商品表
CREATE TABLE `order_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `product_name` VARCHAR(100),
  `price` DECIMAL(10,2) NOT NULL,
  `quantity` INT NOT NULL,
  `subtotal` DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`)
);
```

## 目录结构

```
order/
├── cmd/
├── config/
├── internal/
│   ├── domain/
│   │   ├── order.go
│   │   └── item.go
│   ├── repo/
│   ├── service/
│   └── handler.go
└── README.md
```
