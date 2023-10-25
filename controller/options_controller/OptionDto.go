package option_controller

import "theing/gin-template/model"

// UserDTO 用户数据传输对象，相当于是一个格式的定义和转换。
type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

// 转换的函数 将model.User转换为UserDto
func ToUserDto(user model.User) UserDto {
	return UserDto{ // 返回一个新的UserDto格式
		Name:      user.Username,
		Telephone: user.Tel,
	}
}
