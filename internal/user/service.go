package user

import (
	"context"
	"errors"
	"time"

	"code_scene/config"
	"code_scene/global"
	"code_scene/internal/domain"
	"code_scene/internal/user/repo"
	"code_scene/pkg/jwt"
	"code_scene/pkg/utils"

	"gorm.io/gorm"
)

const (
	CodeExpire       = 300  // 验证码有效期5分钟
	LoginMaxFail     = 5    // 最大登录失败次数
	LoginFailExpire  = 900  // 登录失败计数器过期时间15分钟
	SendCodeInterval = 60   // 发送验证码间隔60秒
)

// UserService 用户服务
type UserService struct {
	userRepo  *repo.UserRepo
	userRedis *repo.UserRedisRepo
	jwtUtil   *jwt.JWT
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		userRepo:  repo.NewUserRepo(),
		userRedis: repo.NewUserRedisRepo(),
		jwtUtil:   jwt.NewJWT(&config.AppConfig.JWT),
	}
}

// Register 邮箱注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 1. 校验验证码
	code, err := s.userRedis.GetEmailCode(ctx, req.Email)
	if err != nil {
		return nil, global.CodeExpired
	}
	if code != req.Code {
		return nil, global.CodeError
	}
	// 验证后删除验证码（防止重复使用）
	s.userRedis.DelEmailCode(ctx, req.Email)

	// 2. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, global.UserAlreadyExist
	}

	// 3. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, global.UserAlreadyExist
	}

	// 4. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. 创建用户
	user := &domain.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 6. 删除验证码
	s.userRedis.DelEmailCode(ctx, req.Email)

	return &RegisterResponse{
		UserID: user.ID,
	}, nil
}

// PhoneRegister 手机号注册
func (s *UserService) PhoneRegister(ctx context.Context, req *PhoneRegisterRequest) (*RegisterResponse, error) {
	// 1. 校验验证码
	code, err := s.userRedis.GetPhoneCode(ctx, req.Phone)
	if err != nil {
		return nil, global.CodeExpired
	}
	if code != req.Code {
		return nil, global.CodeError
	}

	// 2. 检查手机号是否已存在
	exists, err := s.userRepo.ExistsByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, global.UserAlreadyExist
	}

	// 3. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户
	user := &domain.User{
		Username: "phone_" + req.Phone,
		Password: hashedPassword,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		UserID: user.ID,
	}, nil
}

// Login 账号密码登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 检查登录失败次数
	failCount, _ := s.userRedis.GetLoginFailCount(ctx, req.Username)
	if failCount >= int64(LoginMaxFail) {
		return nil, errors.New("登录尝试过多，请15分钟后再试")
	}

	// 2. 获取用户（支持用户名/邮箱/手机号）
	var user *domain.User
	var err error

	user, err = s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = s.userRepo.GetByEmail(req.Username)
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 增加登录失败次数
				s.userRedis.IncrLoginFailCount(ctx, req.Username, time.Duration(LoginFailExpire)*time.Second)
				return nil, global.PasswordError
			}
			return nil, err
		}
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		// 增加登录失败次数
		s.userRedis.IncrLoginFailCount(ctx, req.Username, time.Duration(LoginFailExpire)*time.Second)
		return nil, global.PasswordError
	}

	// 4. 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 5. 重置登录失败次数
	s.userRedis.ResetLoginFailCount(ctx, req.Username)

	// 6. 生成Token
	accessToken, refreshToken, err := s.jwtUtil.GenerateAccessAndRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.AppConfig.JWT.AccessExpire,
		UserInfo:     user.ToVO(),
	}, nil
}

// PhoneLogin 手机号验证码登录
func (s *UserService) PhoneLogin(ctx context.Context, req *PhoneLoginRequest) (*LoginResponse, error) {
	// 1. 校验验证码
	code, err := s.userRedis.GetPhoneCode(ctx, req.Phone)
	if err != nil {
		return nil, global.CodeExpired
	}
	if code != req.Code {
		return nil, global.CodeError
	}

	// 2. 获取用户
	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, err
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 4. 生成Token
	accessToken, refreshToken, err := s.jwtUtil.GenerateAccessAndRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.AppConfig.JWT.AccessExpire,
		UserInfo:     user.ToVO(),
	}, nil
}

// SendEmailCode 发送邮箱验证码
func (s *UserService) SendEmailCode(ctx context.Context, req *SendCodeRequest) (*SendCodeResponse, error) {
	// 1. 检查发送频率
	err := s.userRedis.SetSendCodeRateLimit(ctx, "email:"+req.Email, time.Duration(SendCodeInterval)*time.Second)
	if err != nil {
		return nil, global.CodeSendTooFast
	}

	// 2. 生成验证码
	code := utils.GenerateRandomCode(6)

	// 3. 存储验证码
	err = s.userRedis.SetEmailCode(ctx, req.Email, code, time.Duration(CodeExpire)*time.Second)
	if err != nil {
		return nil, err
	}

	// TODO: 实际发送邮件
	// 这里只打印验证码，实际项目中应该调用邮件服务发送
	println("Email Code:", code)

	return &SendCodeResponse{
		ExpireSeconds: CodeExpire,
	}, nil
}

// SendPhoneCode 发送手机验证码
func (s *UserService) SendPhoneCode(ctx context.Context, req *SendPhoneCodeRequest) (*SendCodeResponse, error) {
	// 1. 检查发送频率
	err := s.userRedis.SetSendCodeRateLimit(ctx, "phone:"+req.Phone, time.Duration(SendCodeInterval)*time.Second)
	if err != nil {
		return nil, global.CodeSendTooFast
	}

	// 2. 生成验证码
	code := utils.GenerateRandomCode(6)

	// 3. 存储验证码
	err = s.userRedis.SetPhoneCode(ctx, req.Phone, code, time.Duration(CodeExpire)*time.Second)
	if err != nil {
		return nil, err
	}

	// TODO: 实际发送短信
	// 这里只打印验证码，实际项目中应该调用短信服务发送
	println("Phone Code:", code)

	return &SendCodeResponse{
		ExpireSeconds: CodeExpire,
	}, nil
}

// RefreshToken 刷新Token
func (s *UserService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 1. 解析Refresh Token
	claims, err := s.jwtUtil.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, global.RefreshTokenInvalid
	}

	// 2. 检查Token类型
	if claims.Type != "refresh" {
		return nil, global.RefreshTokenInvalid
	}

	// 3. 生成新的Access Token
	accessToken, err := s.jwtUtil.GenerateToken(claims.UserID, claims.Username, "access")
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   config.AppConfig.JWT.AccessExpire,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*domain.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, err
	}

	return user.ToVO(), nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, userID int64, req *UpdateUserInfoRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.UserNotFound
		}
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.UserNotFound
		}
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return global.PasswordError
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// SendResetPasswordCode 发送重置密码验证码
func (s *UserService) SendResetPasswordCode(ctx context.Context, req *SendResetPasswordCodeRequest) (*SendCodeResponse, error) {
	// 1. 检查邮箱是否存在
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.UserNotFound
		}
		return nil, err
	}

	// 2. 检查发送频率
	err = s.userRedis.SetSendCodeRateLimit(ctx, "reset:"+req.Email, time.Duration(SendCodeInterval)*time.Second)
	if err != nil {
		return nil, global.CodeSendTooFast
	}

	// 3. 生成验证码
	code := utils.GenerateRandomCode(6)

	// 4. 存储验证码
	err = s.userRedis.SetEmailCode(ctx, "reset:"+req.Email, code, time.Duration(CodeExpire)*time.Second)
	if err != nil {
		return nil, err
	}

	// TODO: 实际发送邮件
	_ = user
	println("Reset Password Code:", code)

	return &SendCodeResponse{
		ExpireSeconds: CodeExpire,
	}, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// 1. 校验验证码
	code, err := s.userRedis.GetEmailCode(ctx, "reset:"+req.Email)
	if err != nil {
		return global.CodeExpired
	}
	if code != req.Code {
		return global.CodeError
	}

	// 2. 获取用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.UserNotFound
		}
		return err
	}

	// 3. 加密新密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 4. 更新密码
	if err := s.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return err
	}

	// 5. 删除验证码
	s.userRedis.DelEmailCode(ctx, "reset:"+req.Email)

	return nil
}

// Logout 登出
func (s *UserService) Logout(ctx context.Context, token string) error {
	// 将Token加入黑名单
	// 注意：这里简化处理，实际应该解析Token获取过期时间
	return s.userRedis.SetTokenBlacklist(ctx, token, time.Duration(config.AppConfig.JWT.AccessExpire)*time.Second)
}
