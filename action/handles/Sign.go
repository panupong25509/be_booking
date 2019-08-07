package handles

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/action/repositories"
	"github.com/panupong25509/be_booking_sign/models"
)

func AddSign(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.AddSign(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetAllSign(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.GetAllSign(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetSignByID(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.GetSignByID(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetSignByName(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.GetSignByName(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func DeleteSign(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.DeleteSign(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func UpdateSign(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.UpdateSign(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}
