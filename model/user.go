package model

// 字段表的定义，数据库相关

import "gorm.io/gorm"

type User struct { //定义数据类型和字段
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}
