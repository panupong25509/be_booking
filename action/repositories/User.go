package repositories

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	user := models.User{}
	if !user.CheckParams(data) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	_, err := GetUserByUsername(c, data)
	if err == nil {
		return nil, models.Error{500, "Username นี้มีผู้ใช้แล้ว"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_ = user.CreateModel(data, string(hash))
	db.NewRecord(user)
	err = db.Create(&user)
	if err != nil {
		return nil, err
	}
	success := models.Success{"success"}
	return success, nil
}

func GetUserByUsername(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	if data["username"] == nil {
		return nil, models.Error{400, "ไม่มี username"}
	}
	username := data["username"].(string)
	user := models.Users{}
	db.First(&user, "username = (?)", username)
	if len(user) == 0 {
		return nil, models.Error{400, "ไม่มี username"}
	}
	return user[0], nil
}
