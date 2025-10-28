package middleware

import (
	"runtime"
	"time"

	"theing/gin-template/common"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware 性能监控中间件
func MetricsMiddleware() gin.HandlerFunc {
	metrics := common.GetMetrics()
	
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		
		// 获取请求大小
		requestSize := c.Request.ContentLength
		
		// 处理请求
		c.Next()
		
		// 计算请求持续时间
		duration := time.Since(start)
		
		// 获取响应状态码和大小
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		
		// 记录 HTTP 请求指标
		metrics.RecordHttpRequest(
			method,
			common.SanitizeEndpoint(path),
			common.GetStatusCodeGroup(statusCode),
			duration,
			requestSize,
			int64(responseSize),
		)
		
		// 定期更新系统指标
		updateSystemMetrics(metrics)
	}
}

// MetricsHandler Prometheus 指标处理器
func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

// updateSystemMetrics 更新系统指标
func updateSystemMetrics(metrics *common.Metrics) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// 更新内存使用
	metrics.UpdateSystemMemory(m.HeapAlloc, m.StackInuse, m.Sys)
	
	// 更新协程数量
	metrics.UpdateGoroutineCount(runtime.NumGoroutine())
	
	// 更新数据库连接数（如果数据库已初始化）
	if common.DB != nil {
		sqlDB, err := common.DB.DB()
		if err == nil {
			stats := sqlDB.Stats()
			metrics.UpdateDatabaseConnections(stats.Idle, stats.InUse, stats.OpenConnections)
		}
	}
}

// DatabaseMetricsMiddleware 数据库查询监控中间件
func DatabaseMetricsMiddleware() gin.HandlerFunc {
	metrics := common.GetMetrics()
	
	return func(c *gin.Context) {
		// 在请求上下文中存储指标记录器
		c.Set("metrics", metrics)
		c.Next()
	}
}

// RecordDatabaseQuery 记录数据库查询（供数据库操作使用）
func RecordDatabaseQuery(c *gin.Context, operation, table string, duration time.Duration, err error) {
	if metrics, exists := c.Get("metrics"); exists {
		m := metrics.(*common.Metrics)
		status := "success"
		if err != nil {
			status = "error"
		}
		m.RecordDatabaseQuery(operation, table, status, duration)
	}
}

// RecordCacheMetrics 记录缓存指标
func RecordCacheMetrics(c *gin.Context, operation, cacheType, keyPrefix string, hit bool) {
	metrics := common.GetMetrics()
	
	if hit {
		metrics.RecordCacheHit(cacheType, keyPrefix)
	} else {
		metrics.RecordCacheMiss(cacheType, keyPrefix)
	}
	
	metrics.RecordCacheOperation(operation, cacheType, "success")
}

// RecordJWTMetrics 记录 JWT 相关指标
func RecordJWTMetrics(c *gin.Context, tokenType string, validated bool, errorType string) {
	metrics := common.GetMetrics()
	
	if validated {
		metrics.RecordJWTTokenValidated("success")
		if tokenType == "issued" {
			metrics.RecordJWTTokenIssued("user")
		}
	} else {
		metrics.RecordJWTTokenValidated("error")
		if errorType != "" {
			metrics.RecordJWTValidationError(errorType)
		}
	}
}

// RecordBusinessMetrics 记录业务指标
func RecordBusinessMetrics(c *gin.Context, action, entityType, status string) {
	metrics := common.GetMetrics()
	
	switch action {
	case "register":
		metrics.RecordUserRegistration(status)
	case "login":
		metrics.RecordUserLogin(status, entityType)
	}
}

// PerformanceMetrics 性能指标结构
type PerformanceMetrics struct {
	RequestCount    int64         `json:"request_count"`
	AverageLatency  time.Duration `json:"average_latency"`
	ErrorRate       float64       `json:"error_rate"`
	CacheHitRate    float64       `json:"cache_hit_rate"`
	DatabaseQueries int64         `json:"database_queries"`
	MemoryUsage     uint64        `json:"memory_usage"`
	GoroutineCount  int           `json:"goroutine_count"`
}

// GetPerformanceMetrics 获取性能指标摘要
func GetPerformanceMetrics() PerformanceMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// 这里简化处理，实际应用中应该从指标中获取真实数据
	return PerformanceMetrics{
		RequestCount:    0, // 应该从 HTTP 请求计数器获取
		AverageLatency:  0, // 应该从 HTTP 延迟直方图计算
		ErrorRate:       0, // 应该从错误计数器计算
		CacheHitRate:    0, // 应该从缓存指标计算
		DatabaseQueries: 0, // 应该从数据库查询计数器获取
		MemoryUsage:     m.HeapAlloc,
		GoroutineCount:  runtime.NumGoroutine(),
	}
}
