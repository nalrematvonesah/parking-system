package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"session-service/internal/nats"
	"session-service/internal/repository"
	"time"

	"go.uber.org/zap"
)

// Repo abstracts the storage so we can mock it in tests.
type Repo interface {
	CreateSession(ctx context.Context, userID, vehicleID, slotID int64, start time.Time) (int64, error)
	EndSession(ctx context.Context, sessionID int64, endTime time.Time) (*repository.Session, error)
	GetSession(ctx context.Context, sessionID int64) (*repository.Session, error)
	CreatePayment(ctx context.Context, sessionID int64, amount float64) error
	ListActiveByUser(ctx context.Context, userID int64) ([]repository.Session, error)
	ListByUser(ctx context.Context, userID int64) ([]repository.Session, error)
}

type ParkingPort interface {
	AssignSlot(ctx context.Context, requestID string) (int64, error)
	ReleaseSlot(ctx context.Context, slotID int64) error
}

type CachePort interface {
	SetActiveSession(ctx context.Context, sessionID int64) error
	DeleteActive(ctx context.Context, sessionID int64) error
}

type EventPublisher interface {
	PublishParkingStarted(e nats.ParkingStartedEvent) error
	PublishPaymentCompleted(e nats.PaymentCompletedEvent) error
}

type Service struct {
	repo         Repo
	parking      ParkingPort
	cache        CachePort
	publisher    EventPublisher
	log          *zap.Logger
	pricePerHour float64
}

func New(
	repo Repo,
	parking ParkingPort,
	cache CachePort,
	publisher EventPublisher,
	log *zap.Logger,
	pricePerHour float64,
) *Service {
	return &Service{
		repo:         repo,
		parking:      parking,
		cache:        cache,
		publisher:    publisher,
		log:          log,
		pricePerHour: pricePerHour,
	}
}

// StartSession reserves a parking slot and opens a billing session.
// If the session can't be persisted, the previously reserved slot is released.
func (s *Service) StartSession(
	ctx context.Context,
	userID, vehicleID int64,
	requestID string,
) (*repository.Session, error) {
	if userID == 0 || vehicleID == 0 {
		return nil, errors.New("user_id and vehicle_id are required")
	}

	slotID, err := s.parking.AssignSlot(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("assign slot: %w", err)
	}

	start := time.Now().UTC()
	sessionID, err := s.repo.CreateSession(ctx, userID, vehicleID, slotID, start)
	if err != nil {
		// best-effort rollback of the parking-side reservation
		if releaseErr := s.parking.ReleaseSlot(ctx, slotID); releaseErr != nil {
			s.log.Warn("failed to release slot after persist failure",
				zap.Int64("slot_id", slotID), zap.Error(releaseErr))
		}
		return nil, fmt.Errorf("create session: %w", err)
	}

	if err := s.cache.SetActiveSession(ctx, sessionID); err != nil {
		s.log.Warn("cache SetActiveSession failed", zap.Error(err))
	}

	if err := s.publisher.PublishParkingStarted(nats.ParkingStartedEvent{
		SessionID: sessionID,
		UserID:    userID,
		VehicleID: vehicleID,
		SlotID:    slotID,
		StartUnix: start.Unix(),
	}); err != nil {
		s.log.Warn("publish parking.started failed", zap.Error(err))
	}

	return &repository.Session{
		ID:        sessionID,
		UserID:    userID,
		VehicleID: vehicleID,
		SlotID:    slotID,
		StartTime: start,
	}, nil
}

// EndSession closes a session, releases the slot, computes price and stores payment.
func (s *Service) EndSession(ctx context.Context, sessionID int64) (*repository.Session, float64, error) {
	end := time.Now().UTC()

	session, err := s.repo.EndSession(ctx, sessionID, end)
	if err != nil {
		return nil, 0, err
	}

	if err := s.parking.ReleaseSlot(ctx, session.SlotID); err != nil {
		s.log.Warn("release slot failed", zap.Int64("slot_id", session.SlotID), zap.Error(err))
	}

	amount := s.CalculatePrice(session.StartTime, end)

	if err := s.repo.CreatePayment(ctx, sessionID, amount); err != nil {
		return nil, 0, fmt.Errorf("create payment: %w", err)
	}

	if err := s.cache.DeleteActive(ctx, sessionID); err != nil {
		s.log.Warn("cache DeleteActive failed", zap.Error(err))
	}

	if err := s.publisher.PublishPaymentCompleted(nats.PaymentCompletedEvent{
		SessionID: sessionID,
		UserID:    session.UserID,
		Amount:    amount,
		EndUnix:   end.Unix(),
	}); err != nil {
		s.log.Warn("publish payment.completed failed", zap.Error(err))
	}

	return session, amount, nil
}

// GetSession returns a session by id or nil if it doesn't exist.
func (s *Service) GetSession(ctx context.Context, sessionID int64) (*repository.Session, error) {
	return s.repo.GetSession(ctx, sessionID)
}

// PriceFor returns the accrued price of an active session as of now.
func (s *Service) PriceFor(ctx context.Context, sessionID int64) (float64, time.Duration, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}
	if session == nil {
		return 0, 0, errors.New("session not found")
	}
	end := time.Now().UTC()
	if session.EndTime != nil {
		end = *session.EndTime
	}
	elapsed := end.Sub(session.StartTime)
	return s.CalculatePrice(session.StartTime, end), elapsed, nil
}

// CalculatePrice rounds up to the next hour, with a minimum of one hour.
func (s *Service) CalculatePrice(start, end time.Time) float64 {
	hours := end.Sub(start).Hours()
	if hours < 1 {
		hours = 1
	} else {
		hours = math.Ceil(hours)
	}
	return hours * s.pricePerHour
}

func (s *Service) ListActiveByUser(ctx context.Context, userID int64) ([]repository.Session, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}
	return s.repo.ListActiveByUser(ctx, userID)
}

func (s *Service) ListByUser(ctx context.Context, userID int64) ([]repository.Session, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}
	return s.repo.ListByUser(ctx, userID)
}
