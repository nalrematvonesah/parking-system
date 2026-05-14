package service

import (
	"context"
	"session-service/internal/cache"
	"session-service/internal/nats"
	"session-service/internal/repository"
	"time"
)

type Service struct {
	repo       *repository.PostgresRepo
	parking    *ParkingClient
	cache      *cache.Redis
	natsClient *nats.Client
}

func New(
	repo *repository.PostgresRepo,
	parking *ParkingClient,
	cache *cache.Redis,
	natsClient *nats.Client,
) *Service {
	return &Service{
		repo:       repo,
		parking:    parking,
		cache:      cache,
		natsClient: natsClient,
	}
}

func (s *Service) StartSession(
	ctx context.Context,
	vehicleID string,
) error {
	slotID, err := s.parking.AssignSlot(ctx)
	if err != nil {
		return err
	}
	err = s.cache.SetSession(
		ctx,
		"session:"+vehicleID,
		"active",
	)
	err = s.natsClient.Publish(
		"ParkingStarted",
		[]byte(vehicleID),
	)

	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return s.repo.CreateSession(
		ctx,
		vehicleID,
		slotID,
		time.Now(),
	)
}

func (s *Service) CalculatePrice(start time.Time) float64 {
	duration := time.Since(start).Hours()

	if duration < 1 {
		duration = 1
	}

	return duration * 500
}

func (s *Service) ProcessPayment(
	ctx context.Context,
	sessionID int64,
	amount float64,
) error {
	return s.repo.CreatePayment(
		ctx,
		sessionID,
		amount,
	)
}

func (s *Service) EndSession(
	ctx context.Context,
	sessionID int64,
	startTime time.Time,
) error {
	endTime := time.Now()

	err := s.repo.EndSession(
		ctx,
		sessionID,
		endTime,
	)

	if err != nil {
		return err
	}

	price := s.CalculatePrice(startTime)

	err = s.ProcessPayment(
		ctx,
		sessionID,
		price,
	)

	if err != nil {
		return err
	}

	err = s.natsClient.Publish(
		"PaymentCompleted",
		[]byte("payment completed"),
	)

	if err != nil {
		return err
	}

	return nil
}
