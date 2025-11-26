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
	ErrNotEnoughBalance                  = errors.New("на счету недостаточно средств")
	ErrInvalidRequestFormat              = errors.New("неверный формат запроса")
	ErrInvalidOrderNumberFormat          = errors.New("неверный формат номера заказа")
	ErrOrderNumberMustBeNumeric          = errors.New("номер заказа должен быть числом")
	ErrLoginIsMissing                    = errors.New("отсутствует логин")
	ErrPasswordIsMissing                 = errors.New("отсутствует пароль")
	ErrPortIsRequired                    = errors.New("порт обязателен")
	ErrPortNotNumber                     = errors.New("порт должен быть числом")
	ErrAddressInvalidFormat              = errors.New("адрес должен быть в формате host:port")
	ErrOrderNotFound                     = errors.New("заказ не найден")
	ErrNegativeWithdrawal                = errors.New("нельзя списать отрицательно количество баллов")
)
