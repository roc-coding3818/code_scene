package main

import (
	"code_scene/modules/user/config"
	"code_scene/modules/user/global"
	"code_scene/modules/user/handler"
	"code_scene/modules/user/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库
	if err := global.InitDB(&config.AppConfig.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 3. 初始化 Redis
	if err := global.InitRedis(&config.AppConfig.Redis); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	log.Println("Redis连接成功")

	// 4. 初始化 Gin
	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 5. 注册路由
	router.RegisterRoutes(r)

	// 6. 启动服务
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("服务启动成功，监听地址: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
