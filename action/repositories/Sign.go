package repositories

import (
	"io"
	"os"

	"github.com/labstack/echo"
)

func AddSign(c echo.Context, data) error {

}

func GetAllSign(c echo.Context) (interface{}, interface{}) {
	
}

func GetSignByID(c echo.Context) (interface{}, interface{}) {
	
}

func GetSignByName(c echo.Context) (interface{}, interface{}) {
	
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
