package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"theing/gin-template/common"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`      // 业务错误码
	Message   string      `json:"message"`   // 错误信息
	Data      interface{} `json:"data"`      // 响应数据
	RequestID string      `json:"request_id"` // 请求ID
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// ErrorResponseStruct 错误响应结构
type ErrorResponseStruct struct {
	Code      int    `json:"code"`      // 业务错误码
	Message   string `json:"message"`   // 错误信息
	Details   string `json:"details"`   // 详细错误信息
	RequestID string `json:"request_id"` // 请求ID
	Timestamp int64  `json:"timestamp"` // 时间戳
	TraceID   string `json:"trace_id,omitempty"` // 追踪ID
}

// ErrorHandlingMiddleware 全局错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录panic堆栈信息
		log.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())
		
		// 返回服务器内部错误
	ErrorResponse(c, common.ErrInternalError)
		c.Abort()
	})
}

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// 根据错误类型处理
			switch e := err.Err.(type) {
			case *common.AppError:
				// 应用自定义错误
				ErrorResponse(c, e)
			default:
				// 未知错误
				log.Printf("Unknown error: %v", err)
				ErrorResponse(c, common.ErrInternalError)
			}
			
			c.Abort()
			return
		}
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		requestID, _ := param.Keys["request_id"].(string)
		return fmt.Sprintf("[%s] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			time.Now().Format("2006-01-02 15:04:05"),
			requestID,
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	requestID, _ := c.Get("request_id")
	response := Response{
		Code:      int(common.CodeSuccess),
		Message:   common.GetErrorMessage(common.CodeSuccess),
		Data:      data,
		RequestID: requestID.(string),
		Timestamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, appErr *common.AppError) {
	requestID, _ := c.Get("request_id")
	
	// 记录错误日志
	if appErr.Code >= common.CodeInternalError {
		log.Printf("Internal Error [%s]: %s", requestID.(string), appErr.Error())
	}
	
	response := ErrorResponseStruct{
		Code:      int(appErr.Code),
		Message:   appErr.Message,
		Details:   appErr.Details,
		RequestID: requestID.(string),
		Timestamp: time.Now().Unix(),
		TraceID:   generateTraceID(),
	}
	
	c.JSON(appErr.HTTPStatus, response)
}

// HandleError 处理错误的便捷函数
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*common.AppError); ok {
		ErrorResponse(c, appErr)
	} else {
		// 包装为应用错误
		appErr := common.NewAppError(common.CodeInternalError, "内部服务器错误", err.Error())
		ErrorResponse(c, appErr)
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	return fmt.Sprintf("trace_%d", time.Now().UnixNano())
}
