package routers

// 用于分离路由

import (
	controller "theing/gin_study/controller"
	"theing/gin_study/controller/admin_controller"
	option_controller "theing/gin_study/controller/options_controller"
	middleware "theing/gin_study/middleware"

	gin "github.com/gin-gonic/gin"
)

// 接收一个gin 引擎，返回一个引擎，不是很懂
func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)                     // 导入包中的模块
	r.POST("/api/auth/login", controller.UserLogin)                       // 用户登录
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info) // middleware.AuthMiddleware()，表示利用中间件包含用户的信息。
	r.GET("/api/auth/info2", option_controller.Login)
	r.GET("/api/GetIndustryList", option_controller.GetIndustryList)     // 获取行业领域列表json格式
	r.GET("/api/GetProfessionList", option_controller.GetProfessionList) // 获取专业选项分类
	r.GET("/api/admin/AdminLogin", admin_controller.AdminLogin)          // 获取专业选项分类

	return r
}
