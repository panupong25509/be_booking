package repositories

import (
	"io"
	"log"
	"os"
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
	log.Print(sign)
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
	_, err = UploadImgForBooking(c, code)
	if err != nil {
		return nil, err
	}
	db.NewRecord(newBooking)
	db.Create(&newBooking)
	return newBooking, nil
}

func UploadImgForBooking(c echo.Context, code string) (interface{}, interface{}) {
	file, err := c.FormFile("file")
	if err != nil {
		return err, nil
	}
	src, err := file.Open()
	if err != nil {
		return err, nil
	}
	defer src.Close()
	dst, err := os.Create(`D:\fe_booking_sign\public\img\applicant_poster\` + code + `.jpg`)
	if err != nil {
		return err, nil
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err, nil
	}

	return nil, nil
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

func GetBookingById(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	Booking := []models.Booking{}
	db.First(&Booking, c.Param("id"))
	user, err := GetUserByIduuid(c, Booking[0].ApplicantID)
	if err != nil {
		return nil, err
	}
	Booking[0].Applicant = user.(models.User)
	sign, err := GetSignByIDForPage(Booking[0].SignID)
	if err != nil {
		return nil, err
	}
	Booking[0].Sign = sign.(models.Sign)
	return &Booking, nil
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

func GetAllBookingThisYear() models.Bookings {
	db := db.DbManager()
	bookings := models.Bookings{}
	db.Where("extract(year from created_at) = (?)", time.Now().Year()).Find(&bookings)
	log.Print(bookings)
	return bookings
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

func Paginator(c echo.Context, allrequest, allbookingthisyear, page, count int, booking []models.Booking) (interface{}, interface{}) {
	pageint := page - 1
	paginator = pagination.NewPaginator(c.Request(), count, len(booking))
	idrange := NewSlice((pageint * count), count, 1)
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
	allpage := int((len(booking) / count))
	if len(booking)%count != 0 {
		allpage++
	}
	Pagination := models.Paginator{allrequest, allbookingthisyear, int(allpage), Bookings}
	return Pagination, nil
}

func GetBookingAdmin(c echo.Context) (interface{}, interface{}) {
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

	bookingsThisYear := GetAllBookingThisYear()
	page, _ := strconv.Atoi(c.Param("page"))
	Pagination, err := Paginator(c, len(booking), len(bookingsThisYear), page, 10, booking)
	if err != nil {
		return nil, err
	}
	return Pagination, nil
}

func GetBookingUser(c echo.Context) (interface{}, interface{}) {
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

	page, _ := strconv.Atoi(c.Param("page"))
	Pagination, err := Paginator(c, 0, 0, page, 10, booking)
	if err != nil {
		return nil, err
	}
	return Pagination, nil
}

func GetBookingByFilter(c echo.Context) (interface{}, interface{}) {
	bookings := []models.Booking{}
	db := db.DbManager()
	db = db.Find(&bookings)
	log.Print(c.Param("month"))
	log.Print(c.Param("year"))
	log.Print(c.Param("signid"))
	log.Print(c.Param("organization"))
	if c.Param("month") != "null" {
		db = db.Where("extract(month from created_at) = (?)", c.Param("month")).Find(&bookings)
	}
	if c.Param("year") != "null" {
		db = db.Where("extract(year from created_at) = (?)", c.Param("year")).Find(&bookings)
	}
	if c.Param("signid") != "null" {
		db = db.Where("sign_id = (?)", c.Param("signid")).Find(&bookings)
	}
	bookingsOrganization := []models.Booking{}
	if c.Param("organization") != "null" {
		for index, booking := range bookings {
			log.Print(index, booking)
			userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
			if err != nil {
				return nil, err
			}
			user := userInterface.(models.User)
			if user.Organization == c.Param("organization") {
				booking.Applicant = user
				sign, err := GetSignByIDForPage(booking.SignID)
				if err != nil {
					return nil, err
				}
				booking.Sign = sign.(models.Sign)
				bookingsOrganization = append(bookingsOrganization, booking)
			}
		}
	} else {
		for index, booking := range bookings {
			log.Print(index, booking)
			userInterface, err := GetUserByIduuid(c, booking.ApplicantID)
			if err != nil {
				return nil, err
			}
			user := userInterface.(models.User)
			booking.Applicant = user
			sign, err := GetSignByIDForPage(booking.SignID)
			if err != nil {
				return nil, err
			}
			booking.Sign = sign.(models.Sign)
			bookingsOrganization = append(bookingsOrganization, booking)
		}
	}
	pagenumber, _ := strconv.Atoi(c.Param("page"))
	Paginator, err := Paginator(c, 0, 0, pagenumber, 5, bookingsOrganization)
	if err != nil {
		return nil, err
	}
	return Paginator, nil
}
func GetSummaryMonth(c echo.Context) (interface{}, interface{}) {
	type Summary struct {
		Month        float64
		Sign         string
		Organization string
		Total        int
	}
	type Summarys struct {
		Summarys []Summary `json:"summarys"`
		Total    int       `json:"total"`
	}
	db := db.DbManager()
	summarys := Summarys{}
	selectSQL := "extract(month from bookings.created_at) as month,signs.name as sign,  users.organization as organization, count(organization) as total"
	joinUser := "join users on bookings.applicant_id = users.id"
	joinSign := "join signs on bookings.sign_id = signs.id"
	whereMonth := "extract(month from bookings.created_at) = " + c.Param("month")
	if c.Param("month") == "null" {
		whereMonth = "extract(month from bookings.created_at) != 0"
	}
	whereSign := "signs.name = (?)"
	if c.Param("sign") == "null" {
		whereSign = "signs.name NOT LIKE (?)"
	}
	whereOrganization := "users.organization = (?)"
	if c.Param("organization") == "null" {
		whereOrganization = "users.organization NOT LIKE (?)"
	}
	log.Print(whereSign)
	groupby := "month,sign, organization"
	db.Table("bookings").Select(selectSQL).Joins(joinUser).Joins(joinSign).Where(whereMonth).Where(whereSign, c.Param("sign")).Where(whereOrganization, c.Param("organization")).Group(groupby).Scan(&summarys.Summarys)
	for _, value := range summarys.Summarys {
		summarys.Total = summarys.Total + value.Total
	}
	return summarys, nil
}
