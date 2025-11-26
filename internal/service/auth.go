package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (s *Service) Register(ctx context.Context, user *types.User) (sessionToken string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	sessionToken = uuid.NewString()

	err = s.repo.WithTx(ctx, func(txRepo Repository) error {
		userID, err := txRepo.CreateUser(ctx, user)
		if err != nil {
			return err
		}

		if err := txRepo.CreateBalance(ctx, userID); err != nil {
			return err
		}

		return txRepo.CreateSession(ctx, &types.Session{
			Token:     sessionToken,
			UserID:    userID,
			ExpiresAt: time.Now().Add(time.Duration(s.config.SessionTTL)),
		})
	})
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (s *Service) Login(ctx context.Context, user *types.User) (sessionToken string, err error) {
	dbUser, err := s.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return "", errcodes.ErrInvalidCredentials
	}

	sessionToken = uuid.NewString()
	return sessionToken, s.repo.CreateSession(ctx, &types.Session{
		Token:     sessionToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Duration(s.config.SessionTTL)),
	})
}
