package types

import (
	"time"
)

type Session struct {
	Token     string    `json:"token" db:"token"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}
