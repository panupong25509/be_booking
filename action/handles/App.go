package handles

import (
	"net/http"

	"github.com/labstack/echo"
)

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "welcome")
}
