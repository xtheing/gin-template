package middleware

// 认证中间件

import (
	"net/http"
	"strings"
	"theing/gin_study/common"
	"theing/gin_study/model"

	"github.com/gin-gonic/gin"
)

//gin 的中间件就是一个函数，返回一个handlerfunc
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 authorization header
		tokenString := c.GetHeader("Authorization") // ! 这里应该也可以获取header中的其他的字段

		// validate token formate
		// 如果为空或者不是以Bearer开头，那就是没有token，或者说是错误的token，返回权限不足
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() // 抛弃这一次的请求。
			return
		}
		// 如果格式是正确的，那就提取token的有效部分
		tokenString = tokenString[7:]

		// 解析token ，函数卸载jwt中
		token, claims, err := common.ParseToken(tokenString) // 通过包进行引用这个方法。
		if err != nil || !token.Valid {                      // 如果解析失败，或者解析后的token无效， || !token.Valid 表示或者token是无效的。
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() // 抛弃这一次的请求。
			return
		}
		// 验证通过,获取token中的userid
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户，如果用户不存在
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() // 抛弃这一次的请求。
			return
		}

		// 如果用户存在，将user的信息写入上下文。
		c.Set("user", user) // 自己理解为相当于写入缓存中，为登录状态了。
		c.Next()

		// 接下来就要创建一个用户获取用户信息的路由
	}

}
