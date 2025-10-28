package routers

// 用于分离路由

import (
	"theing/gin-template/common"
	controller "theing/gin-template/controller"
	"theing/gin-template/controller/admin_controller"
	option_controller "theing/gin-template/controller/options_controller"
	errorMiddleware "theing/gin-template/middleware"

	gin "github.com/gin-gonic/gin"
)

// 接收一个gin 引擎，返回一个引擎，不是很懂
func CollectRoute(r *gin.Engine) *gin.Engine {
	// 初始化监控指标
	common.InitMetrics()
	
	// 添加全局中间件
	r.Use(errorMiddleware.LoggingMiddleware())           // 日志中间件
	r.Use(errorMiddleware.RequestIDMiddleware())        // 请求ID中间件
	r.Use(errorMiddleware.ErrorHandlingMiddleware())     // 全局错误处理中间件
	r.Use(errorMiddleware.ErrorHandlerMiddleware())       // 统一错误处理中间件
	r.Use(errorMiddleware.MetricsMiddleware())           // 性能监控中间件
	r.Use(errorMiddleware.DatabaseMetricsMiddleware())    // 数据库监控中间件

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.Register)                     // 用户注册
			auth.POST("/login", controller.UserLogin)                       // 用户登录
			auth.GET("/info", errorMiddleware.AuthMiddleware(), controller.Info) // 获取用户信息（需要认证）
		}

		// 选项相关路由
		options := api.Group("/options")
		{
			options.GET("/industry", option_controller.GetIndustryList)     // 获取行业领域列表
			options.GET("/profession", option_controller.GetProfessionList) // 获取专业选项分类
		}

		// 管理员相关路由
		admin := api.Group("/admin")
		{
			admin.POST("/login", admin_controller.AdminLogin) // 管理员登录
		}

		// 健康检查路由
		health := api.Group("/health")
		{
			health.GET("/", controller.HealthCheck)      // 系统健康检查
			health.GET("/database", controller.DatabaseHealth) // 数据库健康检查
			health.GET("/stats", controller.DatabaseStats)     // 数据库统计信息
			health.GET("/info", controller.SystemInfo)          // 系统信息
		}

		// 监控指标路由
		api.GET("/metrics", errorMiddleware.MetricsHandler()) // Prometheus 指标

		// 兼容旧路由
		api.GET("/auth/info2", option_controller.Login)
		api.GET("/GetIndustryList", option_controller.GetIndustryList)     // 获取行业领域列表json格式
		api.GET("/GetProfessionList", option_controller.GetProfessionList) // 获取专业选项分类
		api.GET("/admin/AdminLogin", admin_controller.AdminLogin)          // 获取专业选项分类
	}

	return r
}
