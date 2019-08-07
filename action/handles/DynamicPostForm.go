package handles

import (
	"github.com/labstack/echo"
)

func DynamicPostForm(c echo.Context) map[string]interface{} {
	c.Request().ParseForm()
	param := c.Request().PostForm
	m := make(map[string]interface{})
	for key, value := range param {
		m[key] = value[0]
	}
	return m
}
