package domain

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Status    int            `gorm:"type:tinyint;default:1" json:"status"` // 1:正常 2:禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserInfo 用户信息VO
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
}

// ToVO 转换为VO
func (u *User) ToVO() *UserInfo {
	return &UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}
}
