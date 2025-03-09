package models

import (
	"time"
)

type GoogleOauth struct {
	UserId       int
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}
