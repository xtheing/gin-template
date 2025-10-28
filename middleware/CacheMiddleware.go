package middleware

import (
	"net/http"
	"strings"
	"time"

	"theing/gin-template/common"

	"github.com/gin-gonic/gin"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Duration time.Duration // 缓存时间
	Key      string        // 缓存键
	Enabled  bool          // 是否启用缓存
}

// CacheMiddleware 缓存中间件
func CacheMiddleware(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		// 生成缓存键
		cacheKey := generateCacheKey(c, config.Key)

		// 尝试从缓存获取响应
		cacheHelper := common.NewCacheHelper()
		var cachedResponse map[string]interface{}
		err := cacheHelper.GetJSON(c.Request.Context(), cacheKey, &cachedResponse)
		if err == nil {
			// 缓存命中，返回缓存的响应
			c.JSON(http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		// 缓存未命中，继续处理请求
		c.Next()

		// 只缓存成功的响应
		if c.Writer.Status() == http.StatusOK {
			// 获取响应数据
			responseData, exists := c.Get("response_data")
			if exists {
				// 设置缓存
				cacheHelper.SetJSON(c.Request.Context(), cacheKey, responseData, config.Duration)
			}
		}
	}
}

// CacheKeyGenerator 自定义缓存键生成器
func CacheKeyGenerator(keyGenerator func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key := keyGenerator(c); key != "" {
			c.Set("cache_key", key)
		}
		c.Next()
	}
}

// CacheResponse 缓存响应数据的辅助函数
func CacheResponse(c *gin.Context, data interface{}) {
	c.Set("response_data", data)
}

// generateCacheKey 生成缓存键
func generateCacheKey(c *gin.Context, baseKey string) string {
	// 基础键
	key := baseKey

	// 添加路径信息
	if path := c.Request.URL.Path; path != "" {
		key += ":" + strings.ReplaceAll(path, "/", "_")
	}

	// 添加查询参数
	if query := c.Request.URL.RawQuery; query != "" {
		key += ":" + strings.ReplaceAll(query, "&", "_")
	}

	// 添加用户信息（如果有）
	if userID, exists := c.Get("user_id"); exists {
		key += ":user_" + userID.(string)
	}

	return key
}

// GetUserCacheKey 获取用户相关的缓存键
func GetUserCacheKey(userID interface{}, suffix string) string {
	return "user:" + userID.(string) + ":" + suffix
}

// GetOptionCacheKey 获取选项数据的缓存键
func GetOptionCacheKey(optionType string, suffix string) string {
	return "options:" + optionType + ":" + suffix
}

// InvalidateUserCache 清除用户相关缓存
func InvalidateUserCache(c *gin.Context, userID interface{}) error {
	cacheHelper := common.NewCacheHelper()
	pattern := "user:" + userID.(string) + ":*"
	return cacheHelper.DeletePattern(c.Request.Context(), pattern)
}

// InvalidateOptionCache 清除选项相关缓存
func InvalidateOptionCache(c *gin.Context, optionType string) error {
	cacheHelper := common.NewCacheHelper()
	pattern := "options:" + optionType + ":*"
	return cacheHelper.DeletePattern(c.Request.Context(), pattern)
}

// CacheableResponse 可缓存的响应结构
type CacheableResponse struct {
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	ExpiresIn int64       `json:"expires_in"`
}

// WrapCacheableResponse 包装可缓存的响应
func WrapCacheableResponse(data interface{}, expiresIn time.Duration) CacheableResponse {
	return CacheableResponse{
		Data:      data,
		Timestamp: time.Now().Unix(),
		ExpiresIn: int64(expiresIn.Seconds()),
	}
}
