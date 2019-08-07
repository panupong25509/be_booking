package repositories

import (
	"io"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
)

func AddSign(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	sign := models.Sign{}
	if !sign.CheckParamPostForm(data) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	samename, err := GetSignByName(c, data)
	if err != nil {
		return nil, err
	}
	if samename != nil {
		return nil, models.Error{400, "ชื่อป้ายนี้มีอยู่ในระบบแล้ว"}
	}
	file, err := UploadImg(c, data["signname"].(string))
	if err != nil {
		return nil, err
	}
	sign.CreateSignModel(data, file.(string))
	db.Create(&sign)
	return models.Success{200, "success"}, nil
}

func GetAllSign(c echo.Context) (interface{}, interface{}) {

}

func GetSignByID(c echo.Context) (interface{}, interface{}) {

}

func GetSignByName(c echo.Context) (interface{}, interface{}) {

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
	if !sign.CheckParamPostForm(data) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	file, err := UploadImg(c, data["signname"].(string))
	if err != nil {
		return nil, err
	}
	sign.CreateSignModel(data, file.(string))
	oldSign, err := GetSignById(c, sign.ID)
	if err != nil {
		return nil, err
	}
	oldPicture := oldSign.(models.Sign).Picture
	os.Remove(`D:\fe_booking_sign\public\img\` + oldPicture)
	db.Update(&sign)
	return models.Success{200, "success"}, nil
}

func UploadImg(c echo.Context, filename string) (interface{}, interface{}) {
	file, err := c.FormFile("file")
	if err != nil {
		return err, nil
	}
	src, err := file.Open()
	if err != nil {
		return err, nil
	}
	defer src.Close()
	dst, err := os.Create(`D:\fe_booking_sign\public\img\` + filename)
	if err != nil {
		return err, nil
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err, nil
	}

	return nil, nil
}
