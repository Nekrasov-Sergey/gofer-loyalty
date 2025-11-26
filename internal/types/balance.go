package types

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	ID        int64           `db:"id"`
	Current   decimal.Decimal `db:"current"`
	Withdrawn decimal.Decimal `db:"withdrawn"`
	UserID    int64           `db:"user_id"`
}
