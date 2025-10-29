package common

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 常用错误码定义
const (
	// 成功
	CodeSuccess ErrorCode = 0

	// 客户端错误 1000-1999
	CodeInvalidParams    ErrorCode = 1001 // 无效参数
	CodeUnauthorized     ErrorCode = 1002 // 未授权
	CodeForbidden        ErrorCode = 1003 // 禁止访问
	CodeNotFound         ErrorCode = 1004 // 资源不存在
	CodeUserExists       ErrorCode = 1005 // 用户已存在
	CodeUserNotFound     ErrorCode = 1006 // 用户不存在
	CodePasswordError    ErrorCode = 1007 // 密码错误
	CodeTokenInvalid     ErrorCode = 1008 // Token无效
	CodeTokenExpired     ErrorCode = 1009 // Token过期
	CodePasswordTooWeak  ErrorCode = 1010 // 密码过于简单
	CodeBadRequest       ErrorCode = 1011 // 请求错误
	CodeTooManyRequests  ErrorCode = 1012 // 请求过于频繁
	CodeValidationFailed ErrorCode = 1013 // 验证失败

	// 业务逻辑错误 2000-2999
	CodeBusinessError ErrorCode = 2001 // 业务逻辑错误
	CodeDataExists    ErrorCode = 2002 // 数据已存在
	CodeDataNotFound  ErrorCode = 2003 // 数据不存在

	// 服务器错误 5000-5999
	CodeInternalError      ErrorCode = 5001 // 内部服务器错误
	CodeDatabaseError      ErrorCode = 5002 // 数据库错误
	CodeNetworkError       ErrorCode = 5003 // 网络错误
	CodeServiceUnavailable ErrorCode = 5004 // 服务不可用
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode `json:"code"`              // 错误码
	Message    string    `json:"message"`           // 错误信息
	Details    string    `json:"details,omitempty"` // 详细信息
	HTTPStatus int       `json:"-"`                 // HTTP状态码，不返回给前端
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewAppError 创建应用错误
func NewAppError(code ErrorCode, message string, details string) *AppError {
	httpStatus := getHTTPStatusByCode(code)
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
	}
}

// getHTTPStatusByCode 根据错误码获取HTTP状态码
func getHTTPStatusByCode(code ErrorCode) int {
	switch {
	case code == CodeSuccess:
		return http.StatusOK
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 3000:
		return http.StatusUnprocessableEntity
	case code >= 5000 && code < 6000:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// 预定义错误
var (
	ErrInvalidParams      = NewAppError(CodeInvalidParams, "参数错误", "")
	ErrUnauthorized       = NewAppError(CodeUnauthorized, "未授权访问", "")
	ErrForbidden          = NewAppError(CodeForbidden, "禁止访问", "")
	ErrNotFound           = NewAppError(CodeNotFound, "资源不存在", "")
	ErrUserExists         = NewAppError(CodeUserExists, "用户已存在", "")
	ErrUserNotFound       = NewAppError(CodeUserNotFound, "用户不存在", "")
	ErrPasswordError      = NewAppError(CodePasswordError, "密码错误", "")
	ErrTokenInvalid       = NewAppError(CodeTokenInvalid, "Token无效", "")
	ErrTokenExpired       = NewAppError(CodeTokenExpired, "Token已过期", "")
	ErrPasswordTooWeak    = NewAppError(CodePasswordTooWeak, "密码过于简单", "")
	ErrInternalError      = NewAppError(CodeInternalError, "内部服务器错误", "")
	ErrDatabaseError      = NewAppError(CodeDatabaseError, "数据库错误", "")
	ErrNetworkError       = NewAppError(CodeNetworkError, "网络错误", "")
	ErrServiceUnavailable = NewAppError(CodeServiceUnavailable, "服务不可用", "")
)

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code ErrorCode) string {
	messages := map[ErrorCode]string{
		CodeSuccess:            "成功",
		CodeInvalidParams:      "参数错误",
		CodeUnauthorized:       "未授权访问",
		CodeForbidden:          "禁止访问",
		CodeNotFound:           "资源不存在",
		CodeUserExists:         "用户已存在",
		CodeUserNotFound:       "用户不存在",
		CodePasswordError:      "密码错误",
		CodeTokenInvalid:       "Token无效",
		CodeTokenExpired:       "Token已过期",
		CodePasswordTooWeak:    "密码过于简单",
		CodeBadRequest:         "请求错误",
		CodeTooManyRequests:    "请求过于频繁",
		CodeValidationFailed:   "验证失败",
		CodeBusinessError:      "业务逻辑错误",
		CodeDataExists:         "数据已存在",
		CodeDataNotFound:       "数据不存在",
		CodeInternalError:      "内部服务器错误",
		CodeDatabaseError:      "数据库错误",
		CodeNetworkError:       "网络错误",
		CodeServiceUnavailable: "服务不可用",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}
	return "未知错误"
}
