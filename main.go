package main

import (
	"theing/gin_study/common"

	"github.com/gin-gonic/gin"
)

func main() {

	db := common.InitDB() // 初始化数据库
	defer db.DB()
	r := gin.Default()
	r = CollectRoute(r) // 路由中的collectroute，是一个gin的引擎，返回的也是一个引擎，可以说是代理服务。

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
