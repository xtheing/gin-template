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

// 获取行业领域列表json格式
func GetIndustryList(c *gin.Context) {
	industryList := GetIndustry()
	response.Response(
		c,
		http.StatusOK,
		200,
		gin.H{"industry": industryList},
		"获取成功")
}

// 获取专业选项分类
func GetProfessionList(c *gin.Context) {
	professionList := GetProfession()
	if professionList == nil {
		response.Response(
			c,
			http.StatusOK,
			200,
			gin.H{"profession": nil},
			"获取失败")
	} else {
		response.Response(
			c,
			http.StatusOK,
			200,
			gin.H{"profession": professionList},
			"获取成功")
	}
}
