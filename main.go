package main

import (
	"theing/gin_study/common"

	"github.com/gin-gonic/gin"
)

func main() {

	db := common.GetDB()
	defer db.DB()
	r := gin.Default()
	r = CollectRoute(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
