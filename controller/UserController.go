package controller

// 逻辑相关

import (
	"log"
	"net/http"
	"theing/gin_study/common"
	"theing/gin_study/dto"
	"theing/gin_study/model"
	"theing/gin_study/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户注册功能
func Register(c *gin.Context) {
	DB := common.GetDB() //引入DB实例

	// 获取参数
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不能少于6位",
		})
		return
	}

	// 如果名称没有传入，给一个10位的随机字符串
	if len(name) == 0 {
		name = utils.RandomString(10)
	}

	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已存在",
		})
		return
	}

	// 创建用户，用户的密码是不能明文保存的，所有所有的密码都应该加密保存，写法如下，都是通用的。
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // ! 这里的加密方式自己也应该研究一下。
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ // 返回前端一个错误，这里是一个系统基本的错误。
			"code": 500,
			"msg":  "加密错误",
		})
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword), // 创建密码的时候不能明文
	}
	DB.Create(&newUser)
	// 返回结果
	c.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 用户登录
func Login(c *gin.Context) {
	// 获取参数

	DB := common.GetDB() // 引入 DB实例
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不能少于6位",
		})
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}

	// 判断密码是否正确
	// 判定用户密码的时候就用bcrypt的方法进行判定，第一个参数是原始的加密后的密码，第二个参数就是需要对比的密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil { // 如果有err 就提示密码错误。这里应该也是和if同级的一个代码。
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}
	// 发放token

	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ // 给前端一个报错的提示。
			"code": 500,
			"msg":  "系统异常",
		})
		log.Printf("token generate error : %v", err) // 遇到了这个问题记录一下日志。
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { // 找到了就不为零，找不到应该就是零了
		return true
	}
	return false
}

// 登录用户获取自己的信息
func Info(c *gin.Context) {
	user, _ := c.Get("user")                                                                           // 直接获取登录用户的id和信息，应该就是gin.Context的作用而获取到的。
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}}) // 返回结果，这里的user是一个model.User类型的，所以可以直接转换成dto.ToUserDto
	// ? 接下来需要将我们的中间件用来保护用户信息的接口。路由中
	// todo 这里返回的用户信息不应该是用户所有的信息，需要进行设置，需要封装一个返回的固定格式。
}
