package redis

import (
	"context"

	"code_scene/shared/config"
	"github.com/go-redis/redis/v8"
)

// Client 全局Redis客户端
var Client *redis.Client

// Init 初始化Redis连接
func Init(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return Client.Ping(context.Background()).Err()
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}
