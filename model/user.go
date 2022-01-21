package model

// 字段表的定义，数据库相关

type User struct { //定义数据类型和字段，一直没有明白数据库表 是 users，也不是user
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"type:varchar(20);not null"`
	Tel      string `gorm:"varchar(110);not null;unique"`
	Password string `gorm:"size:255;not null"`
}
