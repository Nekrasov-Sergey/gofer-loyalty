package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/dbutils"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (p *Postgres) CreateOrder(ctx context.Context, order *types.Order) error {
	const q = `insert into orders (number, status, accrual, user_id, uploaded_at)
values (:number, :status, :accrual, :user_id, :uploaded_at)`

	return dbutils.NamedExec(ctx, p.db, q, order)
}

func (p *Postgres) GetOrdersByUserID(ctx context.Context, userID int64) (orders []types.Order, err error) {
	const q = `select id, number, status, accrual, user_id, uploaded_at
from orders
where user_id = :user_id
order by uploaded_at desc`

	args := map[string]any{
		"user_id": userID,
	}
	return orders, dbutils.NamedSelect(ctx, p.db, &orders, q, args)
}

func (p *Postgres) GetOrderByNumber(ctx context.Context, orderNumber string) (order *types.Order, err error) {
	const q = `select id, number, status, accrual, user_id, uploaded_at
from orders
where number = :number`

	args := map[string]any{
		"number": orderNumber,
	}

	order = &types.Order{}
	if err := dbutils.NamedGet(ctx, p.db, order, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

func (p *Postgres) GetOrdersForProcessing(ctx context.Context, limit int) (orders []types.Order, err error) {
	const q = `select id, number, status, accrual, user_id, uploaded_at
from orders
where status in ('NEW', 'PROCESSING')
order by uploaded_at
limit :limit`

	args := map[string]any{
		"limit": limit,
	}
	return orders, dbutils.NamedSelect(ctx, p.db, &orders, q, args)
}

func (p *Postgres) UpdateOrder(ctx context.Context, order *types.Order) error {
	const q = `update orders
set status  = :status,
    accrual = :accrual
where id = :id`

	return dbutils.NamedExec(ctx, p.db, q, order)
}
