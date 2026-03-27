package user

import (
	"context"
	"fmt"
	"time"

	"code_scene/global"

	"github.com/go-redis/redis/v8"
)

// RedisKeyPrefix Redis Key前缀
const (
	CodePrefix         = "user:code:"         // 验证码前缀
	TokenBlacklistPrefix = "user:token:blacklist:" // Token黑名单前缀
	LoginFailPrefix    = "user:login:fail:"  // 登录失败次数前缀
	SendCodeRateLimitPrefix = "user:send:code:"  // 发送验证码限流前缀
)

// UserRedisRepo 用户Redis仓库
type UserRedisRepo struct {
	client *redis.Client
}

// NewUserRedisRepo 创建用户Redis仓库
func NewUserRedisRepo() *UserRedisRepo {
	return &UserRedisRepo{client: global.Redis}
}

// SetEmailCode 设置邮箱验证码
func (r *UserRedisRepo) SetEmailCode(ctx context.Context, email, code string, expire time.Duration) error {
	return r.client.Set(ctx, CodePrefix+"email:"+email, code, expire).Err()
}

// GetEmailCode 获取邮箱验证码
func (r *UserRedisRepo) GetEmailCode(ctx context.Context, email string) (string, error) {
	return r.client.Get(ctx, CodePrefix+"email:"+email).Result()
}

// DelEmailCode 删除邮箱验证码
func (r *UserRedisRepo) DelEmailCode(ctx context.Context, email string) error {
	return r.client.Del(ctx, CodePrefix+"email:"+email).Err()
}

// SetPhoneCode 设置手机验证码
func (r *UserRedisRepo) SetPhoneCode(ctx context.Context, phone, code string, expire time.Duration) error {
	return r.client.Set(ctx, CodePrefix+"phone:"+phone, code, expire).Err()
}

// GetPhoneCode 获取手机验证码
func (r *UserRedisRepo) GetPhoneCode(ctx context.Context, phone string) (string, error) {
	return r.client.Get(ctx, CodePrefix+"phone:"+phone).Result()
}

// SetTokenBlacklist 设置Token到黑名单
func (r *UserRedisRepo) SetTokenBlacklist(ctx context.Context, token string, expire time.Duration) error {
	return r.client.Set(ctx, TokenBlacklistPrefix+token, "1", expire).Err()
}

// IsTokenBlacklist 检查Token是否在黑名单
func (r *UserRedisRepo) IsTokenBlacklist(ctx context.Context, token string) (bool, error) {
	result, err := r.client.Get(ctx, TokenBlacklistPrefix+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "1", nil
}

// IncrLoginFailCount 增加登录失败次数
func (r *UserRedisRepo) IncrLoginFailCount(ctx context.Context, username string, expire time.Duration) (int64, error) {
	key := LoginFailPrefix + username
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	// 首次设置过期时间
	if count == 1 {
		r.client.Expire(ctx, key, expire)
	}
	return count, nil
}

// GetLoginFailCount 获取登录失败次数
func (r *UserRedisRepo) GetLoginFailCount(ctx context.Context, username string) (int64, error) {
	count, err := r.client.Get(ctx, LoginFailPrefix+username).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// ResetLoginFailCount 重置登录失败次数
func (r *UserRedisRepo) ResetLoginFailCount(ctx context.Context, username string) error {
	return r.client.Del(ctx, LoginFailPrefix+username).Err()
}

// SetSendCodeRateLimit 设置发送验证码限流
func (r *UserRedisRepo) SetSendCodeRateLimit(ctx context.Context, key string, expire time.Duration) error {
	// 使用Redis的SETNX实现限流
	result, err := r.client.SetNX(ctx, SendCodeRateLimitPrefix+key, "1", expire).Result()
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("rate limit exceeded")
	}
	return nil
}

// CheckSendCodeRateLimit 检查发送验证码限流
func (r *UserRedisRepo) CheckSendCodeRateLimit(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, SendCodeRateLimitPrefix+key).Result()
	if err != nil {
		return false, err
	}
	return result == 0, nil
}
