package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 监控指标结构
type Metrics struct {
	// HTTP 请求指标
	HttpRequestsTotal    *prometheus.CounterVec
	HttpRequestDuration  *prometheus.HistogramVec
	HttpRequestSize      *prometheus.HistogramVec
	HttpResponseSize     *prometheus.HistogramVec

	// 数据库指标
	DatabaseConnections  *prometheus.GaugeVec
	DatabaseQueryTotal   *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec

	// 缓存指标
	CacheHitsTotal       *prometheus.CounterVec
	CacheMissesTotal     *prometheus.CounterVec
	CacheOperationsTotal *prometheus.CounterVec

	// JWT 指标
	JWTTokensIssued      *prometheus.CounterVec
	JWTTokensValidated   *prometheus.CounterVec
	JWTValidationErrors  *prometheus.CounterVec

	// 系统指标
	SystemErrorsTotal     *prometheus.CounterVec
	SystemPanicTotal      *prometheus.CounterVec
	SystemMemoryUsage     *prometheus.GaugeVec
	SystemGoroutineCount  *prometheus.GaugeVec

	// 业务指标
	UserRegistrations     *prometheus.CounterVec
	UserLogins           *prometheus.CounterVec
	ActiveSessions       *prometheus.GaugeVec
}

var (
	// 全局指标实例
	metrics *Metrics

	// 自定义标签
	constLabels = prometheus.Labels{
		"service": "gin-template",
		"version": "2.0.0",
	}
)

// InitMetrics 初始化监控指标
func InitMetrics() *Metrics {
	metrics = &Metrics{
		// HTTP 请求指标
		HttpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_requests_total",
				Help:        "Total number of HTTP requests",
				ConstLabels: constLabels,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HttpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_request_duration_seconds",
				Help:        "HTTP request duration in seconds",
				Buckets:     []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5, 10},
				ConstLabels: constLabels,
			},
			[]string{"method", "endpoint"},
		),
		HttpRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_request_size_bytes",
				Help:        "HTTP request size in bytes",
				Buckets:     prometheus.ExponentialBuckets(100, 10, 8),
				ConstLabels: constLabels,
			},
			[]string{"method", "endpoint"},
		),
		HttpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_response_size_bytes",
				Help:        "HTTP response size in bytes",
				Buckets:     prometheus.ExponentialBuckets(100, 10, 8),
				ConstLabels: constLabels,
			},
			[]string{"method", "endpoint"},
		),

		// 数据库指标
		DatabaseConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        "database_connections",
				Help:        "Number of database connections",
				ConstLabels: constLabels,
			},
			[]string{"state"}, // idle, in_use, open
		),
		DatabaseQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "database_queries_total",
				Help:        "Total number of database queries",
				ConstLabels: constLabels,
			},
			[]string{"operation", "table", "status"}, // select, insert, update, delete
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "database_query_duration_seconds",
				Help:        "Database query duration in seconds",
				Buckets:     []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
				ConstLabels: constLabels,
			},
			[]string{"operation", "table"},
		),

		// 缓存指标
		CacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "cache_hits_total",
				Help:        "Total number of cache hits",
				ConstLabels: constLabels,
			},
			[]string{"cache_type", "key_prefix"},
		),
		CacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "cache_misses_total",
				Help:        "Total number of cache misses",
				ConstLabels: constLabels,
			},
			[]string{"cache_type", "key_prefix"},
		),
		CacheOperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "cache_operations_total",
				Help:        "Total number of cache operations",
				ConstLabels: constLabels,
			},
			[]string{"operation", "cache_type", "status"}, // get, set, delete
		),

		// JWT 指标
		JWTTokensIssued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "jwt_tokens_issued_total",
				Help:        "Total number of JWT tokens issued",
				ConstLabels: constLabels,
			},
			[]string{"user_type"},
		),
		JWTTokensValidated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "jwt_tokens_validated_total",
				Help:        "Total number of JWT tokens validated",
				ConstLabels: constLabels,
			},
			[]string{"status"},
		),
		JWTValidationErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "jwt_validation_errors_total",
				Help:        "Total number of JWT validation errors",
				ConstLabels: constLabels,
			},
			[]string{"error_type"},
		),

		// 系统指标
		SystemErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "system_errors_total",
				Help:        "Total number of system errors",
				ConstLabels: constLabels,
			},
			[]string{"error_type", "component"},
		),
		SystemPanicTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "system_panics_total",
				Help:        "Total number of system panics",
				ConstLabels: constLabels,
			},
			[]string{"component"},
		),
		SystemMemoryUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        "system_memory_usage_bytes",
				Help:        "System memory usage in bytes",
				ConstLabels: constLabels,
			},
			[]string{"type"}, // heap, stack, sys
		),
		SystemGoroutineCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        "system_goroutines",
				Help:        "Number of goroutines",
				ConstLabels: constLabels,
			},
			[]string{},
		),

		// 业务指标
		UserRegistrations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "user_registrations_total",
				Help:        "Total number of user registrations",
				ConstLabels: constLabels,
			},
			[]string{"status"},
		),
		UserLogins: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "user_logins_total",
				Help:        "Total number of user logins",
				ConstLabels: constLabels,
			},
			[]string{"status", "user_type"},
		),
		ActiveSessions: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        "active_sessions",
				Help:        "Number of active user sessions",
				ConstLabels: constLabels,
			},
			[]string{"user_type"},
		),
	}

	return metrics
}

// GetMetrics 获取全局指标实例
func GetMetrics() *Metrics {
	if metrics == nil {
		return InitMetrics()
	}
	return metrics
}

// RecordHttpRequest 记录 HTTP 请求指标
func (m *Metrics) RecordHttpRequest(method, endpoint, statusCode string, duration time.Duration, requestSize, responseSize int64) {
	m.HttpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.HttpRequestSize.WithLabelValues(method, endpoint).Observe(float64(requestSize))
	m.HttpResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

// RecordDatabaseQuery 记录数据库查询指标
func (m *Metrics) RecordDatabaseQuery(operation, table, status string, duration time.Duration) {
	m.DatabaseQueryTotal.WithLabelValues(operation, table, status).Inc()
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// UpdateDatabaseConnections 更新数据库连接数
func (m *Metrics) UpdateDatabaseConnections(idle, inUse, open int) {
	m.DatabaseConnections.WithLabelValues("idle").Set(float64(idle))
	m.DatabaseConnections.WithLabelValues("in_use").Set(float64(inUse))
	m.DatabaseConnections.WithLabelValues("open").Set(float64(open))
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit(cacheType, keyPrefix string) {
	m.CacheHitsTotal.WithLabelValues(cacheType, keyPrefix).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss(cacheType, keyPrefix string) {
	m.CacheMissesTotal.WithLabelValues(cacheType, keyPrefix).Inc()
}

// RecordCacheOperation 记录缓存操作
func (m *Metrics) RecordCacheOperation(operation, cacheType, status string) {
	m.CacheOperationsTotal.WithLabelValues(operation, cacheType, status).Inc()
}

// RecordJWTTokenIssued 记录 JWT 令牌签发
func (m *Metrics) RecordJWTTokenIssued(userType string) {
	m.JWTTokensIssued.WithLabelValues(userType).Inc()
}

// RecordJWTTokenValidated 记录 JWT 令牌验证
func (m *Metrics) RecordJWTTokenValidated(status string) {
	m.JWTTokensValidated.WithLabelValues(status).Inc()
}

// RecordJWTValidationError 记录 JWT 验证错误
func (m *Metrics) RecordJWTValidationError(errorType string) {
	m.JWTValidationErrors.WithLabelValues(errorType).Inc()
}

// RecordSystemError 记录系统错误
func (m *Metrics) RecordSystemError(errorType, component string) {
	m.SystemErrorsTotal.WithLabelValues(errorType, component).Inc()
}

// RecordSystemPanic 记录系统恐慌
func (m *Metrics) RecordSystemPanic(component string) {
	m.SystemPanicTotal.WithLabelValues(component).Inc()
}

// UpdateSystemMemory 更新系统内存使用
func (m *Metrics) UpdateSystemMemory(heap, stack, sys uint64) {
	m.SystemMemoryUsage.WithLabelValues("heap").Set(float64(heap))
	m.SystemMemoryUsage.WithLabelValues("stack").Set(float64(stack))
	m.SystemMemoryUsage.WithLabelValues("sys").Set(float64(sys))
}

// UpdateGoroutineCount 更新协程数量
func (m *Metrics) UpdateGoroutineCount(count int) {
	m.SystemGoroutineCount.WithLabelValues().Set(float64(count))
}

// RecordUserRegistration 记录用户注册
func (m *Metrics) RecordUserRegistration(status string) {
	m.UserRegistrations.WithLabelValues(status).Inc()
}

// RecordUserLogin 记录用户登录
func (m *Metrics) RecordUserLogin(status, userType string) {
	m.UserLogins.WithLabelValues(status, userType).Inc()
}

// UpdateActiveSessions 更新活跃会话数
func (m *Metrics) UpdateActiveSessions(userType string, count int) {
	m.ActiveSessions.WithLabelValues(userType).Set(float64(count))
}

// GetStatusCodeGroup 获取状态码分组
func GetStatusCodeGroup(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}

// SanitizeEndpoint 清理端点名称
func SanitizeEndpoint(endpoint string) string {
	// 将路径参数替换为占位符
	// 例如: /api/users/123 -> /api/users/:id
	// 简单实现，可以根据需要扩展
	return endpoint
}
