package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/panupong25509/be_booking_sign/config"
)

var db *gorm.DB
var err error

func Init() {
	configuration := config.GetConfig()
	path := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", configuration.DB_TYPE, configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_HOST, configuration.DB_PORT, configuration.DB_NAME)
	db, err = gorm.Open(configuration.DB_TYPE, path)
	// defer db.Close()
	if err != nil {
		panic("DB Connection Error")
	}
	// db.AutoMigrate(&model.User{})

}

func DbManager() *gorm.DB {
	return db
}
