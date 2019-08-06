package route

import (
	api "github.com/JewlyTwin/echo-restful-api/action"
	"github.com/labstack/echo"
)

func Init() *echo.Echo {
	e := echo.New()

	// e.GET("/", api.Home)

	e.GET("/users", api.GetUsers)
	return e
}
