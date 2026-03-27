package handler

import (
	"net/http"
	"code_scene/modules/user/internal/service"
	"code_scene/modules/user/global"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 注册
// @Summary 用户注册
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.Success(c, resp)
}

// Login 登录
// @Summary 用户登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.Success(c, resp)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.Response
// @Router /api/v1/user/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	req := &service.GetUserInfoRequest{
		UserID: userID.(int64),
	}

	resp, err := h.userService.GetUserInfo(c.Request.Context(), req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.Success(c, resp)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateUserInfoRequest true "更新请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	var req service.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}
	req.UserID = userID.(int64)

	resp, err := h.userService.UpdateUserInfo(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.Success(c, resp)
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdatePasswordRequest true "修改密码请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	var req service.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}
	req.UserID = userID.(int64)

	err := h.userService.UpdatePassword(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.SuccessWithMsg(c, "密码修改成功", nil)
}

// SendCode 发送短信验证码
// @Summary 发送短信验证码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body service.SendCodeRequest true "发送验证码请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/send_code [post]
func (h *UserHandler) SendCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	err := h.userService.SendCode(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.SuccessWithMsg(c, "验证码发送成功", nil)
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body service.RefreshTokenRequest true "刷新Token请求"
// @Success 200 {object} global.Response
// @Router /api/v1/user/refresh_token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FailWithCode(c, global.ParamError)
		return
	}

	resp, err := h.userService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(global.ErrCode); ok {
			global.FailWithCode(c, e)
			return
		}
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.Success(c, resp)
}

// Logout 登出
// @Summary 登出
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.Response
// @Router /api/v1/user/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		global.FailWithCode(c, global.Unauthorized)
		return
	}

	err := h.userService.Logout(c.Request.Context(), userID.(int64))
	if err != nil {
		global.FailWithCode(c, global.ServerError)
		return
	}

	global.SuccessWithMsg(c, "登出成功", nil)
}

// Ping 健康检查
// @Summary 健康检查
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} global.Response
// @Router /ping [get]
func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "pong",
	})
}
