package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
)

// MD5 MD5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 SHA256加密
func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateRandomKey 生成随机密钥
func GenerateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// HashPassword 密码哈希
func HashPassword(password string) string {
	return MD5(password + "salt") // 生产环境建议使用bcrypt
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	// 实际项目中使用分布式ID生成器
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Int64ToString int64转string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 string转int64
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
