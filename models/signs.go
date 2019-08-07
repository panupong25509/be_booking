package models

import (
	"strconv"
	"time"
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

func (s *Sign) CheckParamPostForm(data map[string]interface{}) bool {
	if data["signname"] == nil {
		return false
	}
	if data["location"] == nil {
		return false
	}
	if data["limitdate"] == nil {
		return false
	}
	if data["beforebooking"] == nil {
		return false
	}
	return true
}

func (s *Sign) CreateSignModel(data map[string]interface{}, namepic string) {
	if data["id"] != nil {
		s.ID, _ = strconv.Atoi(data["id"].(string))
	}
	s.Name = data["signname"].(string)
	s.Location = data["location"].(string)
	s.Limitdate, _ = strconv.Atoi(data["limitdate"].(string))
	s.Beforebooking, _ = strconv.Atoi(data["beforebooking"].(string))
	s.Picture = namepic
}
