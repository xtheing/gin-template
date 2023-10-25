package common

// 初始化数据库相关

import (
	"fmt"
	"theing/gin-template/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("db_name")
	fmt.Println("driverName: ", driverName)
	if driverName == "mysql" {
		host := viper.GetString("mysql_host")
		port := viper.GetString("mysql_port")
		database := viper.GetString("mysql_database")
		username := viper.GetString("mysql_username")
		password := viper.GetString("mysql_password")
		charset := viper.GetString("mysql_charset")
		args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
			username,
			password,
			host,
			port,
			database,
			charset)
		db, err := gorm.Open(mysql.Open(args))
		if err != nil {
			panic("mysql数据库连接失败: " + err.Error())
		}
		db.AutoMigrate(&model.User{}) // 调用的是类名，自动创建数据表
		DB = db
		return db
	} else if driverName == "postgres" {
		// db, err := gorm.Open(postgres.Open(args))
		// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		host := viper.GetString("postgres_host")
		port := viper.GetString("postgres_port")
		database := viper.GetString("postgres_database")
		username := viper.GetString("postgres_username")
		password := viper.GetString("postgres_password")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			host,
			username,
			password,
			database,
			port)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("pgsql数据库连接失败" + err.Error())
		}
		db.AutoMigrate(&model.User{}) // 调用的是类名，自动创建数据表
		DB = db
		return db
	}
	panic("pgsql数据库连接失败")
}

// 定义一个方法来获取DB实例，需要在controller中引入
func GetDB() *gorm.DB {
	isPrintSql := viper.GetString("is_print_sql") // 判定是否打印sql的一个debug模式
	if isPrintSql == "true" {
		return DB.Debug()
	} else if isPrintSql == "false" {
		return DB
	} else {
		return DB
	}
}
