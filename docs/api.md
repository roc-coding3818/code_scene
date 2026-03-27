# 接口设计规范

## RESTful 规范

### URL 设计
```
GET    /users              # 获取用户列表
GET    /users/:id          # 获取单个用户
POST   /users              # 创建用户
PUT    /users/:id          # 更新用户
DELETE /users/:id          # 删除用户
```

### HTTP 方法

| 方法 | 说明 | 幂等 |
|------|------|------|
| GET | 查询 | 是 |
| POST | 创建 | 否 |
| PUT | 完整更新 | 是 |
| PATCH | 部分更新 | 否 |
| DELETE | 删除 | 是 |

### 状态码

#### 2xx 成功
- 200 OK
- 201 Created
- 204 No Content

#### 4xx 客户端错误
- 400 Bad Request - 参数错误
- 401 Unauthorized - 未授权
- 403 Forbidden - 禁止访问
- 404 Not Found - 资源不存在
- 422 Unprocessable Entity - 参数校验失败

#### 5xx 服务端错误
- 500 Internal Server Error

## 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "test"
  }
}
```

### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 失败响应
```json
{
  "code": 400,
  "message": "参数错误"
}
```

## 错误码规范

### 通用错误（1-99）
- 1: 成功
- 400: 参数错误
- 401: 未授权
- 403: 禁止访问
- 404: 资源不存在
- 500: 服务器内部错误

### 模块错误（1000+）
- 用户模块: 1001-1999
- 订单模块: 2001-2999
- 支付模块: 3001-3999

## 接口版本

通过 URL 路径版本控制：
```
/api/v1/users
/api/v2/users
```

## 参数校验

### 请求参数
- 必填参数使用 `binding:"required"`
- 字符串长度限制
- 数值范围限制
- 格式校验（邮箱、手机号等）

### 示例
```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6,max=20"`
    Email    string `json:"email" binding:"required,email"`
}
```
