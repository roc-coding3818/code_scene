package jwt

import (
	"code_scene/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// JWT JWT工具
type JWT struct {
	secret        []byte
	accessExpire int64
	refreshExpire int64
}

// NewJWT 创建JWT实例
func NewJWT(cfg *config.JWTConfig) *JWT {
	return &JWT{
		secret:        []byte(cfg.Secret),
		accessExpire:  cfg.AccessExpire,
		refreshExpire: cfg.RefreshExpire,
	}
}

// GenerateToken 生成Token
func (j *JWT) GenerateToken(userID int64, username, tokenType string) (string, error) {
	var expire time.Duration
	if tokenType == "refresh" {
		expire = time.Duration(j.refreshExpire) * time.Second
	} else {
		expire = time.Duration(j.accessExpire) * time.Second
	}

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "code_scene",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// GenerateAccessAndRefreshToken 生成Access和Refresh Token
func (j *JWT) GenerateAccessAndRefreshToken(userID int64, username string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateToken(userID, username, "access")
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateToken(userID, username, "refresh")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证Token是否有效
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查是否是 access token
	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
