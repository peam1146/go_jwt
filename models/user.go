package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"unique"`
	Password []byte `json:"password" binding:"required"`
}
