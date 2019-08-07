package repositories

import (
	"github.com/JewlyTwin/be_booking_sign/models"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db, err := ConnectDB(c)
	if err != nil {
		return nil, err
	}
	user := models.User{}
	if !user.CheckParams(data) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	_, err = GetUserByUsername(c, data)
	if err == nil {
		return nil, models.Error{500, "Username นี้มีผู้ใช้แล้ว"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_ = user.CreateModel(data, string(hash))
	err = db.Create(&user)
	if err != nil {
		return nil, err
	}
	return Success(nil), nil
}
