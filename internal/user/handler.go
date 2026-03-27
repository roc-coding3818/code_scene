package user

import (
	"code_scene/global"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: NewUserService(),
	}
}

// Register 注册
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// PhoneRegister 手机号注册
func (h *UserHandler) PhoneRegister(c *gin.Context) {
	var req PhoneRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.PhoneRegister(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// Login 登录
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// PhoneLogin 手机号登录
func (h *UserHandler) PhoneLogin(c *gin.Context) {
	var req PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.PhoneLogin(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// SendEmailCode 发送邮箱验证码
func (h *UserHandler) SendEmailCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.SendEmailCode(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// SendPhoneCode 发送手机验证码
func (h *UserHandler) SendPhoneCode(c *gin.Context) {
	var req SendPhoneCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.SendPhoneCode(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// RefreshToken 刷新Token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	resp, err := h.service.GetUserInfo(c.Request.Context(), userID.(int64))
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	err := h.service.UpdateUserInfo(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.SuccessWithMsg(c, "更新成功", nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	err := h.service.ChangePassword(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.SuccessWithMsg(c, "密码修改成功", nil)
}

// SendResetPasswordCode 发送重置密码验证码
func (h *UserHandler) SendResetPasswordCode(c *gin.Context) {
	var req SendResetPasswordCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.service.SendResetPasswordCode(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.Success(c, resp)
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.SuccessWithMsg(c, "密码重置成功", nil)
}

// Logout 登出
func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	err := h.service.Logout(c.Request.Context(), token)
	if err != nil {
		global.Fail(c, 400, err.Error())
		return
	}

	global.SuccessWithMsg(c, "登出成功", nil)
}
