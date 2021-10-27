package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 一般数据返回封装
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// 成功的返回
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

// 失败的返回
func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 400, data, msg)
}
