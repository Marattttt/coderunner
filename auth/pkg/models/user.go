package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int
	Name       string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
