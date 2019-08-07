package handles

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/action/repositories"
	"github.com/panupong25509/be_booking_sign/models"
)

func Login(c echo.Context) error {
	data := DynamicPostForm(c)
	jwt, err := repositories.Login(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status.Message)
	}
	return c.JSON(200, jwt)
}

func Register(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.Register(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetUserByUsername(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.GetUserByUsername(c, data)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}

func GetUserById(c echo.Context) error {
	success, err := repositories.GetUserById(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, success)
}
