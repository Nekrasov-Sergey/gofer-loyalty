package app

import (
	"context"
	"net/http"
	"time"

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

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		a.logger.Info().Msgf("Приложение запущено на %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error().Err(err).Msg("Ошибка при запуске приложения")
			cancel()
		}
	}()

	<-ctx.Done()
	a.logger.Info().Msg("Остановка приложения...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := a.shutdown(shutdownCtx); err != nil {
		return err
	}

	a.logger.Info().Msg("Приложение остановлено")
	return nil
}

func (a *App) shutdown(shutdownCtx context.Context) error {
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return errors.Wrap(err, "ошибка при остановке приложения")
	}
	return nil
}
