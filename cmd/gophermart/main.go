package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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
	client := resty.New()
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

	s := service.New(repo, client, cfg, l)
	go s.RunAccrualWorker(ctx)

	h := rest.New(s, cfg, l)
	h.RegisterRoutes(r, repo)

	a := app.New(r, cfg.RunAddress, l)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- a.Run()
	}()

	select {
	case <-ctx.Done():
	case err := <-serverErr:
		if err != nil {
			return err
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		return err
	}

	l.Info().Msg("Приложение остановлено")
	return nil
}
