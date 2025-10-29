package controller

import (
	"theing/gin-template/common"
	"theing/gin-template/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck 系统健康检查
func HealthCheck(c *gin.Context) {
	status := common.CheckSystemHealth()

	if status.Status == "healthy" {
		response.Success(c, status, "系统健康")
	} else {
		// 使用新的错误响应格式
		appErr := common.NewAppError(common.CodeServiceUnavailable, "系统不健康", "")
		response.FailWithError(c, appErr)
	}
}

// DatabaseHealth 数据库健康检查
func DatabaseHealth(c *gin.Context) {
	// 获取完整系统状态，然后提取数据库部分
	status := common.CheckSystemHealth()
	dbStatus := status.Database

	if dbStatus.Status == "healthy" {
		response.Success(c, dbStatus, "数据库连接正常")
	} else {
		appErr := common.NewAppError(common.CodeDatabaseError, "数据库连接异常", dbStatus.Error)
		response.FailWithError(c, appErr)
	}
}

// DatabaseStats 数据库统计信息
func DatabaseStats(c *gin.Context) {
	stats := common.GetDatabaseStats()

	if stats["status"] == "connected" {
		response.Success(c, stats, "获取数据库统计信息成功")
	} else {
		appErr := common.NewAppError(common.CodeDatabaseError, "数据库连接异常", "")
		response.FailWithError(c, appErr)
	}
}

// SystemInfo 系统信息
func SystemInfo(c *gin.Context) {
	info := gin.H{
		"service":     "gin-template",
		"version":     "2.0.0",
		"environment": gin.Mode(),
		"go_version":  "go1.19+",
		"features": []string{
			"统一错误处理",
			"请求ID追踪",
			"密码强度验证",
			"JWT安全配置",
			"数据库连接池",
			"健康检查",
		},
	}

	response.Success(c, info, "获取系统信息成功")
}
