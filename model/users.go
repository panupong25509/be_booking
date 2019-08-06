package model

import uuid "github.com/satori/go.uuid"

// type User struct {
// 	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
// 	Username     string    `gorm:"type:varchar(255)"`
// 	Password     string    `gorm:"type:varchar(255)"`
// 	Fname        string    `gorm:"type:varchar(255)"`
// 	Lname        string    `gorm:"type:varchar(255)"`
// 	Organization string    `gorm:"type:varchar(255)"`
// 	Email        string    `gorm:"type:varchar(255)"`
// 	Role         string    `gorm:"type:varchar(255)"`
// }

type User struct {
	ID uuid.UUID
	// Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string
}
