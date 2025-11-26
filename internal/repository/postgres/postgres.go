package postgres

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/service"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/dbutils"
)

type Postgres struct {
	db     sqlx.ExtContext
	rawDB  *sqlx.DB
	logger zerolog.Logger
}

func New(databaseDSN string, logger zerolog.Logger) (*Postgres, error) {
	if err := migrateDB(databaseDSN, logger); err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("pgx", databaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось подключиться к БД")
	}

	logger.Info().Msg("Установлено подключение к БД")

	return &Postgres{
		db:     db,
		rawDB:  db,
		logger: logger,
	}, nil
}

func (p *Postgres) Close() error {
	if err := p.rawDB.Close(); err != nil {
		return errors.Wrap(err, "не удалось закрыть соединения с БД")
	}
	p.logger.Info().Msg("Закрыто соединение с БД")
	return nil
}

func (p *Postgres) WithTx(ctx context.Context, fn func(txRepo service.Repository) error) error {
	return dbutils.WrapTxx(ctx, p.rawDB, nil, func(tx *sqlx.Tx) error {
		txRepo := &Postgres{
			db:     tx,
			rawDB:  nil,
			logger: p.logger.With().Str("scope", "tx").Logger(),
		}
		return fn(txRepo)
	})
}
