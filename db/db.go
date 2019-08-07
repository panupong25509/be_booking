package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func Init() {

	db, err = gorm.Open("postgres", "postgres://postgres:postgres@103.86.49.57:5432/practice?sslmode=disable")
	// defer db.Close()
	if err != nil {
		panic("DB Connection Error")
	}
	// db.AutoMigrate(&model.User{})

}

func DbManager() *gorm.DB {
	return db
}
