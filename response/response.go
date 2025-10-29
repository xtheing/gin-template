package response

import (
	"theing/gin-template/common"
	"theing/gin-template/middleware"

	"github.com/gin-gonic/gin"
)

// Success 成功响应 - 使用新的统一响应格式
func Success(ctx *gin.Context, data interface{}, msg string) {
	// 如果没有自定义消息，使用默认成功消息
	if msg == "" {
		msg = common.GetErrorMessage(common.CodeSuccess)
	}

	// 构造响应数据
	responseData := gin.H{
		"message": msg,
		"data":    data,
	}

	// 使用新的统一响应函数
	middleware.SuccessResponse(ctx, responseData)
}

// Fail 失败响应 - 使用新的统一错误格式
func Fail(ctx *gin.Context, data interface{}, msg string) {
	appErr := common.NewAppError(common.CodeInvalidParams, msg, "")
	middleware.ErrorResponse(ctx, appErr)
}

// FailWithError 使用应用错误返回失败响应
func FailWithError(ctx *gin.Context, appErr *common.AppError) {
	middleware.ErrorResponse(ctx, appErr)
}

// Response 保留原有接口以兼容旧代码，但内部使用新格式
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	if httpStatus >= 200 && httpStatus < 300 {
		// 成功响应
		responseData := gin.H{
			"message": msg,
			"data":    data,
		}
		middleware.SuccessResponse(ctx, responseData)
	} else {
		// 错误响应
		appErr := common.NewAppError(common.ErrorCode(code), msg, "")
		middleware.ErrorResponse(ctx, appErr)
	}
}

// PaginationResponse 分页响应
func PaginationResponse(ctx *gin.Context, list interface{}, total int64, page int, pageSize int, msg string) {
	if msg == "" {
		msg = "查询成功"
	}

	responseData := gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"message":   msg,
	}

	middleware.SuccessResponse(ctx, responseData)
}
