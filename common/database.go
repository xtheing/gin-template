package common

// 初始化数据库相关

import (
	"fmt"
	"theing/gin-template/model"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("db_name")
	fmt.Println("driverName: ", driverName)
	
	var db *gorm.DB
	var err error
	
	if driverName == "mysql" {
		db, err = initMySQLDB()
	} else if driverName == "postgres" {
		db, err = initPostgresDB()
	} else {
		panic("不支持的数据库类型: " + driverName)
	}
	
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	
	// 配置连接池
	configureConnectionPool(db)
	
	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
	
	DB = db
	return db
}

// initMySQLDB 初始化 MySQL 数据库连接
func initMySQLDB() (*gorm.DB, error) {
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
	
	return gorm.Open(mysql.Open(args), &gorm.Config{})
}

// initPostgresDB 初始化 PostgreSQL 数据库连接
func initPostgresDB() (*gorm.DB, error) {
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
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// configureConnectionPool 配置数据库连接池
func configureConnectionPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("获取底层数据库连接失败: %v\n", err)
		return
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间
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
