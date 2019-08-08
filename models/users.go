package models

import (
	"time"

	"github.com/labstack/echo"

	"github.com/gofrs/uuid"
)

func (u User) TableName() string {
	return "users"
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"-" db:"username"`
	Password     string    `json:"-" db:"password"`
	Fname        string    `json:"fname" db:"fname"`
	Lname        string    `json:"lname" db:"lname"`
	Organization string    `json:"organization" db:"organization"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Users []User

func (u *User) CheckParams(c echo.Context) bool {
	if c.FormValue("username") == "" {
		return false
	}
	if c.FormValue("password") == "" {
		return false
	}
	if c.FormValue("fname") == "" {
		return false
	}
	if c.FormValue("lname") == "" {
		return false
	}
	if c.FormValue("organization") == "" {
		return false
	}
	if c.FormValue("email") == "" {
		return false
	}
	if c.FormValue("role") == "" {
		return false
	}
	return true
}

func (u *User) CreateModel(c echo.Context, password string) bool {
	u.ID, _ = uuid.NewV4()
	u.Username = c.FormValue("username")
	u.Password = password
	u.Fname = c.FormValue("fname")
	u.Lname = c.FormValue("lname")
	u.Organization = c.FormValue("organization")
	u.Email = c.FormValue("email")
	u.Role = c.FormValue("role")
	return true
}
