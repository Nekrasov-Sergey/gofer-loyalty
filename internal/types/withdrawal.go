package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type Withdrawal struct {
	ID          int64           `db:"id"`
	OrderNumber string          `json:"order" db:"order_number"`
	Withdrawn   decimal.Decimal `json:"sum" db:"withdrawn"`
	UserID      int64           `db:"user_id"`
	ProcessedAt time.Time       `db:"processed_at"`
}
