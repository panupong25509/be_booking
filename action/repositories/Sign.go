package repositories

import (
	"io"
	"os"
	"strconv"

	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"

	"github.com/labstack/echo"
)

func AddSign(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	if !sign.CheckParamPostForm(c) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	checkSign, err := GetSignByName(c)
	if checkSign != nil {
		return nil, err
	}
	_, err = UploadImg(c)
	if err != nil {
		return nil, err
	}
	sign.CreateSignModel(c)
	db.Create(&sign)
	return models.Success{200, "success"}, nil
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
	signs := []models.Sign{}
	db.Where("sign_name = ?", c.FormValue("signname")).First(&signs)
	if len(signs) != 0 {
		return signs[0], nil
	}
	return nil, models.Error{500, "Not have sign"}
}

func DeleteSign(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	id, _ := strconv.Atoi(data["id"].(string))
	err := db.Find(&sign, id)
	if err != nil {
		return nil, models.Error{400, "ไม่มีป้ายนี้ใน Database"}
	}
	os.Remove(`D:\fe_booking_sign\public\img\` + sign.Picture)
	_ = db.Delete(&sign)
	return models.Success{200, "success"}, nil
}

func UpdateSign(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	if !sign.CheckParamPostForm(c) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	_, err := UploadImg(c)
	if err != nil {
		return nil, err
	}
	sign.CreateSignModel(c)
	oldSign, err := GetSignByID(c, data)
	if err != nil {
		return nil, err
	}
	oldPicture := oldSign.(models.Sign).Picture
	os.Remove(`D:\fe_booking_sign\public\img\` + oldPicture)
	db.Update(&sign)
	return models.Success{200, "success"}, nil
}

func UploadImg(c echo.Context) (interface{}, interface{}) {
	file, err := c.FormFile("file")
	if err != nil {
		return err, nil
	}
	src, err := file.Open()
	if err != nil {
		return err, nil
	}
	defer src.Close()
	dst, err := os.Create(`D:\fe_booking_sign\public\img\` + c.FormValue("signname"))
	if err != nil {
		return err, nil
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err, nil
	}

	return nil, nil
}
