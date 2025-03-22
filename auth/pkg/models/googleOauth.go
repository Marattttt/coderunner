package models

import (
	"time"
)

type GoogleOauth struct {
	ID           int `gorm:"primaryKey"`
	UserID       int `json:"-"`
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}
