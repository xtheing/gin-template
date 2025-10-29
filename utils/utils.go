package utils

// 各种工具包

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// 生成随机字符串的功能，传入一个数字的长度，生成相应长度的随机字符串。
func RandomString(n int) string {
	var letters = []byte("asdfkjj;asdf;lasqpoewitupoizx,cvnb/m")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// PasswordStrength 密码强度类型
type PasswordStrength int

const (
	PasswordWeak   PasswordStrength = iota // 弱密码
	PasswordMedium                         // 中等密码
	PasswordStrong                         // 强密码
)

// PasswordValidation 密码验证结果
type PasswordValidation struct {
	IsValid     bool             // 是否有效
	Strength    PasswordStrength // 密码强度
	Score       int              // 密码评分 (0-100)
	Errors      []string         // 错误信息
	Suggestions []string         // 建议
}

// ValidatePassword 验证密码复杂度
func ValidatePassword(password string) PasswordValidation {
	var validation PasswordValidation
	var errors []string
	var suggestions []string
	score := 0

	// 基本长度检查
	if len(password) < 6 {
		errors = append(errors, "密码长度不能少于6位")
	} else if len(password) >= 6 {
		score += 20
	}
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 10
	}

	// 包含小写字母
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score += 15
	} else {
		errors = append(errors, "密码必须包含小写字母")
		suggestions = append(suggestions, "添加小写字母可以增加密码强度")
	}

	// 包含大写字母
	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score += 15
	} else {
		suggestions = append(suggestions, "添加大写字母可以增加密码强度")
	}

	// 包含数字
	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score += 15
	} else {
		errors = append(errors, "密码必须包含数字")
		suggestions = append(suggestions, "添加数字可以增加密码强度")
	}

	// 包含特殊字符
	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password); matched {
		score += 15
	} else {
		suggestions = append(suggestions, "添加特殊字符可以显著增加密码强度")
	}

	// 常见弱密码检查
	weakPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
	}
	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if strings.Contains(lowerPassword, weak) {
			score -= 30
			errors = append(errors, "密码包含常见弱密码模式")
			break
		}
	}

	// 重复字符检查
	if hasRepeatingChars(password, 3) {
		score -= 10
		suggestions = append(suggestions, "避免使用连续重复的字符")
	}

	// 确保分数在合理范围内
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	validation.Score = score
	validation.Errors = errors
	validation.Suggestions = suggestions
	validation.IsValid = len(errors) == 0 && score >= 60

	// 确定密码强度
	if score < 40 {
		validation.Strength = PasswordWeak
	} else if score < 70 {
		validation.Strength = PasswordMedium
	} else {
		validation.Strength = PasswordStrong
	}

	return validation
}

// hasRepeatingChars 检查是否有连续重复字符
func hasRepeatingChars(s string, maxRepeats int) bool {
	if len(s) < maxRepeats {
		return false
	}

	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= maxRepeats {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// GetPasswordStrengthText 获取密码强度文本描述
func GetPasswordStrengthText(strength PasswordStrength) string {
	switch strength {
	case PasswordWeak:
		return "弱"
	case PasswordMedium:
		return "中等"
	case PasswordStrong:
		return "强"
	default:
		return "未知"
	}
}
