package models

import "time"

func (s Sign) TableName() string {
	return "signs"
}

type Sign struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"sign_name"`
	Location      string    `json:"location" db:"location"`
	Limitdate     int       `json:"limitdate" db:"limitdate"`
	Beforebooking int       `json:"beforebooking" db:"beforebooking"`
	Picture       string    `json:"picture" db:"picture"`
	Booking       []Booking `json:"-" db:"-"  has_many:"bookings"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"-" db:"updated_at"`
}
