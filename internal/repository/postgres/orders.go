package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/dbutils"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (p *Postgres) CreateOrder(ctx context.Context, order types.Order) error {
	const insertQuery = `insert into orders (number, user_id, uploaded_at)
values (:number, :user_id, :uploaded_at)`

	if execErr := dbutils.NamedExec(ctx, p.db, insertQuery, order); execErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(execErr, &pgErr) && pgErr.Code == "23505" {
			const existQuery = `select exists(select 1 from orders where number = :number and user_id = :user_id)`

			var exists bool
			if getErr := dbutils.NamedGet(ctx, p.db, &exists, existQuery, order); getErr != nil {
				return multierr.Append(execErr, errors.Wrapf(getErr, "не удалось получить заказ %s", order.Number))
			}

			if exists {
				return errcodes.ErrOrderAlreadyUploadedByUser
			}
			return errcodes.ErrOrderAlreadyUploadedByAnotherUser
		}
		return errors.Wrapf(execErr, "не удалось создать заказ %s", order.Number)
	}

	return nil
}

func (p *Postgres) GetOrders(ctx context.Context, userID int64) (orders []types.Order, err error) {
	const q = `select number, uploaded_at
from orders
where user_id = :user_id
order by uploaded_at desc`

	args := map[string]any{
		"user_id": userID,
	}
	return orders, dbutils.NamedSelect(ctx, p.db, &orders, q, args)
}
