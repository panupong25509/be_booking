package api

import (
	"net/http"

	"github.com/JewlyTwin/echo-restful-api/db"
	"github.com/JewlyTwin/echo-restful-api/model"
	"github.com/labstack/echo"
)

func GetUsers(c echo.Context) error {
	db := db.DbManager()
	users := model.User{
		FirstName: "supaiwit",
		LastName:  "SDASD",
		Age:       20,
		Email:     "sacxsadsa@gmail.com",
	}
	db.NewRecord(users)
	db.Create(&users)
	// spew.Dump(json.Marshal(users))
	// return c.JSON(http.StatusOK, users)
	return c.JSON(http.StatusOK, users)
}
