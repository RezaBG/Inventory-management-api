package user

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint
	User      User
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}
