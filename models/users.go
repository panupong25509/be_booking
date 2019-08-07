package models

import (
	"time"

	"github.com/gofrs/uuid"
)

func (u User) TableName() string {
	return "users"
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Password     string    `json:"password" db:"password"`
	Fname        string    `json:"fname" db:"fname"`
	Lname        string    `json:"lname" db:"lname"`
	Organization string    `json:"organization" db:"organization"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Users []User

func (u *User) CheckParams(data map[string]interface{}) bool {
	if data["username"] == nil {
		return false
	}
	if data["password"] == nil {
		return false
	}
	if data["fname"] == nil {
		return false
	}
	if data["lname"] == nil {
		return false
	}
	if data["organization"] == nil {
		return false
	}
	if data["email"] == nil {
		return false
	}
	if data["role"] == nil {
		return false
	}
	return true
}

func (u *User) CreateModel(data map[string]interface{}, password string) bool {
	u.ID, _ = uuid.NewV4()
	u.Username = data["username"].(string)
	u.Password = password
	u.Fname = data["fname"].(string)
	u.Lname = data["lname"].(string)
	u.Organization = data["organization"].(string)
	u.Email = data["email"].(string)
	u.Role = data["role"].(string)
	return true
}
