package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

func (s *Service) CreateOrder(ctx context.Context, order types.Order) error {
	return s.repo.CreateOrder(ctx, order)
}

func (s *Service) GetOrders(ctx context.Context, userID int64) ([]types.ResponseOrder, error) {
	orders, err := s.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	externalOrders := make(map[string]types.ExternalOrder)

	mu := sync.Mutex{}

	numWorkers := s.config.NumWorkers
	if s.config.NumWorkers <= 0 {
		numWorkers = 1
	}

	sem := semaphore.NewWeighted(numWorkers)

	g, ctx := errgroup.WithContext(ctx)

	for _, order := range orders {
		g.Go(func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return errors.Wrap(err, "не удалось получить слот семафора")
			}
			defer sem.Release(1)

			externalOrder, err := s.getOrderFromExternalSystem(ctx, order.Number)
			if err != nil {
				return err
			}

			mu.Lock()
			externalOrders[externalOrder.Order] = externalOrder
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "ошибка при ожидании завершения воркеров")
	}

	responseOrders := make([]types.ResponseOrder, 0, len(orders))

	for _, order := range orders {
		externalOrder := externalOrders[order.Number]
		if externalOrder.Status == "" {
			continue
		}

		responseOrders = append(responseOrders, types.ResponseOrder{
			Number:     order.Number,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
			Status:     externalOrder.Status,
			Accrual:    externalOrder.Accrual,
		})
	}

	return responseOrders, nil
}

func (s *Service) getOrderFromExternalSystem(ctx context.Context, orderNumber string) (order types.ExternalOrder, err error) {
	path, err := url.JoinPath(s.config.AccrualSystemAddress, fmt.Sprintf("/api/orders/%s", orderNumber))
	if err != nil {
		return types.ExternalOrder{}, errors.Wrap(err, "Не удалось сформировать url")
	}

	numRetries := s.config.NumRetries
	if numRetries <= 0 {
		numRetries = 1
	}

	for i := 0; i < numRetries; i++ {
		select {
		case <-ctx.Done():
			return types.ExternalOrder{}, errors.WithStack(ctx.Err())
		default:
		}

		resp, err := s.client.R().
			SetContext(ctx).
			SetResult(&order).
			Get(path)
		if err != nil {
			return types.ExternalOrder{}, errors.Wrap(err, "ошибка при выполнении запроса во внешнюю систему")
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			retryAfter := 1
			if header := resp.Header().Get("Retry-After"); header != "" {
				if retryAfter, err = strconv.Atoi(header); err != nil {
					return types.ExternalOrder{}, errors.Wrap(err, fmt.Sprintf("не удалось распарсить заголовок Retry-After: %s", header))
				}
			}

			timer := time.NewTimer(time.Duration(retryAfter) * time.Second)
			select {
			case <-ctx.Done():
				return types.ExternalOrder{}, errors.WithStack(ctx.Err())
			case <-timer.C:
				continue
			}
		}

		if resp.IsError() {
			return types.ExternalOrder{}, errors.Errorf("сервер вернул ошибку: %d, body: %s", resp.StatusCode(), resp.String())
		}

		return order, nil
	}

	return types.ExternalOrder{}, errors.Errorf("превышено максимальное количество запросов во внешнюю систему: %d", numRetries)
}
