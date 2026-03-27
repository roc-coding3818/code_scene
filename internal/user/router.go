package user

import (
	"code_scene/pkg/middleware"

	"github.com/gin-gonic/gin"
)

var userHandler *UserHandler

// RegisterRoutes 注册用户模块路由
func RegisterRoutes(group *gin.RouterGroup) {
	if userHandler == nil {
		userHandler = NewUserHandler()
	}

	// 公开接口
	public := group.Group("/user")
	{
		// 注册
		public.POST("/register", userHandler.Register)
		// 手机号注册
		public.POST("/phone/register", userHandler.PhoneRegister)
		// 登录
		public.POST("/login", userHandler.Login)
		// 手机号登录
		public.POST("/phone/login", userHandler.PhoneLogin)
		// 发送邮箱验证码
		public.POST("/send/email/code", userHandler.SendEmailCode)
		// 发送手机验证码
		public.POST("/send/phone/code", userHandler.SendPhoneCode)
		// 刷新Token
		public.POST("/refresh/token", userHandler.RefreshToken)
		// 发送重置密码验证码
		public.POST("/send/reset/password/code", userHandler.SendResetPasswordCode)
		// 重置密码
		public.POST("/reset/password", userHandler.ResetPassword)
	}

	// 需要认证的接口
	auth := group.Group("/user")
	auth.Use(middleware.JWTAuth())
	{
		// 获取用户信息
		auth.GET("/info", userHandler.GetUserInfo)
		// 更新用户信息
		auth.PUT("/info", userHandler.UpdateUserInfo)
		// 修改密码
		auth.PUT("/password", userHandler.ChangePassword)
		// 登出
		auth.POST("/logout", userHandler.Logout)
	}
}
