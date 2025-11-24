package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/dbutils"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (p *Postgres) CreateUser(ctx context.Context, user types.User) (userID int64, err error) {
	const q = `insert into users (login, password)
values (:login, :password)
returning id`

	if err := dbutils.NamedGet(ctx, p.db, &userID, q, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, errcodes.ErrLoginAlreadyExists
		}
		return 0, errors.Wrapf(err, "не удалось создать пользователя %s", user.Login)
	}

	return userID, nil
}

func (p *Postgres) GetUser(ctx context.Context, login string) (user types.User, err error) {
	const q = `select id, login, password
from users
where login = :login`

	args := map[string]any{
		"login": login,
	}

	if err := dbutils.NamedGet(ctx, p.db, &user, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, errcodes.ErrInvalidCredentials
		}
		return types.User{}, err
	}

	return user, nil
}

func (p *Postgres) CreateSession(ctx context.Context, session types.Session) error {
	const q = `insert into sessions (token, user_id, expires_at)
values (:token, :user_id, :expires_at)`

	if err := dbutils.NamedExec(ctx, p.db, q, session); err != nil {
		return errors.Wrap(err, "не удалось создать сессию")
	}
	return nil
}

func (p *Postgres) GetUserIDByTokenSession(ctx context.Context, token string) (userID int64, err error) {
	const q = `select user_id
from sessions
where token = :token
  and expires_at > now()`

	args := map[string]any{
		"token": token,
	}

	if err := dbutils.NamedGet(ctx, p.db, &userID, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errcodes.ErrUserUnauthorized
		}
		return 0, err
	}

	return userID, nil
}
