package repositories

import (
	"github.com/panupong25509/be_booking_sign/models"
	"github.com/panupong25509/be_booking_sign/db"
	"io"
	"os"

	"github.com/labstack/echo"
)

func AddSign(c echo.Context, data) error {

}

func GetAllSign(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	signs := []models.Sign{}
	db.Find(&signs)
	return signs, nil
}

func GetSignByID(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	// if data["id"]
	db.First(&sign, data["id"])
	return sign, nil
}

func GetSignByName(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	db.Where("sign_name = ?", data["signname"]).First(&user)
	return sign, nil
}
func DeleteSign(c echo.Context) (interface{}, interface{}) {
	
}
func UpdateSign(c echo.Context) (interface{}, interface{}) {
	
}

func Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(`D:\fe_booking_sign\public\img\` + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
