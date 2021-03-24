package gorm

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

func Connect(uri string) *gorm.DB{
	db, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}
