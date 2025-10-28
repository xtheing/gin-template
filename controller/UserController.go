package controller

// 逻辑相关

import (
	"log"
	"strings"
	"theing/gin-template/common"
	"theing/gin-template/dto"
	"theing/gin-template/model"
	"theing/gin-template/response"
	"theing/gin-template/utils"

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
		response.FailWithError(c, common.NewAppError(common.CodeInvalidParams, "手机号必须为11位", ""))
		return
	}

	// 使用新的密码验证逻辑
	passwordValidation := utils.ValidatePassword(password)
	if !passwordValidation.IsValid {
		errorMsg := "密码不符合要求：" + strings.Join(passwordValidation.Errors, "；")
		response.FailWithError(c, common.NewAppError(common.CodePasswordTooWeak, errorMsg, strings.Join(passwordValidation.Suggestions, "；")))
		return
	}

	// 如果密码强度较弱，给出建议但不阻止注册
	if passwordValidation.Strength != utils.PasswordStrong {
		log.Printf("用户注册密码强度提醒：%v", strings.Join(passwordValidation.Suggestions, "；"))
	}

	// 如果名称没有传入，给一个10位的随机字符串
	if len(name) == 0 {
		name = utils.RandomString(10)
	}

	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.FailWithError(c, common.ErrUserExists)
		return
	}

	// 创建用户，用户的密码是不能明文保存的，所有所有的密码都应该加密保存，写法如下，都是通用的。
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // ! 这里的加密方式自己也应该研究一下。
	if err != nil {
		response.FailWithError(c, common.NewAppError(common.CodeInternalError, "加密错误", err.Error()))
		return // 返回前端一个错误，这里是一个系统基本的错误。
	}
	newUser := model.User{
		Username:  name,
		Telephone: telephone,
		Password:  string(hasedPassword), // 创建密码的时候不能明文
	}
	if err := DB.Create(&newUser).Error; err != nil {
		response.FailWithError(c, common.NewAppError(common.CodeDatabaseError, "用户创建失败", err.Error()))
		return
	}
	// 返回结果
	response.Success(c, gin.H{"user_id": newUser.ID}, "注册成功")
}

// 用户登录
func UserLogin(c *gin.Context) {
	// 验证参数
	type PostUserLogin struct {
		Telephone string `json:"telephone" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	var login PostUserLogin

	if err := c.ShouldBindJSON(&login); err != nil {
		response.FailWithError(c, common.NewAppError(common.CodeInvalidParams, "参数错误", err.Error()))
		return
	}

	print(login.Password, login.Telephone)

	// 数据验证
	if len(login.Telephone) != 11 {
		response.FailWithError(c, common.NewAppError(common.CodeInvalidParams, "手机号必须为11位", ""))
		return
	}

	// 密码长度基本检查
	if len(login.Password) < 6 {
		response.FailWithError(c, common.NewAppError(common.CodeInvalidParams, "密码不能少于6位", ""))
		return
	}

	// 判断手机号是否存在
	DB := common.GetDB() // 引入 DB实例
	var user model.User
	if err := DB.Where("telephone = ?", login.Telephone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.FailWithError(c, common.ErrUserNotFound)
		} else {
			response.FailWithError(c, common.NewAppError(common.CodeDatabaseError, "查询用户失败", err.Error()))
		}
		return
	}

	// 判断密码是否正确
	// 判定用户密码的时候就用bcrypt的方法进行判定，第一个参数是原始的加密后的密码，第二个参数就是需要对比的密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil { // 如果有err 就提示密码错误。这里应该也是和if同级的一个代码。
		response.FailWithError(c, common.ErrPasswordError)
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.FailWithError(c, common.NewAppError(common.CodeInternalError, "token 发放失败", err.Error()))
		log.Printf("token generate error : %v", err) // 遇到了这个问题记录一下日志。
		return
	}

	// 返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0 // 找到了就不为零，找不到应该就是零了
}

// 登录用户获取自己的信息
func Info(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.FailWithError(c, common.ErrUnauthorized)
		return
	}

	response.Success(c, gin.H{"user": dto.ToUserDto(user.(model.User))}, "获取用户信息成功")
}
