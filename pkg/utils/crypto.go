package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomCode 生成指定位数的随机数字验证码
func GenerateRandomCode(length int) string {
	code := make([]byte, length)
	_, err := rand.Read(code)
	if err != nil {
		return ""
	}

	// 转换为数字
	for i := range code {
		code[i] = byte(int(code[i]) % 10)
	}

	return string(code)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// RSAEncrypt RSA加密（示例，实际使用需配置公钥）
func RSAEncrypt(data, pubKey []byte) ([]byte, error) {
	// 这里简化处理，实际应使用真实的RSA公钥
	// 演示用，实际建议使用 crypto/rsa
	return data, nil
}

// RSADecrypt RSA解密
func RSADecrypt(data, privKey []byte) ([]byte, error) {
	return data, nil
}

// SHA256Hash SHA256哈希
func SHA256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
