package service

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/config"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

type Repository interface {
	CreateUser(ctx context.Context, user types.User) (userID int64, err error)
	GetUser(ctx context.Context, login string) (user types.User, err error)
	CreateSession(ctx context.Context, session types.Session) error
	GetUserIDByTokenSession(ctx context.Context, token string) (userID int64, err error)

	CreateOrder(ctx context.Context, order types.Order) error
	GetOrders(ctx context.Context, userID int64) (orders []types.Order, err error)
}

type Service struct {
	repo   Repository
	client *resty.Client
	config *config.Config
	logger zerolog.Logger
}

func New(repo Repository, client *resty.Client, config *config.Config, logger zerolog.Logger) *Service {
	s := &Service{
		repo:   repo,
		client: client,
		config: config,
		logger: logger,
	}
	return s
}
