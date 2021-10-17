package main

// 用于分离路由

import (
	"theing/gin_study/controller"

	"github.com/gin-gonic/gin"
)

// 接收一个gin 引擎，返回一个引擎
func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register) // 导入包中的模块
	r.POST("/api/auth/login", controller.Login)       // 导入包中的模块
	return r
}
