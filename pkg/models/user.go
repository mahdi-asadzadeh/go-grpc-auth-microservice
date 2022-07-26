package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"column:email;UNIQUE; not null"`
	Password string `gorm:"column:password;not null"`
}
