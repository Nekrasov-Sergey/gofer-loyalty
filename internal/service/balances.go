package service

import (
	"context"
	"time"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
)

func (s *Service) GetBalance(ctx context.Context, userID int64) (*types.BalanceResponse, error) {
	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &types.BalanceResponse{
		Current:   balance.Current.InexactFloat64(),
		Withdrawn: balance.Withdrawn.InexactFloat64(),
	}, nil
}

func (s *Service) WithdrawBalance(ctx context.Context, withdrawal *types.Withdrawal) error {
	return s.repo.WithTx(ctx, func(txRepo Repository) error {
		if err := txRepo.WithdrawBalance(ctx, withdrawal); err != nil {
			return err
		}
		return txRepo.AddWithdrawal(ctx, withdrawal)
	})
}

func (s *Service) GetWithdrawals(ctx context.Context, userID int64) ([]types.WithdrawalResponse, error) {
	withdrawals, err := s.repo.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawalsResponse := make([]types.WithdrawalResponse, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		withdrawalsResponse = append(withdrawalsResponse, types.WithdrawalResponse{
			Order:       withdrawal.OrderNumber,
			Sum:         withdrawal.Withdrawn.InexactFloat64(),
			ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
		})
	}

	return withdrawalsResponse, nil
}
