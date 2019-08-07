package repositories

import (
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/JewlyTwin/be_booking_sign/mailers"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
	"github.com/siredwin/pongorenderer/renderer"
)

func AddBooking(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	sign, err := GetSignByID(c, data)
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
	err = db.Create(&newBooking)
	if err != nil {
		return nil, models.Error{500, "Can't Create to Database"}
	}
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

func GetPaginateAdmin(c echo.Context) error {
	db := db.DbManager()
	// Lets use the Forbes top 7.
	booking := []models.Booking{}
	db.Where("status = 'pending'").Find(&booking)
	log.Print(booking)
	// usernames := []string{"Larry Ellison", "Carlos Slim Helu",
	// "Mark Zuckerberg", "Amancio Ortega ", "Jeff Bezos", " Warren Buffett ", "Bill Gates"}
	// sets paginator with the current offset (from the url query param)
	postsPerPage := 2
	paginator = pagination.NewPaginator(c.Request(), postsPerPage, len(booking))
	// fetch the next posts "postsPerPage"
	idrange := NewSlice(paginator.Offset(), postsPerPage, 1)
	//create a new page list that shows up on html
	Bookings := []models.Booking{}
	for _, num := range idrange {
		//Prevent index out of range errors
		if num <= len(booking)-1 {
			myuser := booking[num]
			Bookings = append(Bookings, myuser)
		}
	}
	// set the paginator in context
	// also set the page list in context
	// if you also have more data, set it context
	data = pongo2.Context{"paginator": paginator, "posts": Bookings}
	log.Print(data)
	return c.Render(http.StatusOK, "templates/page.html", data)
}

// func GetPaginateAdmin(page string, c echo.Context) (interface{}, interface{}) {
// 	jwtReq, err := GetJWT(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tokens, err := DecodeJWT(jwtReq.(string))
// 	if err != nil {
// 		return nil, err
// 	}
// 	db := db.DbManager()
// 	if tokens["Role"] != "admin" {
// 		return nil, models.Error{500, "You not Admin"}
// 	}

// 	numberPage, _ := strconv.Atoi(page)
// 	// db.Offset(100).Limit(20).Find(&interface, "id = ?", id)
// 	q := db.Paginate(numberPage, 10)
// 	booking := []models.Booking{}
// 	err = q.Where("status = 'pending'").All(&booking)

// 	bookings := []models.Booking{}
// 	for _, value := range booking {
// 		user, err := GetUserByIduuid(c, value.ApplicantID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		value.Applicant = user.(models.User)
// 		sign, err := GetSignByID(c, value.SignID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		value.Sign = sign.(models.Sign)
// 		bookings = append(bookings, value)
// 	}
// 	bookingJson := models.Page{numberPage, bookings, q.Paginator.TotalPages}
// 	return &bookingJson, nil
// }

// func GetPaginateUser(page string, order string, c echo.Context) (interface{}, interface{}) {
// 	jwtReq, err := GetJWT(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tokens, err := DecodeJWT(jwtReq.(string))
// 	if err != nil {
// 		return nil, err
// 	}
// 	db := db.DbManager()
// 	numberPage, _ := strconv.Atoi(page)
// 	q := db.Paginate(numberPage, 10)
// 	booking := []models.Booking{}
// 	err = q.Where("applicant_id = (?)", tokens["UserID"]).Order(order).All(&booking)
// 	bookings := []models.Booking{}
// 	for _, value := range booking {
// 		user, err := GetUserByIduuid(c, value.ApplicantID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		value.Applicant = user.(models.User)
// 		sign, err := GetSignById(c, value.SignID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		value.Sign = sign.(models.Sign)
// 		bookings = append(bookings, value)
// 	}
// 	bookingJson := models.Page{numberPage, bookings, q.Paginator.TotalPages}
// 	return &bookingJson, nil

// }
