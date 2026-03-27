# 支付系统模块

## 模块介绍

支付系统负责对接第三方支付渠道，处理支付订单的创建、回调、对账等核心流程。

## 功能清单

### 支付渠道
- [ ] 微信支付（Native/H5/JSAPI）
- [ ] 支付宝（PC/手机/WAP）
- [ ] 银联支付

### 支付流程
- [ ] 支付下单
- [ ] 支付回调处理
- [ ] 支付结果查询
- [ ] 支付取消/退款

### 安全保障
- [ ] 回调验签
- [ ] 回调幂等
- [ ] 金额校验
- [ ] 回调通知

### 财务对账
- [ ] 日对账文件下载
- [ ] 账单核对
- [ ] 异常订单处理

## 核心知识点

- 第三方支付接入流程
- 支付回调验签
- 支付幂等设计
- 资金安全要点
- 对账系统设计

## 数据表

```sql
-- 支付记录表
CREATE TABLE `payments` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `payment_no` VARCHAR(50) NOT NULL COMMENT '支付流水号',
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单号',
  `amount` DECIMAL(10,2) NOT NULL,
  `channel` VARCHAR(20) NOT NULL COMMENT '支付渠道',
  `channel_order_no` VARCHAR(100) COMMENT '渠道订单号',
  `status` TINYINT NOT NULL COMMENT '支付状态',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_no` (`payment_no`),
  KEY `idx_order_no` (`order_no`)
);

-- 退款记录表
CREATE TABLE `refunds` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `refund_no` VARCHAR(50) NOT NULL,
  `payment_id` BIGINT NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL,
  `reason` VARCHAR(500),
  `status` TINYINT,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_no` (`refund_no`)
);
```

## 目录结构

```
payment/
├── cmd/
├── config/
├── internal/
│   ├── domain/
│   │   ├── payment.go
│   │   └── refund.go
│   ├── repo/
│   ├── service/
│   │   ├── wechat.go
│   │   └── alipay.go
│   └── handler.go
└── README.md
```
