package route

import (
	"github.com/labstack/echo"
	api "github.com/panupong25509/be_booking_sign/action/handles"
)

func Init() *echo.Echo {
	e := echo.New()

	// e.GET("/", api.GetUsers)
	e.POST("/login", api.Login)
	e.POST("/register", api.Register)
	e.GET("/user", api.GetUserById)
	return e
}
