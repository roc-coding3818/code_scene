package user

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名 3-20位
	Password string `json:"password" binding:"required,min=6,max=20"` // 密码 6-20位
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Code     string `json:"code" binding:"required,len=6"`             // 邮箱验证码
}

// PhoneRegisterRequest 手机号注册请求
type PhoneRegisterRequest struct {
	Phone    string `json:"phone" binding:"required,len=11"`   // 手机号
	Password string `json:"password" binding:"required,min=6"` // 密码
	Code     string `json:"code" binding:"required,len=6"`      // 短信验证码
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名/邮箱/手机号
	Password string `json:"password" binding:"required"` // 密码
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone string `json:"phone" binding:"required,len=11"` // 手机号
	Code  string `json:"code" binding:"required,len=6"`   // 短信验证码
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
	Type  string `json:"type" binding:"required"`        // 验证码类型: register, login, reset_password
}

// SendPhoneCodeRequest 发送手机验证码请求
type SendPhoneCodeRequest struct {
	Phone string `json:"phone" binding:"required,len=11"` // 手机号
	Type  string `json:"type" binding:"required"`         // 验证码类型
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar"`
}

// SendResetPasswordCodeRequest 发送重置密码验证码请求
type SendResetPasswordCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
