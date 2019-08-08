package models

import (
	"strconv"
	"time"

	"github.com/labstack/echo"

	"github.com/gofrs/uuid"
)

func (b Booking) TableName() string {
	return "bookings"
}

type Paginator struct {
	Allpage  int       `json:"allpage"`
	Bookings []Booking `json:"bookings"`
}

type Booking struct {
	ID          int       `json:"id" db:"id"`
	Code        string    `json:"booking_code" db:"code"`
	ApplicantID uuid.UUID `json:"applicant_id" db:"applicant_id" fk_id:"id"`
	SignID      int       `json:"sign_id" db:"sign_id" fk_id:"id"`
	Description string    `json:"description" db:"description"`
	FirstDate   time.Time `json:"first_date" db:"first_date"`
	LastDate    time.Time `json:"last_date" db:"last_date"`
	Status      string    `json:"status" db:"status"`
	Comment     string    `json:"comment" db:"comment"`
	Applicant   User      `json:"applicant" db:"-"`
	Sign        Sign      `json:"sign" db:"-"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type IDbooking struct {
	ID int `json:"id"`
}

func (b *Booking) ReturnJsonID() IDbooking {
	idbook := IDbooking{b.ID}
	return idbook
}

type BookingDay struct {
	Firstdate time.Time `json:"firstdate"`
	Lastdate  time.Time `json:"lastdate"`
}

type Page struct {
	NumberPage int
	Bookings   []Booking
	TotalPage  int
}

type Bookings []Booking

type BookingDays []BookingDay

func (b *Booking) CreateModel(c echo.Context, code string) bool {
	if c.FormValue("applicant_id") == "" {
		return false
	}
	if c.FormValue("sign_id") == "" {
		return false
	}
	if c.FormValue("description") == "" {
		return false
	}
	if c.FormValue("first_date") == "" {
		return false
	}
	if c.FormValue("last_date") == "" {
		return false
	}
	b.Code = code
	b.ApplicantID, _ = uuid.FromString(c.FormValue("applicant_id"))
	b.SignID, _ = strconv.Atoi(c.FormValue("sign_id"))
	b.Description = c.FormValue("description")
	b.FirstDate, _ = time.Parse("2006-01-02", c.FormValue("first_date"))
	b.LastDate, _ = time.Parse("2006-01-02", c.FormValue("last_date"))
	b.Status = "pending"
	b.Comment = ""
	return true
}
