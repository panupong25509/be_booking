package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
