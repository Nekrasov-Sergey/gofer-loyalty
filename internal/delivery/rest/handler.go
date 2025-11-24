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
	Register(ctx context.Context, user types.User) (sessionToken string, err error)
	Login(ctx context.Context, user types.User) (sessionToken string, err error)

	CreateOrder(ctx context.Context, order types.Order) error
	GetOrders(ctx context.Context, userID int64) ([]types.ResponseOrder, error)
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
	userAPI := r.Group("/api/user")
	userAPI.POST("/register", h.register)
	userAPI.POST("/login", h.login)

	userAuth := r.Group("/api/user", router.SessionMiddleware(repo))
	userAuth.POST("/orders", h.createOrder)
	userAuth.GET("/orders", h.getOrders)
	userAuth.GET("/balance")
	userAuth.POST("/balance/withdraw")
	userAuth.GET("/withdrawals")
}
