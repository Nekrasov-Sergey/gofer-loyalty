package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (s *Service) CreateOrder(ctx context.Context, orderRequest *types.Order) error {
	order, err := s.repo.GetOrderByNumber(ctx, orderRequest.Number)
	if err != nil && !errors.Is(err, errcodes.ErrOrderNotFound) {
		return err
	}

	if err == nil {
		if order.UserID == orderRequest.UserID {
			return errcodes.ErrOrderAlreadyUploadedByUser
		}
		return errcodes.ErrOrderAlreadyUploadedByAnotherUser
	}

	orderRequest.Accrual = decimal.Zero
	orderRequest.Status = types.OrderStatusNew

	return s.repo.CreateOrder(ctx, orderRequest)
}

func (s *Service) GetOrders(ctx context.Context, userID int64) ([]types.OrderResponse, error) {
	orders, err := s.repo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	ordersResponse := make([]types.OrderResponse, 0, len(orders))
	for _, order := range orders {
		ordersResponse = append(ordersResponse, types.OrderResponse{
			Number:     order.Number,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
			Status:     order.Status,
			Accrual:    order.Accrual.InexactFloat64(),
		})
	}

	return ordersResponse, nil
}
