package errcodes

import (
	"github.com/pkg/errors"
)

var (
	ErrLoginAlreadyExists                = errors.New("логин уже занят")
	ErrInvalidCredentials                = errors.New("неверная пара логин/пароль")
	ErrUserUnauthorized                  = errors.New("пользователь не авторизован")
	ErrOrderAlreadyUploadedByUser        = errors.New("заказ уже был загружен этим пользователем")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("заказ уже был загружен другим пользователем")
)
