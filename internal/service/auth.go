package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
)

func (s *Service) Register(ctx context.Context, user types.User) (sessionToken string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return s.createSession(ctx, userID)
}

func (s *Service) Login(ctx context.Context, user types.User) (sessionToken string, err error) {
	dbUser, err := s.repo.GetUser(ctx, user.Login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return "", errcodes.ErrInvalidCredentials
	}

	return s.createSession(ctx, dbUser.ID)
}

func (s *Service) createSession(ctx context.Context, userID int64) (sessionToken string, err error) {
	sessionToken = uuid.NewString()
	if err := s.repo.CreateSession(ctx, sessionToken, userID, time.Now().Add(time.Duration(s.config.SessionTTL))); err != nil {
		return "", err
	}
	return sessionToken, nil
}
