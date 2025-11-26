package types

import (
	"github.com/shopspring/decimal"
)

type ExternalOrder struct {
	Order   string              `json:"order"`
	Status  ExternalOrderStatus `json:"status"`
	Accrual decimal.Decimal     `json:"accrual"`
}

type ExternalOrderStatus string

const (
	ExternalOrderStatusEmpty      ExternalOrderStatus = ""
	ExternalOrderStatusRegistered ExternalOrderStatus = "REGISTERED"
	ExternalOrderStatusProcessing ExternalOrderStatus = "PROCESSING"
	ExternalOrderStatusInvalid    ExternalOrderStatus = "INVALID"
	ExternalOrderStatusProcessed  ExternalOrderStatus = "PROCESSED"
)
