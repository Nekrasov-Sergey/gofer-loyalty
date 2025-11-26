package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

func (s *Service) RunAccrualWorker(ctx context.Context) {
	s.logger.Info().Msg("Запущен воркер")

	workerCount := s.config.WorkerCount
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}

	limit := workerCount * 10
	ordersChan := make(chan *types.Order, limit)

	wg := &sync.WaitGroup{}

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range ordersChan {
				select {
				case <-ctx.Done():
					return
				case <-s.globalPauseCh:
				}
				s.processOrder(ctx, order)
			}
		}()
	}

	ticker := time.NewTicker(time.Duration(s.config.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Msg("Воркер остановлен")
			close(ordersChan)
			wg.Wait()
			return
		case <-ticker.C:
			orders, err := s.repo.GetOrdersForProcessing(ctx, limit)
			if err != nil {
				s.logger.Error().Err(err).Msg("Не удалось получить заказы из БД")
				continue
			}

			if len(orders) == 0 {
				continue
			}

			for i := range orders {
				select {
				case <-ctx.Done():
					close(ordersChan)
					wg.Wait()
					return
				case ordersChan <- &orders[i]:
				}
			}
		}
	}
}

func (s *Service) processOrder(ctx context.Context, order *types.Order) {
	externalOrder, err := s.getOrderFromExternalSystem(ctx, order.Number)
	if err != nil {
		s.logger.Error().Err(err).Msg("Не удалось получить заказ из внешней системы")
		return
	}

	if externalOrder == nil {
		return
	}

	switch externalOrder.Status {
	case types.ExternalOrderStatusEmpty, types.ExternalOrderStatusRegistered:
		return
	case types.ExternalOrderStatusProcessing:
		if order.Status == types.OrderStatusProcessing {
			return
		}
		order.Status = types.OrderStatusProcessing
	case types.ExternalOrderStatusInvalid:
		order.Status = types.OrderStatusInvalid
	case types.ExternalOrderStatusProcessed:
		order.Status = types.OrderStatusProcessed
		order.Accrual = externalOrder.Accrual
	}

	err = s.repo.WithTx(ctx, func(txRepo Repository) error {
		if err := txRepo.UpdateOrder(ctx, order); err != nil {
			return err
		}
		if order.Accrual.IsPositive() {
			return txRepo.AddBalance(ctx, order.UserID, externalOrder.Accrual)
		}
		return nil
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Не удалось обновить заказ")
		return
	}
}

func (s *Service) getOrderFromExternalSystem(ctx context.Context, orderNumber string) (order *types.ExternalOrder, err error) {
	path, err := url.JoinPath(s.config.AccrualSystemAddress, fmt.Sprintf("/api/orders/%s", orderNumber))
	if err != nil {
		return nil, errors.Wrap(err, "не удалось сформировать url")
	}

	numRetries := s.config.NumRetries
	if numRetries <= 0 {
		numRetries = 1
	}

	for i := 0; i < numRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, errors.WithStack(ctx.Err())
		default:
		}

		order = &types.ExternalOrder{}
		resp, err := s.client.R().
			SetContext(ctx).
			SetResult(order).
			Get(path)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка при выполнении запроса к системе расчёта баллов")
		}

		switch resp.StatusCode() {
		case http.StatusOK:
			return order, nil

		case http.StatusNoContent:
			return nil, nil

		case http.StatusTooManyRequests:
			retryAfter := 1
			if header := resp.Header().Get("Retry-After"); header != "" {
				if retryAfter, err = strconv.Atoi(header); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("не удалось распарсить заголовок Retry-After: %s", header))
				}
			}

			pauseDuration := time.Duration(retryAfter) * time.Second

			s.pauseAllWorkers(ctx, pauseDuration)

			select {
			case <-ctx.Done():
				return nil, errors.WithStack(ctx.Err())
			case <-time.After(pauseDuration):
				continue
			}

		default:
			return nil, errors.Errorf("система расчёта баллов вернула ошибку: %d, body: %s", resp.StatusCode(), resp.String())
		}
	}

	return nil, errors.Errorf("превышено максимальное количество запросов к системе расчета баллов: %d", numRetries)
}

func (s *Service) pauseAllWorkers(ctx context.Context, pauseDuration time.Duration) {
	s.globalPauseMu.Lock()
	defer s.globalPauseMu.Unlock()

	select {
	case <-s.globalPauseCh:
	default:
		return
	}

	s.globalPauseCh = make(chan struct{})

	s.logger.Info().Str("duration", pauseDuration.String()).Msg("Получен 429 — все воркеры поставлены на паузу")

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(pauseDuration):
		}

		s.globalPauseMu.Lock()
		close(s.globalPauseCh)
		s.globalPauseMu.Unlock()

		s.logger.Info().Msg("Пауза для воркеров снята")
	}()
}
