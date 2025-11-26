package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/dbutils"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (p *Postgres) GetBalanceByUserID(ctx context.Context, userID int64) (balance *types.Balance, err error) {
	const q = `select id, current, withdrawn, user_id
from balances
where user_id = :user_id`

	args := map[string]any{
		"user_id": userID,
	}

	balance = &types.Balance{}
	return balance, dbutils.NamedGet(ctx, p.db, balance, q, args)
}

func (p *Postgres) CreateBalance(ctx context.Context, userID int64) error {
	const q = `insert into balances (current, withdrawn, user_id)
values (0, 0, :user_id)`

	args := map[string]any{
		"user_id": userID,
	}

	return dbutils.NamedExec(ctx, p.db, q, args)
}

func (p *Postgres) AddBalance(ctx context.Context, userID int64, balance decimal.Decimal) error {
	const q = `update balances
set current = current + :balance
where user_id = :user_id`

	args := map[string]any{
		"user_id": userID,
		"balance": balance,
	}

	return dbutils.NamedExec(ctx, p.db, q, args)
}

func (p *Postgres) WithdrawBalance(ctx context.Context, withdrawal *types.Withdrawal) error {
	const q = `update balances
set current   = current - :withdrawn,
    withdrawn = withdrawn + :withdrawn
where user_id = :user_id
  and current >= :withdrawn
returning current`

	var current decimal.Decimal
	if err := dbutils.NamedGet(ctx, p.db, &current, q, withdrawal); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errcodes.ErrNotEnoughBalance
		}
		return err
	}
	return nil
}

func (p *Postgres) AddWithdrawal(ctx context.Context, withdrawal *types.Withdrawal) error {
	const q = `insert into withdrawals (order_number, withdrawn, user_id, processed_at)
values (:order_number, :withdrawn, :user_id, :processed_at)`

	return dbutils.NamedExec(ctx, p.db, q, withdrawal)
}

func (p *Postgres) GetWithdrawalsByUserID(ctx context.Context, userID int64) (withdrawals []types.Withdrawal, err error) {
	const q = `select id, order_number, withdrawn, user_id, processed_at
from withdrawals
where user_id = :user_id
order by processed_at desc`

	args := map[string]any{
		"user_id": userID,
	}
	return withdrawals, dbutils.NamedSelect(ctx, p.db, &withdrawals, q, args)
}
