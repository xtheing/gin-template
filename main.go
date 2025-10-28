package main

import (
	"fmt"
	"os"
	"theing/gin-template/common"
	"theing/gin-template/routers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	InitConfig() // 项目开始的时候就应该读取配置文件
	isDebug := viper.GetString("is_debug")
	if isDebug == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	db := common.InitDB() // 初始化数据库
	defer db.DB()
	
	// 初始化缓存
	if err := common.InitCache(); err != nil {
		fmt.Printf("缓存初始化失败: %v\n", err)
		// 缓存初始化失败不影响服务启动
	}
	
	r := gin.Default()
	r = routers.CollectRoute(r) // 路由中的collectroute，是一个gin的引擎，返回的也是一个引擎，可以说是代理服务。
	port := viper.GetString("server_port")
	if port != "" {
		r.Run(":" + port) // 如果有port，就使用这个端口号
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 配置功能
func InitConfig() {
	// 获取当前的工作目录
	workDir, _ := os.Getwd()
	// 设置要读取的配置文件
	viper.SetConfigName("application")
	// 设置要读取的文件类型
	viper.SetConfigType("yaml")
	// 设置文件的路径
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.SetEnvPrefix("tpl") // 读取环境变量的前缀为tpl
	viper.AutomaticEnv()      //从环境变量中获取配置
}
