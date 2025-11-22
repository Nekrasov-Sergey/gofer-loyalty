package service

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/config"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

type Repository interface {
	CreateUser(ctx context.Context, user types.User) (userID int64, err error)
	GetUser(ctx context.Context, login string) (user types.User, err error)
	CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error
	GetUserIDBySession(ctx context.Context, token string) (userID int64, err error)
}

type Service struct {
	repo   Repository
	config *config.Config
	logger zerolog.Logger
}

func New(repo Repository, config *config.Config, logger zerolog.Logger) *Service {
	s := &Service{
		repo:   repo,
		config: config,
		logger: logger,
	}
	return s
}
