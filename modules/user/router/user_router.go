package router

import (
	"code_scene/modules/user/handler"
	"code_scene/modules/user/internal/repo"
	"code_scene/modules/user/internal/service"
	"code_scene/shared/config"
	"code_scene/shared/jwt"
	"code_scene/shared/middleware"

	"github.com/gin-gonic/gin"
)

var (
	userHandler *handler.UserHandler
)

// Init 初始化路由
func init() {
	// 初始化repo
	userRepo := repo.NewUserRepo()

	// 初始化service (JWT配置将在RegisterRoutes中传入)
	userService := service.NewUserService(userRepo, nil)

	// 初始化handler
	userHandler = handler.NewUserHandler(userService)
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, jwtCfg *config.JWTConfig) {
	// 初始化JWT
	jwt.Init(jwtCfg)

	// 健康检查
	r.GET("/ping", userHandler.Ping)

	// 用户模块
	v1 := r.Group("/api/v1/user")
	{
		// 公开接口
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/send_code", userHandler.SendCode)
		v1.POST("/refresh_token", userHandler.RefreshToken)

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.JWT())
		{
			auth.GET("/info", userHandler.GetUserInfo)
			auth.PUT("/info", userHandler.UpdateUserInfo)
			auth.PUT("/password", userHandler.UpdatePassword)
			auth.POST("/logout", userHandler.Logout)
		}
	}
}
