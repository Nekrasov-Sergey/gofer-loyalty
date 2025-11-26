package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/config"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/service"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

type Service interface {
	Register(ctx context.Context, user *types.User) (sessionToken string, err error)
	Login(ctx context.Context, user *types.User) (sessionToken string, err error)

	CreateOrder(ctx context.Context, order *types.Order) error
	GetOrders(ctx context.Context, userID int64) ([]types.OrderResponse, error)

	GetBalance(ctx context.Context, userID int64) (*types.BalanceResponse, error)
	WithdrawBalance(ctx context.Context, withdrawal *types.Withdrawal) error
	GetWithdrawals(ctx context.Context, userID int64) ([]types.WithdrawalResponse, error)
}

type Handler struct {
	service Service
	config  *config.Config
	logger  zerolog.Logger
}

func New(service Service, config *config.Config, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		config:  config,
		logger:  logger,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, repo service.Repository) {
	api := r.Group("/api")
	{
		user := api.Group("/user")
		user.POST("/register", h.register)
		user.POST("/login", h.login)
		{
			authorized := user.Group("", router.SessionMiddleware(repo))
			authorized.POST("/orders", h.createOrder)
			authorized.GET("/balance", h.getBalance)
			authorized.POST("/balance/withdraw", h.withdrawBalance)
			{
				compressed := authorized.Group("", router.CompressMiddleware())
				compressed.GET("/orders", h.getOrders)
				compressed.GET("/withdrawals", h.getWithdrawals)
			}
		}
	}
}
