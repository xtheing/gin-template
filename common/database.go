package common

// 初始化数据库相关

import (
	"fmt"
	"theing/gin_study/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// driverName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args))
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{}) // 调用的是类名
	return db
}

// 定义一个方法来获取DB实例
func GetDB() *gorm.DB {
	return DB
}
