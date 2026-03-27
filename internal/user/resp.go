package user

import "code_scene/internal/domain"

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int64            `json:"expires_in"` // AccessToken过期时间(秒)
	UserInfo     *domain.UserInfo `json:"user_info"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	ExpireSeconds int `json:"expire_seconds"` // 验证码有效期(秒)
}
