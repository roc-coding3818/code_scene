package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"code_scene/modules/user/internal/domain"
	"code_scene/modules/user/internal/repo"
	"code_scene/modules/user/global"
	"code_scene/shared/jwt"

	"github.com/go-redis/redis/v8"
	"code_scene/shared/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repo.UserRepo
	jwtCfg   *config.JWTConfig
}

func NewUserService(userRepo *repo.UserRepo, jwtCfg *config.JWTConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Password   string `json:"password" binding:"required,min=6,max=20"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"required"`
	Nickname   string `json:"nickname"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captcha_key" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	User *domain.UserInfo `json:"user"`
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 1. 验证图形验证码
	if err := s.verifyCaptcha(ctx, req.CaptchaKey, req.Captcha); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否存在
	exist, err := s.userRepo.CheckUsernameExist(ctx, req.Username)
	if err != nil {
		return nil, global.ServerError
	}
	if exist {
		return nil, global.UserAlreadyExist
	}

	// 3. 检查邮箱是否存在
	exist, err = s.userRepo.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return nil, global.ServerError
	}
	if exist {
		return nil, global.UserAlreadyExist
	}

	// 4. 检查手机号是否存在
	exist, err = s.userRepo.CheckPhoneExist(ctx, req.Phone)
	if err != nil {
		return nil, global.ServerError
	}
	if exist {
		return nil, global.UserAlreadyExist
	}

	// 5. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, global.ServerError
	}

	// 6. 创建用户
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: nickname,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, global.ServerError
	}

	// 7. 删除图形验证码
	s.deleteCaptcha(ctx, req.CaptchaKey)

	return &RegisterResponse{
		User: user.ToVO(),
	}, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User         *domain.UserInfo `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresAt    int64            `json:"expires_at"`
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, global.ServerError
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, global.PasswordError
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		return nil, global.Forbidden
	}

	// 4. 生成Token
	accessToken, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, global.ServerError
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, global.ServerError
	}

	// 5. 存储RefreshToken到Redis
	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, global.ServerError
	}

	expiresAt := time.Now().Add(time.Duration(s.jwtCfg.AccessExpire) * time.Second).Unix()

	return &LoginResponse{
		User:         user.ToVO(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserID int64 `json:"user_id"`
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*domain.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, global.ServerError
	}

	return user.ToVO(), nil
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, req *UpdateUserInfoRequest) (*domain.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, global.ServerError
	}

	// 检查邮箱是否被占用
	if req.Email != "" && req.Email != user.Email {
		exist, err := s.userRepo.CheckEmailExist(ctx, req.Email)
		if err != nil {
			return nil, global.ServerError
		}
		if exist {
			return nil, global.UserAlreadyExist
		}
		user.Email = req.Email
	}

	// 检查手机号是否被占用
	if req.Phone != "" && req.Phone != user.Phone {
		exist, err := s.userRepo.CheckPhoneExist(ctx, req.Phone)
		if err != nil {
			return nil, global.ServerError
		}
		if exist {
			return nil, global.UserAlreadyExist
		}
		user.Phone = req.Phone
	}

	// 更新其他字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, global.ServerError
	}

	return user.ToVO(), nil
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	UserID      int64  `json:"user_id"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.UserNotFound
		}
		return global.ServerError
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return global.PasswordError
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return global.ServerError
	}

	if err := s.userRepo.UpdatePassword(ctx, req.UserID, string(hashedPassword)); err != nil {
		return global.ServerError
	}

	// 使其RefreshToken失效
	s.deleteRefreshToken(ctx, req.UserID)

	return nil
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
	Type  string `json:"type"` // register, login, reset_password
}

// SendCode 发送短信验证码
func (s *UserService) SendCode(ctx context.Context, req *SendCodeRequest) error {
	// 生成6位数字验证码
	code := generateCode()

	// 存储验证码到Redis，5分钟有效
	key := "sms_code:" + req.Phone + ":" + req.Type
	if err := global.Redis.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		return global.ServerError
	}

	// TODO: 调用短信服务发送验证码
	// 这里模拟发送成功，实际项目中需要集成短信服务
	println("SMS Code:", code) // 实际项目中删除

	return nil
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
	Type  string `json:"type"`
}

// VerifyCode 验证短信验证码
func (s *UserService) VerifyCode(ctx context.Context, req *VerifyCodeRequest) (bool, error) {
	key := "sms_code:" + req.Phone + ":" + req.Type
	storedCode, err := global.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, global.CodeExpired
	}
	if err != nil {
		return false, global.ServerError
	}

	if storedCode != req.Code {
		return false, global.CodeError
	}

	// 验证成功后删除验证码
	s.deleteCaptcha(ctx, key)

	return true, nil
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新Token
func (s *UserService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error) {
	// 1. 验证RefreshToken
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, global.RefreshTokenInvalid
	}

	// 2. 检查RefreshToken是否在Redis中
	key := "refresh_token:" + string(rune(claims.UserID))
	if _, err := global.Redis.Get(ctx, key).Result(); err != nil {
		return nil, global.RefreshTokenInvalid
	}

	// 3. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, global.ServerError
	}

	// 4. 检查用户状态
	if user.Status != 1 {
		return nil, global.Forbidden
	}

	// 5. 生成新Token
	accessToken, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, global.ServerError
	}

	// 6. 生成新的RefreshToken
	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, global.ServerError
	}

	// 7. 更新Redis中的RefreshToken
	s.storeRefreshToken(ctx, user.ID, newRefreshToken)

	expiresAt := time.Now().Add(time.Duration(s.jwtCfg.AccessExpire) * time.Second).Unix()

	return &LoginResponse{
		User:         user.ToVO(),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout 登出
func (s *UserService) Logout(ctx context.Context, userID int64) error {
	return s.deleteRefreshToken(ctx, userID)
}

// storeRefreshToken 存储RefreshToken到Redis
func (s *UserService) storeRefreshToken(ctx context.Context, userID int64, token string) error {
	key := "refresh_token:" + string(rune(userID))
	expire := time.Duration(s.jwtCfg.RefreshExpire) * time.Second
	return global.Redis.Set(ctx, key, token, expire).Err()
}

// deleteRefreshToken 删除RefreshToken
func (s *UserService) deleteRefreshToken(ctx context.Context, userID int64) error {
	key := "refresh_token:" + string(rune(userID))
	return global.Redis.Del(ctx, key).Err()
}

// verifyCaptcha 验证图形验证码
func (s *UserService) verifyCaptcha(ctx context.Context, key, captcha string) error {
	storedCaptcha, err := global.Redis.Get(ctx, "captcha:"+key).Result()
	if err == redis.Nil {
		return global.CodeExpired
	}
	if err != nil {
		return global.ServerError
	}

	if storedCaptcha != captcha {
		return global.CodeError
	}

	return nil
}

// deleteCaptcha 删除图形验证码
func (s *UserService) deleteCaptcha(ctx context.Context, key string) error {
	return global.Redis.Del(ctx, "captcha:"+key).Err()
}

// generateCode 生成6位数字验证码
func generateCode() string {
	code := time.Now().UnixNano() % 1000000
	return fmt.Sprintf("%06d", code)
}
