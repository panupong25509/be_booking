package route

import (
	"github.com/labstack/echo"
	api "github.com/panupong25509/be_booking_sign/action/handles"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", api.Home)
	//Authentication
	e.POST("/login", api.Login)
	e.POST("/register", api.Register)
	//User
	e.GET("/user", api.GetUserById)
	//Sign
	e.POST("/deletesign", api.DeleteSign)
	e.POST("/updatesign", api.UpdateSign)
	e.POST("/addsign", api.AddSign)
	e.GET("/allsign", api.GetAllSign)
	// e.GET("/booking/{page}/{order}", api.GetPaginateUser)
	// e.GET("/sign/{id}", api.GetSignById)
	//Booking
	e.POST("/addbooking", api.AddBooking)
	e.GET("/booking/:id", api.GetBookingById)
	e.GET("/booking/:page/:order", api.GetBookingUser)
	e.GET("/getbookingdays/:id", api.GetBookingDayBySign)
	//admin
	e.GET("/admin/booking/:page", api.GetBookingAdmin)
	e.POST("/admin/booking/approve", api.ApproveBooking)
	e.POST("/admin/booking/reject", api.RejectBooking)

	return e
}
