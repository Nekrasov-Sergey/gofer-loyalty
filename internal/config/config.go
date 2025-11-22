package config

import (
	"flag"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/utils"
)

type Config struct {
	RunAddress           string       `env:"RUN_ADDRESS"`
	DatabaseURI          string       `env:"DATABASE_URI"`
	AccrualSystemAddress string       `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SessionTTL           HourDuration `env:"SESSION_TTL"`
}

func New(logger zerolog.Logger) (*Config, error) {
	runAddress := NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	flag.Var(&runAddress, "a", "адрес и порт запуска сервиса")

	databaseURI := flag.String("d", "postgres://user:pass@localhost:5432/loyalty?sslmode=disable", "адрес подключения к базе данных")

	accrualSystemAddress := NetAddress{
		Host: "localhost",
		Port: 8081,
	}
	flag.Var(&accrualSystemAddress, "r", "адрес системы расчёта начислений")

	sessionTTL := HourDuration(24 * time.Hour)
	flag.Var(&sessionTTL, "t", "время жизни сессии в часах")

	flag.Parse()

	cfg := Config{
		RunAddress:           runAddress.String(),
		DatabaseURI:          utils.Deref(databaseURI),
		AccrualSystemAddress: accrualSystemAddress.String(),
		SessionTTL:           sessionTTL,
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, errors.Wrap(err, "не удалось распарсить переменные окружения в конфиг")
	}

	logger.Info().
		Str("run_address", cfg.RunAddress).
		Str("database_uri", cfg.DatabaseURI).
		Str("accrual_system_address", cfg.AccrualSystemAddress).
		Str("session_ttl", cfg.SessionTTL.String()).
		Msg("Загружена конфигурация приложения")

	return &cfg, nil
}

type NetAddress struct {
	Host string
	Port int
}

func (a *NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return errors.New("адрес должен быть в формате host:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.Wrap(err, "неверный порт")
	}
	host := parts[0]
	if host == "" {
		host = "localhost"
	}
	a.Host = host
	a.Port = port
	return nil
}

type HourDuration time.Duration

func (d *HourDuration) String() string {
	return time.Duration(*d).String()
}

func (d *HourDuration) Set(s string) error {
	hours, err := strconv.Atoi(s)
	if err != nil {
		return errors.Wrap(err, "значение должно быть в часах")
	}
	*d = HourDuration(time.Duration(hours) * time.Hour)
	return nil
}

func (d *HourDuration) UnmarshalText(text []byte) error {
	hours, err := strconv.Atoi(string(text))
	if err != nil {
		return errors.Wrap(err, "значение должно быть в часах")
	}
	*d = HourDuration(time.Duration(hours) * time.Hour)
	return nil
}
