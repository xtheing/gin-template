package option_controller

import (
	"fmt"
	"net/http"
	"theing/gin_study/response"

	"github.com/gin-gonic/gin"
)

// 测试查询功能
func Login(c *gin.Context) {
	userList := TestSelect(0)
	fmt.Println("usersList", userList)
	fmt.Println("userName", userList[0])
	response.Response(
		c,
		http.StatusOK,
		200,
		gin.H{"users": userList},
		"获取成功")
}
