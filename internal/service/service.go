package service

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/config"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(txRepo Repository) error) error

	CreateUser(ctx context.Context, user *types.User) (userID int64, err error)
	GetUserByLogin(ctx context.Context, login string) (user *types.User, err error)
	CreateSession(ctx context.Context, session *types.Session) error
	GetUserIDByToken(ctx context.Context, sessionToken string) (userID int64, err error)

	CreateOrder(ctx context.Context, order *types.Order) error
	GetOrdersByUserID(ctx context.Context, userID int64) (orders []types.Order, err error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*types.Order, error)
	GetOrdersForProcessing(ctx context.Context, limit int) (orders []types.Order, err error)
	UpdateOrder(ctx context.Context, order *types.Order) error

	GetBalanceByUserID(ctx context.Context, userID int64) (balance *types.Balance, err error)
	CreateBalance(ctx context.Context, userID int64) error
	AddBalance(ctx context.Context, userID int64, balance decimal.Decimal) error
	WithdrawBalance(ctx context.Context, withdrawal *types.Withdrawal) error
	AddWithdrawal(ctx context.Context, withdrawal *types.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID int64) (withdrawals []types.Withdrawal, err error)
}

type Service struct {
	repo   Repository
	client *resty.Client
	config *config.Config
	logger zerolog.Logger

	globalPauseMu *sync.Mutex
	globalPauseCh chan struct{}
}

func New(repo Repository, client *resty.Client, config *config.Config, logger zerolog.Logger) *Service {
	s := &Service{
		repo:   repo,
		client: client,
		config: config,
		logger: logger,

		globalPauseMu: &sync.Mutex{},
		globalPauseCh: make(chan struct{}),
	}
	close(s.globalPauseCh)
	return s
}
