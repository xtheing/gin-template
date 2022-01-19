package option_controller

import (
	"theing/gin_study/common"
	"theing/gin_study/model"
)

//测试查询功能
func TestSelect(userId int) (userList []model.User) {
	DB := common.GetDB()
	DB.Raw("select * from users where id > ?", userId).Scan(&userList)
	return
}
