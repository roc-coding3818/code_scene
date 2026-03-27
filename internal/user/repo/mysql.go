package user

import (
	"code_scene/global"
	"code_scene/internal/domain"

	"gorm.io/gorm"
)

// UserRepo 用户仓库
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo() *UserRepo {
	return &UserRepo{db: global.DB}
}

// Create 创建用户
func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepo) GetByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// UpdatePassword 更新密码
func (r *UserRepo) UpdatePassword(id int64, password string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateStatus 更新用户状态
func (r *UserRepo) UpdateStatus(id int64, status int) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(id int64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepo) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByPhone 检查手机号是否存在
func (r *UserRepo) ExistsByPhone(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}
