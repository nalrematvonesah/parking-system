package app

import (
	"notification-service/internal/config"
	"notification-service/internal/consumer"
	"notification-service/internal/email"

	"go.uber.org/zap"
)

type App struct {
	consumer *consumer.Consumer
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	mailer := email.New(
		cfg.SMTP.Host, cfg.SMTP.Port,
		cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From,
		log,
	)

	c, err := consumer.New(cfg.NatsURL, mailer, log)
	if err != nil {
		return nil, err
	}
	if err := c.Start(); err != nil {
		c.Close()
		return nil, err
	}
	return &App{consumer: c}, nil
}

func (a *App) Stop() {
	a.consumer.Close()
}
