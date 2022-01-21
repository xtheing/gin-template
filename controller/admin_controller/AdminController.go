package admin_controller

import (
	"log"
	"net/http"
	"theing/gin_study/common"
	"theing/gin_study/model"
	"theing/gin_study/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// admin登录
func AdminLogin(c *gin.Context) {
	// 获取参数

	DB := common.GetDB() // 引入 DB实例
	tel := c.PostForm("tel")
	password := c.PostForm("password")

	// 数据验证
	if len(tel) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Select("username,tel,password").Table("users").Where("tel = ?", tel).First(&user)
	// DB.Where("tel = ?", tel).First(&user)
	// DB.Raw("select tel from users where id > ?", userId).Scan(&userList)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	// 判断密码是否正确
	// 判定用户密码的时候就用bcrypt的方法进行判定，第一个参数是原始的加密后的密码，第二个参数就是需要对比的密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil { // 如果有err 就提示密码错误。这里应该也是和if同级的一个代码。
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "token 发放失败")
		log.Printf("token generate error : %v", err) // 遇到了这个问题记录一下日志。
		return
	}

	// 返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}
