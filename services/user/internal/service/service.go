package service

import (
	"context"
	"errors"

	"user-service/internal/cache"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  *repository.PostgresRepo
	cache *cache.RedisCache
}

func New(repo *repository.PostgresRepo, cache *cache.RedisCache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) Register(ctx context.Context, email, password string) (int64, error) {
	if email == "" || password == "" {
		return 0, errors.New("email and password are required")
	}

	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateUser(ctx, email, string(hash))
}

func (s *Service) Login(ctx context.Context, email, password string) (int64, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("invalid credentials")
	}

	s.cache.SetSession(ctx, user.ID)
	return user.ID, nil
}

func (s *Service) AddVehicle(ctx context.Context, userID int64, plate string) error {
	if plate == "" {
		return errors.New("plate number is required")
	}
	return s.repo.AddVehicle(ctx, userID, plate)
}

func (s *Service) GetVehicles(ctx context.Context, userID int64) ([]string, error) {
	return s.repo.GetVehiclesByUser(ctx, userID)
}
