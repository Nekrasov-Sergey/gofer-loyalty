package app

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type App struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

func New(handler http.Handler, addr string, logger zerolog.Logger) *App {
	return &App{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}

func (a *App) Run() error {
	a.logger.Info().Msgf("Приложение запущено на %s", a.httpServer.Addr)

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "ошибка при запуске приложения")
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "ошибка при остановке приложения")
	}
	return nil
}
