package common

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`    // 状态：healthy, unhealthy
	Timestamp int64                  `json:"timestamp"` // 时间戳
	Database  DatabaseHealthStatus     `json:"database"`  // 数据库状态
	Services  map[string]interface{}  `json:"services"`  // 其他服务状态
}

// DatabaseHealthStatus 数据库健康状态
type DatabaseHealthStatus struct {
	Status     string        `json:"status"`     // 状态
	Connection string        `json:"connection"` // 连接信息
	Latency   time.Duration `json:"latency"`   // 延迟
	Error      string        `json:"error,omitempty"` // 错误信息
}

// DatabaseHealthChecker 数据库健康检查器
type DatabaseHealthChecker struct {
	db      *gorm.DB
	latency time.Duration // 存储延迟
}

// NewDatabaseHealthChecker 创建数据库健康检查器
func NewDatabaseHealthChecker(db *gorm.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

// CheckHealth 检查数据库健康状态
func (d *DatabaseHealthChecker) CheckHealth() error {
	if d.db == nil {
		return fmt.Errorf("数据库连接为空")
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 测试数据库连接
	start := time.Now()
	err = sqlDB.Ping()
	latency := time.Since(start)

	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 记录延迟
	d.latency = latency
	return nil
}

// GetName 获取检查器名称
func (d *DatabaseHealthChecker) GetName() string {
	return "database"
}

// CheckSystemHealth 检查系统健康状态
func CheckSystemHealth() HealthStatus {
	status := HealthStatus{
		Timestamp: time.Now().Unix(),
		Services:  make(map[string]interface{}),
	}

	// 检查数据库
	dbStatus := checkDatabaseHealth()
	status.Database = dbStatus

	// 检查其他服务
	checkServicesHealth(&status)

	// 确定整体状态
	if dbStatus.Status == "healthy" && allServicesHealthy(status.Services) {
		status.Status = "healthy"
	} else {
		status.Status = "unhealthy"
	}

	return status
}

// checkDatabaseHealth 检查数据库健康状态
func checkDatabaseHealth() DatabaseHealthStatus {
	if DB == nil {
		return DatabaseHealthStatus{
			Status:     "unhealthy",
			Connection: "disconnected",
			Error:      "数据库连接为空",
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return DatabaseHealthStatus{
			Status:     "unhealthy",
			Connection: "error",
			Error:      fmt.Sprintf("获取底层数据库连接失败: %v", err),
		}
	}

	// 测试连接
	start := time.Now()
	err = sqlDB.Ping()
	latency := time.Since(start)

	if err != nil {
		return DatabaseHealthStatus{
			Status:     "unhealthy",
			Connection: "failed",
			Error:      fmt.Sprintf("数据库连接失败: %v", err),
		}
	}

	return DatabaseHealthStatus{
		Status:     "healthy",
		Connection: "connected",
		Latency:   latency,
	}
}

// checkServicesHealth 检查其他服务健康状态
func checkServicesHealth(status *HealthStatus) {
	// 检查JWT服务
	jwtStatus := checkJWTService()
	status.Services["jwt"] = jwtStatus

	// 检查缓存服务
	cacheStatus := checkCacheService()
	status.Services["cache"] = cacheStatus

	// 可以添加更多服务检查
	// status.Services["redis"] = checkRedisService()
}

// checkJWTService 检查JWT服务状态
func checkJWTService() map[string]interface{} {
	// 简单检查JWT配置是否正确
	jwtSecret := getJWTKey()
	
	if jwtSecret == nil || len(jwtSecret) == 0 {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  "JWT密钥未配置",
		}
	}

	return map[string]interface{}{
		"status": "healthy",
		"config": "configured",
	}
}

// checkCacheService 检查缓存服务状态
func checkCacheService() map[string]interface{} {
	if Cache == nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  "缓存客户端未初始化",
		}
	}

	// 尝试从缓存获取一个测试键
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testKey := "health_check_test"
	testValue := "ok"
	
	// 设置测试键
	if err := Cache.Set(ctx, testKey, testValue, 5*time.Second); err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  fmt.Sprintf("缓存写入失败: %v", err),
		}
	}

	// 读取测试键
	value, err := Cache.Get(ctx, testKey)
	if err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  fmt.Sprintf("缓存读取失败: %v", err),
		}
	}

	// 删除测试键
	Cache.Delete(ctx, testKey)

	if value != testValue {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  "缓存数据不一致",
		}
	}

	config := GetCacheConfig()
	return map[string]interface{}{
		"status": "healthy",
		"config": config,
	}
}

// allServicesHealthy 检查所有服务是否健康
func allServicesHealthy(services map[string]interface{}) bool {
	for _, service := range services {
		if serviceMap, ok := service.(map[string]interface{}); ok {
			if status, exists := serviceMap["status"]; exists {
				if status != "healthy" {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// GetDatabaseStats 获取数据库统计信息
func GetDatabaseStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	stats := sqlDB.Stats()
	
	return map[string]interface{}{
		"status":         "connected",
		"open_connections": stats.OpenConnections,
		"in_use":         stats.InUse,
		"idle":           stats.Idle,
		"max_open_conns": stats.MaxOpenConnections,
		"wait_count":     stats.WaitCount,
		"wait_duration":   stats.WaitDuration.String(),
		"max_idle_closed": stats.MaxIdleClosed,
		"max_idle_time":   "0", // Go 1.15+ 中 MaxIdleTime 字段不存在
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}
