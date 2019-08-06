package db

import (
	"github.com/JewlyTwin/echo-restful-api/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func Init() {

	db, err = gorm.Open("mysql", "PiOwrb3U68:jg5xygPrLx@tcp(remotemysql.com:3306)/PiOwrb3U68")
	// defer db.Close()
	if err != nil {
		panic("DB Connection Error")
	}
	db.AutoMigrate(&model.User{})

}

func DbManager() *gorm.DB {
	return db
}
