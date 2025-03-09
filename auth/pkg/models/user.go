package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email string
	CreatedAt time.Time
	DeletedAt sql.NullTime
	Google    *GoogleOauth
}
