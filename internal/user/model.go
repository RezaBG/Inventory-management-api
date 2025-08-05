package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}
