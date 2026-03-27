package global

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
	})
}

// FailWithCode 使用预定义错误码
func FailWithCode(c *gin.Context, errCode ErrCode) {
	Fail(c, errCode.Code, errCode.Message)
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, data *PageData) {
	Success(c, data)
}

// ErrCode 错误码
type ErrCode struct {
	Code    int
	Message string
}

func (e ErrCode) Error() string {
	return e.Message
}

// CommonErr 通用错误
var (
	ParamError            = ErrCode{Code: 400, Message: "参数错误"}
	Unauthorized          = ErrCode{Code: 401, Message: "未授权"}
	Forbidden             = ErrCode{Code: 403, Message: "禁止访问"}
	NotFound              = ErrCode{Code: 404, Message: "资源不存在"}
	ServerError           = ErrCode{Code: 500, Message: "服务器内部错误"}
)

// UserErr 用户模块错误
var (
	UserNotFound         = ErrCode{Code: 1001, Message: "用户不存在"}
	UserAlreadyExist     = ErrCode{Code: 1002, Message: "用户已存在"}
	PasswordError        = ErrCode{Code: 1003, Message: "密码错误"}
	CodeError            = ErrCode{Code: 1004, Message: "验证码错误"}
	CodeExpired          = ErrCode{Code: 1005, Message: "验证码已过期"}
	CodeSendTooFast      = ErrCode{Code: 1006, Message: "验证码发送过快"}
	TokenInvalid         = ErrCode{Code: 1007, Message: "Token无效"}
	TokenExpired         = ErrCode{Code: 1008, Message: "Token已过期"}
	RefreshTokenInvalid  = ErrCode{Code: 1009, Message: "Refresh Token无效"}
)
