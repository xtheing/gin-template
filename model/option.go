package model

import "github.com/jinzhu/gorm/dialects/postgres"

// 字段表的定义，数据库相关

type Option_industry struct { //定义数据类型和字段，一直没有明白数据库表 是 users，也不是user
	Name           string         `json:"name" gorm:"type:varchar(256);not null"`
	Industry_jsonb postgres.Jsonb `json:"industry_jsonb" gorm:"type:jsonb;not null"`
}

type Option_major struct {
	Uuid       string
	Major      string
	Major_json postgres.Jsonb
}
