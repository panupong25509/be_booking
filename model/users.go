package model

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
