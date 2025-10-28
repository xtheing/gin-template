package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecureKey 生成安全的随机密钥
func GenerateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("生成随机密钥失败: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateJWTKey 生成适用于 JWT 的安全密钥（32字节）
func GenerateJWTKey() (string, error) {
	return GenerateSecureKey(32)
}
