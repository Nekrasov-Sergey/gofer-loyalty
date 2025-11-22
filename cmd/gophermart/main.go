package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/app"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/config"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/delivery/rest"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/repository/postgres"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/service"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Приложение завершилось с ошибкой")
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	l := logger.New()

	r := router.New(gin.ReleaseMode, l)

	cfg, err := config.New(l)
	if err != nil {
		return err
	}

	repo, err := postgres.New(cfg.DatabaseURI, l)
	if err != nil {
		return err
	}
	defer multierr.AppendInvoke(&err, multierr.Close(repo))

	s := service.New(repo, cfg, l)

	h := rest.New(s, cfg, l)
	h.RegisterRoutes(r, repo)

	a := app.New(r, cfg.RunAddress, l)
	return a.Run(ctx)
}
