package route

import (
	"github.com/labstack/echo"
	api "github.com/panupong25509/be_booking_sign2/action"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", api.GetUsers)
	return e
}
