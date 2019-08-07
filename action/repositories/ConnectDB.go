package repositories

import (
	"github.com/labstack/echo"
)

func ConnectDB(c echo.Context) (*pop.Connection, interface{}) {
	db, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return nil, models.Error{500, "can't connect db"}
	}
	return db, nil
}
