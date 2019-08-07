package route

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/action/handles"
)

func Init() *echo.Echo {
	e := echo.New()

	e.POST("/register", handles.Register)
	return e
}
