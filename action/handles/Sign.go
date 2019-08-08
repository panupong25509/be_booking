package handles

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/action/repositories"
	"github.com/panupong25509/be_booking_sign/models"
)

func AddSign(c echo.Context) error {
	success, err := repositories.AddSign(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetAllSign(c echo.Context) error {
	signs, err := repositories.GetAllSign(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, signs)
}

func GetSignByID(c echo.Context) error {
	sign, err := repositories.GetSignByID(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, sign)
}

func GetSignByName(c echo.Context) error {
	success, err := repositories.GetSignByName(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func DeleteSign(c echo.Context) error {
	success, err := repositories.DeleteSign(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func UpdateSign(c echo.Context) error {
	success, err := repositories.UpdateSign(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}
