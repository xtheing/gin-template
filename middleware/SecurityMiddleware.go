package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"theing/gin-template/common"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableCORS          bool          `json:"enable_cors"`
	EnableRateLimit     bool          `json:"enable_rate_limit"`
	MaxRequestsPerMin   int           `json:"max_requests_per_min"`
	EnableCSRF          bool          `json:"enable_csrf"`
	EnableSecureHeaders bool          `json:"enable_secure_headers"`
	TrustedProxies      []string      `json:"trusted_proxies"`
	RequestTimeout      time.Duration `json:"request_timeout"`
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头部
		if config.EnableSecureHeaders {
			setSecurityHeaders(c)
		}

		// CORS 处理
		if config.EnableCORS {
			handleCORS(c)
		}

		// 请求超时处理
		if config.RequestTimeout > 0 {
			ctx, cancel := context.WithTimeout(c.Request.Context(), config.RequestTimeout)
			defer cancel()
			c.Request = c.Request.WithContext(ctx)
		}

		// IP 白名单检查
		if len(config.TrustedProxies) > 0 {
			clientIP := c.ClientIP()
			if !isTrustedIP(clientIP, config.TrustedProxies) {
				appErr := common.NewAppError(common.CodeUnauthorized, "IP地址不被信任", "")
				c.JSON(http.StatusForbidden, appErr)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	// 简化的限流实现，生产环境建议使用 Redis 或专门的服务
	clients := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// 清理过期记录
		if requests, exists := clients[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			clients[clientIP] = validRequests
		}

		// 检查是否超过限制
		if requests, exists := clients[clientIP]; exists && len(requests) >= maxRequests {
			appErr := common.NewAppError(common.CodeTooManyRequests, "请求过于频繁，请稍后再试", "")
			c.JSON(http.StatusTooManyRequests, appErr)
			c.Abort()
			return
		}

		// 记录当前请求
		clients[clientIP] = append(clients[clientIP], now)

		c.Next()
	}
}

// InputValidationMiddleware 输入验证中间件
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 SQL 注入
		if containsSQLInjection(c.Request.URL.RawQuery) {
			appErr := common.NewAppError(common.CodeBadRequest, "请求参数包含非法字符", "")
			c.JSON(http.StatusBadRequest, appErr)
			c.Abort()
			return
		}

		// 检查 XSS
		if containsXSS(c.Request.URL.RawQuery) {
			appErr := common.NewAppError(common.CodeBadRequest, "请求参数包含非法字符", "")
			c.JSON(http.StatusBadRequest, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CSRFProtectionMiddleware CSRF 保护中间件
func CSRFProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于 GET 请求跳过 CSRF 检查
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// 检查 CSRF Token
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("csrf_token")
		}

		if token == "" || !validateCSRFToken(token, c) {
			appErr := common.NewAppError(common.CodeUnauthorized, "CSRF token 无效", "")
			c.JSON(http.StatusForbidden, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// setSecurityHeaders 设置安全头部
func setSecurityHeaders(c *gin.Context) {
	// 防止点击劫持
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")

	// XSS 保护
	c.Header("X-XSS-Protection", "1; mode=block")

	// 强制 HTTPS（生产环境）
	if gin.Mode() == gin.ReleaseMode {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	// 内容安全策略
	c.Header("Content-Security-Policy", "default-src 'self'")

	// 隐藏服务器信息
	c.Header("Server", "")

	// 限制跨域请求
	c.Header("Access-Control-Allow-Credentials", "false")
}

// handleCORS 处理 CORS
func handleCORS(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")

	// 在开发环境允许所有源
	if gin.Mode() == gin.DebugMode {
		c.Header("Access-Control-Allow-Origin", "*")
	} else if origin != "" {
		// 生产环境应该检查白名单
		c.Header("Access-Control-Allow-Origin", origin)
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
	c.Header("Access-Control-Max-Age", "86400")

	// 处理预检请求
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}

// isTrustedIP 检查IP是否在信任列表中
func isTrustedIP(ip string, trustedProxies []string) bool {
	for _, trusted := range trustedProxies {
		if ip == trusted {
			return true
		}
	}
	return false
}

// containsSQLInjection 检查SQL注入
func containsSQLInjection(input string) bool {
	if input == "" {
		return false
	}

	// 常见的 SQL 注入模式
	patterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(or|and)\s+\d+\s*=\s*\d+`,
		`(?i)(['"];|\/\*|\*\/|--|#|\/\*)`,
		`(?i)(script|javascript|vbscript|onload|onerror)`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}

	return false
}

// containsXSS 检查XSS攻击
func containsXSS(input string) bool {
	if input == "" {
		return false
	}

	// 常见的 XSS 模式
	patterns := []string{
		`(?i)(<script|</script|<iframe|</iframe|<object|</object)`,
		`(?i)(javascript:|vbscript:|data:)`,
		`(?i)(onload|onerror|onclick|onmouseover)`,
		`(?i)(eval\(|alert\(|confirm\(|prompt\()`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return true
		}
	}

	return false
}

// validateCSRFToken 验证CSRF Token
func validateCSRFToken(token string, c *gin.Context) bool {
	// 这里应该实现真正的 CSRF token 验证逻辑
	// 简化实现，生产环境需要更复杂的逻辑
	if token == "" {
		return false
	}

	// 可以从 session 中获取期望的 token
	// 或者使用双重提交 Cookie 等方法
	return len(token) > 10 // 简化的长度检查
}

// generateCSRFToken 生成CSRF Token
func generateCSRFToken(c *gin.Context) string {
	// 简化的 token 生成，生产环境应该使用更安全的方法
	return "csrf_token_" + time.Now().Format("20060102150405")
}

// SanitizeInput 清理输入
func SanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	return input
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	// 简化的邮箱验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return common.NewAppError(common.CodeValidationFailed, "密码长度至少8位", "")
	}

	if len(password) > 128 {
		return common.NewAppError(common.CodeValidationFailed, "密码长度不能超过128位", "")
	}

	// 检查是否包含数字
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	// 检查是否包含大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 检查是否包含小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 检查是否包含特殊字符
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	count := 0
	if hasNumber {
		count++
	}
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasSpecial {
		count++
	}

	if count < 3 {
		return common.NewAppError(common.CodeValidationFailed, "密码必须包含数字、大小写字母、特殊字符中的至少3种", "")
	}

	return nil
}

// RateLimiter 简单的内存限流器
type RateLimiter struct {
	requests map[string][]time.Time
	window   time.Duration
	max      int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		max:      max,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	now := time.Now()

	// 清理过期请求
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// 检查是否超过限制
	if requests, exists := rl.requests[key]; exists && len(requests) >= rl.max {
		return false
	}

	// 记录当前请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}
