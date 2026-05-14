package service

import (
	"context"
	"errors"
	"fmt"

	"user-service/internal/email"
	"user-service/internal/observability"
	"user-service/internal/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repo interface {
	CreateUser(ctx context.Context, email, hashedPassword string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	GetUserByID(ctx context.Context, id int64) (*repository.User, error)
	AddVehicle(ctx context.Context, userID int64, plate string) error
	DeleteVehicle(ctx context.Context, userID int64, plate string) error
	GetVehiclesByUser(ctx context.Context, userID int64) ([]string, error)
}

type EventPublisher interface {
	PublishUserRegistered(userID int64, email string) error
}

type Service struct {
	repo      Repo
	publisher EventPublisher
	mailer    *email.Sender
	log       *zap.Logger
}

func New(repo Repo, publisher EventPublisher, mailer *email.Sender, log *zap.Logger) *Service {
	return &Service{repo: repo, publisher: publisher, mailer: mailer, log: log}
}

func (s *Service) Register(ctx context.Context, emailAddr, password string) (int64, error) {
	if emailAddr == "" || password == "" {
		return 0, errors.New("email and password are required")
	}

	existing, err := s.repo.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	id, err := s.repo.CreateUser(ctx, emailAddr, string(hash))
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	observability.IncRegistrations()

	if err := s.publisher.PublishUserRegistered(id, emailAddr); err != nil {
		s.log.Warn("failed to publish user.registered event", zap.Error(err))
	}

	go func() {
		if err := s.mailer.SendWelcome(emailAddr); err != nil {
			s.log.Warn("failed to send welcome email", zap.String("to", emailAddr), zap.Error(err))
		} else if s.mailer.IsConfigured() {
			s.log.Info("welcome email sent", zap.String("to", emailAddr))
		}
	}()

	return id, nil
}

func (s *Service) Login(ctx context.Context, emailAddr, password string) (int64, error) {
	if emailAddr == "" || password == "" {
		return 0, errors.New("email and password are required")
	}

	user, err := s.repo.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("invalid credentials")
	}

	observability.IncActiveSessions()
	return user.ID, nil
}

func (s *Service) Logout(ctx context.Context, userID int64) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	observability.DecActiveSessions()
	return nil
}

func (s *Service) AddVehicle(ctx context.Context, userID int64, plate string) error {
	if plate == "" {
		return errors.New("plate number is required")
	}
	return s.repo.AddVehicle(ctx, userID, plate)
}

func (s *Service) DeleteVehicle(ctx context.Context, userID int64, plate string) error {
	if plate == "" {
		return errors.New("plate number is required")
	}
	return s.repo.DeleteVehicle(ctx, userID, plate)
}

func (s *Service) GetVehicles(ctx context.Context, userID int64) ([]string, error) {
	return s.repo.GetVehiclesByUser(ctx, userID)
}
