package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID         int64           `db:"id"`
	Number     string          `db:"number"`
	Status     OrderStatus     `db:"status"`
	Accrual    decimal.Decimal `db:"accrual"`
	UserID     int64           `db:"user_id"`
	UploadedAt time.Time       `db:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
