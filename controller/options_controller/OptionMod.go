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

// 获取行业领域列表json格式
func GetIndustry() (industry_jsonb model.Option_industry) {
	DB := common.GetDB()
	// DB.Raw("select uuid, name, industry_jsonb from option_industry where uuid = 'OPI20211230113204VCJHWVHT' and name = 'industry_jsonb';").Scan(&industry)
	DB.Select("uuid, name, industry_jsonb").Table("option_industry").Where("uuid = ? and name = ?", "OPI20211230113204VCJHWVHT", "industry_jsonb").Scan(&industry_jsonb)
	return
}

// 获取专业选项分类
func GetProfession() (professionList []model.Option_major) {
	DB := common.GetDB()
	// DB.Raw("select uuid, major,major_json from option_major where uuid = 'OMJ2021122814344084JTNUCY';").Scan(&professionList)
	DB.Select("uuid, major,major_json").Table("option_major").Where("uuid = ?", "OMJ2021122814344084JTNUCY").Scan(&professionList)
	return
}
