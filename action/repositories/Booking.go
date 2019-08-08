package repositories

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/mailer"
	"github.com/panupong25509/be_booking_sign/models"
	"github.com/siredwin/pongorenderer/renderer"
)

func AddBooking(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	sign, err := GetSignByID(c)
	if err != nil {
		return nil, err
	}
	code := GenCodeBooking(c, sign.(models.Sign))
	newBooking := models.Booking{}
	if !newBooking.CreateModel(c, code) {
		return nil, models.Error{400, "Please complete all fields"}
	}
	validate, err := ValidateBookingTime(newBooking, sign.(models.Sign))
	if err != nil {
		return nil, err
	}
	if !validate {
		return nil, models.Error{400, "Busy date"}
	}
	db.NewRecord(newBooking)
	db.Create(&newBooking)
	return newBooking, nil
}

func GenCodeBooking(c echo.Context, sign models.Sign) string {
	code := sign.Name + "CODE" + c.FormValue("first_date") + c.FormValue("last_date")
	return code
}

func ValidateBookingTime(newBooking models.Booking, sign models.Sign) (bool, interface{}) {
	db := db.DbManager()
	bookings := models.Bookings{}
	db.Where("last_date >= (?) and first_date <= (?) and sign_id = (?)",
		newBooking.FirstDate, newBooking.LastDate, newBooking.SignID).Find(&bookings)
	if len(bookings) != 0 {
		return false, models.Error{500, "Busy date"}
	}
	if CheckDate(newBooking.FirstDate, newBooking.LastDate) > sign.Limitdate {
		return false, models.Error{500, "Please book within " + strconv.Itoa(sign.Limitdate) + " days"}
	}
	if CheckDate(time.Now(), newBooking.FirstDate) < sign.Beforebooking {
		return false, models.Error{500, "Please book before " + strconv.Itoa(sign.Beforebooking) + " days"}
	}
	return true, nil
}
func CheckDate(D1 time.Time, D2 time.Time) int {
	diff := D2.Sub(D1)
	allDay := int(diff.Hours()/24) + 1 //first-last
	day := D1
	sunday := 0
	for day.Before(D2) {
		if int(day.Weekday()) == 0 {
			sunday = sunday + 1
			day = day.AddDate(0, 0, 7)
		} else {
			day = day.AddDate(0, 0, 1)
		}
	}
	weekday := sunday * 2 // weekday in firstdate - lastdate
	return allDay - weekday
}

func GetBookingDaysBySign(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	bookings := models.Bookings{}
	bookingdate := time.Now().Format("2006-01-02")
	db.Where("( last_date >= (?) or first_date >= (?) ) and sign_id = (?)", bookingdate, bookingdate, c.Param("id")).Find(&bookings)
	days := models.BookingDays{}
	for _, value := range bookings {
		days = append(days, models.BookingDay{value.FirstDate, value.LastDate})
	}
	return days, nil
}

func RejectBooking(c echo.Context) (interface{}, interface{}) {
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string))
	if err != nil {
		return nil, err
	}
	db := db.DbManager()
	if tokens["Role"] != "admin" {
		return nil, models.Error{500, "You not Admin"}
	}
	comment := c.FormValue("comment")
	booking := models.Booking{}
	db.Find(&booking, c.FormValue("id"))
	if booking.Status == "reject" {
		return nil, models.Error{500, "This booking id is status reject"}
	}
	if booking.Status == "approve" {
		return nil, models.Error{500, "This booking id is status approve"}
	}
	booking.Comment = comment
	booking.Status = "reject"
	userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
	db.Save(&booking)
	user, err := userInterface.(models.User)
	mailer.SendEmail("Your booking Rejected", user.Email, "reject")
	return models.Success{200, "reject success"}, nil
}

func ApproveBooking(c echo.Context) (interface{}, interface{}) {
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string))
	if err != nil {
		return nil, err
	}
	db := db.DbManager()
	if tokens["Role"] != "admin" {
		return nil, models.Error{500, "You not Admin"}
	}
	booking := models.Booking{}
	db.Find(&booking, c.FormValue("id"))
	if booking.Status == "approve" {
		return nil, models.Error{500, "This booking approved"}
	}
	if booking.Status == "reject" {
		return nil, models.Error{500, "This booking rejected"}
	}
	booking.Status = "approve"
	db.Save(&booking)
	userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
	user, err := userInterface.(models.User)
	mailer.SendEmail("Your booking approved", user.Email, "approve")
	return models.Success{200, "Approve success"}, nil
}

// func DeleteBooking(c echo.Context) (interface{}, interface{}) {
// 	db, err := ConnectDB(c)
// 	if err != nil {
// 		return nil, models.Error{500, "Can't connect Database"}
// 	}
// 	data := DynamicPostForm(c)
// 	booking := models.Booking{}
// 	id, _ := strconv.Atoi(data["id"].(string))
// 	err = db.Find(&booking, id)
// 	if err != nil {
// 		return nil, models.Error{500, "Data มีปัญหาไม่สามารถยกเลิกได้"}
// 	}
// 	_ = db.Destroy(&booking)
// 	return Success(nil), nil
// }

// func SendMail(c echo.Context) (interface{}, interface{}) {
// 	mailers.SendWelcomeEmails(c.Response(), "Test", "panupong.jkn@gmail.com", true)
// 	return nil, nil
// }

var (
	paginator    = &pagination.Paginator{}
	data         = pongo2.Context{}
	MainRenderer = renderer.Renderer{Debug: true}
)

//generator
func NewSlice(start, count, step int) []int {
	s := make([]int, count)
	for i := range s {
		s[i] = start
		start += step
	}
	return s
}

func GetPaginateAdmin(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string))
	if err != nil {
		return nil, err
	}
	if tokens["Role"] != "admin" {
		return nil, models.Error{500, "You not Admin"}
	}
	booking := []models.Booking{}
	db.Where("status = 'pending'").Find(&booking)
	postsPerPage := 3
	pagestring, _ := strconv.Atoi(c.Param("page"))
	pageint := pagestring - 1
	paginator = pagination.NewPaginator(c.Request(), postsPerPage, len(booking))
	idrange := NewSlice((pageint * postsPerPage), postsPerPage, 1)
	Bookings := []models.Booking{}
	for _, num := range idrange {
		if num <= len(booking)-1 {
			myuser := booking[num]
			user, err := GetUserByIduuid(c, myuser.ApplicantID)
			if err != nil {
				return nil, err
			}
			myuser.Applicant = user.(models.User)
			sign, err := GetSignByIDForPage(myuser.SignID)
			if err != nil {
				return nil, err
			}
			myuser.Sign = sign.(models.Sign)
			Bookings = append(Bookings, myuser)
		}
	}
	allpage := int((len(booking) / postsPerPage))
	if len(booking)%postsPerPage != 0 {
		allpage++
	}
	Pagination := models.Paginator{int(allpage), Bookings}
	return Pagination, nil
}

func GetPaginateUser(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string))
	if err != nil {
		return nil, err
	}
	booking := []models.Booking{}
	db.Order(c.Param("order")).Where("applicant_id = (?)", tokens["UserID"]).Find(&booking)
	postsPerPage := 5
	pagestring, _ := strconv.Atoi(c.Param("page"))
	pageint := pagestring - 1
	paginator = pagination.NewPaginator(c.Request(), postsPerPage, len(booking))
	idrange := NewSlice((pageint * postsPerPage), postsPerPage, 1)
	Bookings := []models.Booking{}
	for _, num := range idrange {
		if num <= len(booking)-1 {
			myuser := booking[num]
			user, err := GetUserByIduuid(c, myuser.ApplicantID)
			if err != nil {
				return nil, err
			}
			myuser.Applicant = user.(models.User)
			sign, err := GetSignByIDForPage(myuser.SignID)
			if err != nil {
				return nil, err
			}
			myuser.Sign = sign.(models.Sign)
			Bookings = append(Bookings, myuser)
		}
	}
	allpage := int((len(booking) / postsPerPage))
	if len(booking)%postsPerPage != 0 {
		allpage++
	}
	Pagination := models.Paginator{int(allpage), Bookings}
	return Pagination, nil
}
