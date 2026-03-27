package middleware

import (
	"code_scene/config"
	"code_scene/global"
	userRepo "code_scene/internal/user/repo"
	"code_scene/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	jwtUtil := jwt.NewJWT(&config.AppConfig.JWT)
	redisRepo := userRepo.NewUserRedisRepo()

	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			global.FailWithCode(c, global.Unauthorized)
			c.Abort()
			return
		}

		// 解析Token (Bearer xxx)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			global.FailWithCode(c, global.Unauthorized)
			c.Abort()
			return
		}

		token := parts[1]

		// 检查Token是否在黑名单
		isBlacklist, _ := redisRepo.IsTokenBlacklist(c.Request.Context(), token)
		if isBlacklist {
			global.FailWithCode(c, global.TokenInvalid)
			c.Abort()
			return
		}

		// 解析Token
		claims, err := jwtUtil.ValidateToken(token)
		if err != nil {
			if err.Error() == "Token has been revoked" {
				global.FailWithCode(c, global.TokenInvalid)
			} else {
				global.FailWithCode(c, global.TokenExpired)
			}
			c.Abort()
			return
		}

		// 将用户信息存入Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
