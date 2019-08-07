package repositories

import (
	"strconv"
	"time"

	"github.com/JewlyTwin/be_booking_sign/mailers"
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
)

func AddBooking(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	signID, _ := strconv.Atoi(data["sign_id"].(string))
	sign, err := GetSignByID(c, data)
	if err != nil {
		return nil, err
	}
	code := GenCodeBooking(data, sign.(models.Sign))
	newBooking := models.Booking{}
	if !newBooking.CreateModel(data, code) {
		return nil, models.Error{400, "Please complete all fields"}
	}
	validate, err := ValidateBookingTime(newBooking, sign.(models.Sign))
	if err != nil {
		return nil, err
	}
	if !validate {
		return nil, models.Error{400, "Busy date"}
	}
	err = db.Create(&newBooking)
	if err != nil {
		return nil, models.Error{500, "Can't Create to Database"}
	}
	return newBooking, nil
}

func GenCodeBooking(data map[string]interface{}, sign models.Sign) string {
	code := sign.Name + "CODE" + data["first_date"].(string) + data["last_date"].(string)
	return code
}

func ValidateBookingTime(newBooking models.Booking, sign models.Sign) (bool, interface{}) {
	db := db.DbManager()
	bookings := models.Bookings{}
	db.Where("last_date >= (?) and first_date <= (?) and sign_id = (?)", newBooking.FirstDate, newBooking.LastDate, newBooking.SignID).Find(&bookings)
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
	signid, _ := strconv.Atoi(c.Param("id"))
	err := db.Where("( last_date >= (?) or first_date >= (?) ) and sign_id = (?)", bookingdate, bookingdate, signid).Find(&bookings)
	if err != nil {
		return nil, models.Error{400, "DB"}
	}
	days := models.BookingDays{}
	for _, value := range bookings {
		days = append(days, models.BookingDay{value.FirstDate, value.LastDate})
	}
	return days, nil
}

func RejectBooking(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
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
	idbooking, _ := strconv.Atoi(data["id"].(string))
	comment := data["comment"].(string)
	booking := models.Booking{}
	err = db.Find(&booking, idbooking)
	if err != nil {
		return nil, models.Error{500, "Can't Select data form Database"}
	}
	if booking.Status == "reject" {
		return nil, models.Error{500, "This booking id is status reject"}
	}
	if booking.Status == "approve" {
		return nil, models.Error{500, "This booking id is status approve"}
	}
	booking.Comment = comment
	booking.Status = "reject"
	// userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
	// user, err := userInterface.(models.User)
	db.Update(&booking)
	// mailers.SendWelcomeEmails(c, "Your booking Rejected", user.Email, false)
	return models.Success{200, "success"}, nil
}

func ApproveBooking(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
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
	id, _ := strconv.Atoi(data["id"].(string))
	booking := models.Booking{}
	err = db.Find(&booking, id)
	if err != nil {
		return nil, err
	}
	if booking.Status == "approve" {
		return nil, models.Error{500, "This booking approved"}
	}
	if booking.Status == "reject" {
		return nil, models.Error{500, "This booking rejected"}
	}
	booking.Status = "approve"
	// userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
	// user, err := userInterface.(models.User)
	err = db.Update(&booking)
	if err != nil {
		return nil, err
	}
	// mailers.SendWelcomeEmails(c, "Your booking approved", user.Email, true)
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

func SendMail(c echo.Context) (interface{}, interface{}) {
	mailers.SendWelcomeEmails(c.Response(), "Test", "panupong.jkn@gmail.com", true)
	return nil, nil
}

func GetPaginateAdmin(page string, c echo.Context) (interface{}, interface{}) {
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

	numberPage, _ := strconv.Atoi(page)
	q := db.Paginate(numberPage, 10)
	booking := []models.Booking{}
	err = q.Where("status = 'pending'").All(&booking)

	bookings := []models.Booking{}
	for _, value := range booking {
		user, err := GetUserByIduuid(c, value.ApplicantID)
		if err != nil {
			return nil, err
		}
		value.Applicant = user.(models.User)
		sign, err := GetSignByID(c, value.SignID)
		if err != nil {
			return nil, err
		}
		value.Sign = sign.(models.Sign)
		bookings = append(bookings, value)
	}
	bookingJson := models.Page{numberPage, bookings, q.Paginator.TotalPages}
	return &bookingJson, nil
}

func GetPaginateUser(page string, order string, c echo.Context) (interface{}, interface{}) {
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string), "bookingsign")
	if err != nil {
		return nil, err
	}
	db := db.DbManager()
	numberPage, _ := strconv.Atoi(page)
	q := db.Paginate(numberPage, 10)
	booking := []models.Booking{}
	err = q.Where("applicant_id = (?)", tokens["UserID"]).Order(order).All(&booking)
	bookings := []models.Booking{}
	for _, value := range booking {
		user, err := GetUserByIduuid(c, value.ApplicantID)
		if err != nil {
			return nil, err
		}
		value.Applicant = user.(models.User)
		sign, err := GetSignById(c, value.SignID)
		if err != nil {
			return nil, err
		}
		value.Sign = sign.(models.Sign)
		bookings = append(bookings, value)
	}
	bookingJson := models.Page{numberPage, bookings, q.Paginator.TotalPages}
	return &bookingJson, nil

}
