package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Email     string
	CreatedAt time.Time
	DeletedAt sql.NullTime
	Google    *GoogleOauth
}
