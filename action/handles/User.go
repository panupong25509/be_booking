package handles

import (
	"github.com/labstack/echo"
)

func Register(c echo.Context) error {
	data := DynamicPostForm(c)
	success, err := repositories.Register(c, data)
	if err != nil {
		status := err.(models.Error)
		// return c.Render(status.Code, r.JSON(status))
	}
	// return c.Render(200, r.JSON(success))
}