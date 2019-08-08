package handles

import (
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/action/repositories"
	"github.com/panupong25509/be_booking_sign/models"
)

func AddBooking(c echo.Context) error {
	newBooking, err := repositories.AddBooking(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	booking := newBooking.(models.Booking)
	return c.JSON(200, booking.ReturnJsonID())
}

// func DeleteBooking(c echo.Context) error {
// 	data := DynamicPostForm(c)

// 	deleteBooking, err := repositories.DeleteBooking(c, data)
// 	if err != nil {
// 		status := err.(models.Error)
// 		return c.Render(status.Code, r.JSON(status))
// 	}
// 	return c.Render(200, r.JSON(deleteBooking))
// }

func GetBookingDayBySign(c echo.Context) error {
	days, err := repositories.GetBookingDaysBySign(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, days)
}

func ApproveBooking(c echo.Context) error {
	message, err := repositories.ApproveBooking(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, message)
}

func RejectBooking(c echo.Context) error {
	message, err := repositories.RejectBooking(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, message)
}

func GetBookingAdmin(c echo.Context) error {
	booking, _ := repositories.GetBookingAdmin(c)
	// if err != nil {
	// 	status := err.(models.Error)
	// 	return c.JSON(status.Code, status)
	// }
	return c.JSON(200, booking)
}

func GetBookingUser(c echo.Context) error {
	booking, err := repositories.GetBookingUser(c)
	if err != nil {
		status := err.(models.Error)
		return c.JSON(status.Code, status)
	}
	return c.JSON(200, booking)
}

// func GetPaginateUser(c echo.Context) error {
// 	page := c.Param("page")
// 	order := c.Param("order")
// 	allBooking, err := repositories.GetPaginateUser(page, order, c)
// 	if err != nil {
// 		status := err.(models.Error)
// 		return c.JSON(status.Code, status)
// 	}
// 	return c.JSON(200, allBooking)
// }
