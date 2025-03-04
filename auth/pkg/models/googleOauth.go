package models

import (
	"database/sql"
	"time"
)

type GoogleOauth struct {
	UserId       int
	AccessToken  string
	RefreshToken sql.NullString
	ExpiresIn    int
	CreatedAt time.Time
}
