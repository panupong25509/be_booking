package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign2/action/repositories"
)

func GetUsers(c echo.Context) error {
	// db := db.DbManager()
	// users := model.User{
	// 	Username:     "JewlyTwin",
	// 	Password:     "Jew",
	// 	Fname:        "supaiwit",
	// 	Lname:        "likitwarasad",
	// 	Organization: "SIT",
	// 	Email:        "sacxsadsa@gmail.com",
	// 	Role:         "Admin",
	// 	CreatedAt:    time.Now(),
	// 	UpdatedAt:    time.Now(),
	// }
	// db.NewRecord(users)
	// db.Create(&users)
	// spew.Dump(json.Marshal(users))
	// return c.JSON(http.StatusOK, users)

	test := repositories.EncodeJWT("bookingsign")
	log.Print(test)
	testde, _ := repositories.DecodeJWT(test, "bookingsign")
	log.Print(testde)
	return c.JSON(http.StatusOK, "users")
}
