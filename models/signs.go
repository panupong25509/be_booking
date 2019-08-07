package models

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
)

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

func (s *Sign) CheckParamPostForm(c echo.Context) bool {
	if c.FormValue("signname") == "" {
		return false
	}
	if c.FormValue("location") == "" {
		return false
	}
	if c.FormValue("limitdate") == "" {
		return false
	}
	if c.FormValue("beforebooking") == "" {
		return false
	}
	return true
}

func (s *Sign) CreateSignModel(c echo.Context) {
	if c.FormValue("id") != "" {
		s.ID, _ = strconv.Atoi(c.FormValue("id"))
	}
	s.Name = c.FormValue("signname")
	s.Location = c.FormValue("location")
	s.Limitdate, _ = strconv.Atoi(c.FormValue("limitdate"))
	s.Beforebooking, _ = strconv.Atoi(c.FormValue("beforebooking"))
	s.Picture = c.FormValue("signname")
}
