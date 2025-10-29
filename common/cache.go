package common

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// CacheClient 缓存客户端接口
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Flush(ctx context.Context) error
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string
}

var (
	// Cache 全局缓存客户端
	Cache CacheClient
	// CacheConfig 缓存配置
	cacheConfig CacheConfig
	// ErrCacheNotFound 缓存未找到错误
	ErrCacheNotFound = fmt.Errorf("缓存不存在")
)

// InitCache 初始化缓存连接
func InitCache() error {
	// 从配置文件读取缓存配置
	cacheConfig = CacheConfig{
		Host:     viper.GetString("cache_host"),
		Port:     viper.GetInt("cache_port"),
		Password: viper.GetString("cache_password"),
		DB:       viper.GetInt("cache_db"),
		Prefix:   viper.GetString("cache_prefix"),
	}

	// 设置默认值
	if cacheConfig.Host == "" {
		cacheConfig.Host = "localhost"
	}
	if cacheConfig.Port == 0 {
		cacheConfig.Port = 6379
	}
	if cacheConfig.Prefix == "" {
		cacheConfig.Prefix = "gin_template:"
	}

	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password: cacheConfig.Password,
		DB:       cacheConfig.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("缓存连接失败: %v", err)
	}

	Cache = &RedisCache{client: rdb}
	fmt.Printf("缓存连接成功: %s:%d\n", cacheConfig.Host, cacheConfig.Port)
	return nil
}

// Get 获取缓存值
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	fullKey := r.getFullKey(key)
	return r.client.Get(ctx, fullKey).Result()
}

// Set 设置缓存值
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fullKey := r.getFullKey(key)
	return r.client.Set(ctx, fullKey, value, expiration).Err()
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.getFullKey(key)
	return r.client.Del(ctx, fullKey).Err()
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.getFullKey(key)
	result, err := r.client.Exists(ctx, fullKey).Result()
	return result > 0, err
}

// Flush 清空所有缓存
func (r *RedisCache) Flush(ctx context.Context) error {
	// 只清空当前应用的缓存
	pattern := r.getFullKey("*")
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}

// getFullKey 获取完整的缓存键
func (r *RedisCache) getFullKey(key string) string {
	return cacheConfig.Prefix + key
}

// GetCacheConfig 获取缓存配置
func GetCacheConfig() CacheConfig {
	return cacheConfig
}

// CacheHelper 缓存辅助函数
type CacheHelper struct {
	cache CacheClient
}

// NewCacheHelper 创建缓存辅助器
func NewCacheHelper() *CacheHelper {
	return &CacheHelper{cache: Cache}
}

// GetJSON 获取JSON格式的缓存数据
func (h *CacheHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	value, err := h.cache.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("缓存不存在")
		}
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

// SetJSON 设置JSON格式的缓存数据
func (h *CacheHelper) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return h.cache.Set(ctx, key, string(data), expiration)
}

// GetOrSet 获取缓存或设置新值
func (h *CacheHelper) GetOrSet(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 尝试从缓存获取
	var result interface{}
	err := h.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	// 缓存不存在，执行函数获取数据
	result, err = fn()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := h.SetJSON(ctx, key, result, expiration); err != nil {
		// 缓存设置失败不影响主流程，只记录日志
		fmt.Printf("缓存设置失败: %v\n", err)
	}

	return result, nil
}

// DeletePattern 根据模式删除缓存
func (h *CacheHelper) DeletePattern(ctx context.Context, pattern string) error {
	if redisCache, ok := h.cache.(*RedisCache); ok {
		fullPattern := cacheConfig.Prefix + pattern
		keys, err := redisCache.client.Keys(ctx, fullPattern).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			return redisCache.client.Del(ctx, keys...).Err()
		}
	}
	return fmt.Errorf("当前缓存类型不支持模式删除")
}

// GetUserCacheKey 生成用户缓存键
func GetUserCacheKey(userID, suffix string) string {
	return fmt.Sprintf("user:%s:%s", userID, suffix)
}

// GetOptionCacheKey 生成选项缓存键
func GetOptionCacheKey(optionType, suffix string) string {
	return fmt.Sprintf("options:%s:%s", optionType, suffix)
}
