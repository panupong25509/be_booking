package models

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

func (b Booking) TableName() string {
	return "booking"
}

type Booking struct {
	ID          int       `json:"id" db:"id"`
	Code        string    `json:"booking_code" db:"booking_code"`
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

func (b *Booking) CreateModel(data map[string]interface{}, code string) bool {
	if data["applicant_id"] == nil {
		return false
	}
	if data["sign_id"] == nil {
		return false
	}
	if data["description"] == nil {
		return false
	}
	if data["first_date"] == nil {
		return false
	}
	if data["last_date"] == nil {
		return false
	}
	b.Code = code
	b.ApplicantID, _ = uuid.FromString(data["applicant_id"].(string))
	b.SignID, _ = strconv.Atoi(data["sign_id"].(string))
	b.Description = data["description"].(string)
	b.FirstDate, _ = time.Parse("2006-01-02", data["first_date"].(string))
	b.LastDate, _ = time.Parse("2006-01-02", data["last_date"].(string))
	b.Status = "pending"
	b.Comment = ""
	return true
}
