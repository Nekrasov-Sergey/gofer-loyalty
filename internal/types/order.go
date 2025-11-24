package types

import (
	"time"
)

type Order struct {
	ID         int64     `db:"id"`
	Number     string    `db:"number"`
	UserID     int64     `db:"user_id"`
	UploadedAt time.Time `db:"uploaded_at"`
}
