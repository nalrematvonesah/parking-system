package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"session-service/internal/nats"
	"session-service/internal/repository"
	"session-service/internal/service"

	"go.uber.org/zap"
)

// ---- mocks ----

type mockRepo struct {
	nextID   int64
	sessions map[int64]*repository.Session
	payments map[int64]float64
	failNext bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		sessions: make(map[int64]*repository.Session),
		payments: make(map[int64]float64),
		nextID:   1,
	}
}

func (m *mockRepo) CreateSession(_ context.Context, userID, vehicleID, slotID int64, start time.Time) (int64, error) {
	if m.failNext {
		m.failNext = false
		return 0, errors.New("db down")
	}
	id := m.nextID
	m.nextID++
	m.sessions[id] = &repository.Session{
		ID: id, UserID: userID, VehicleID: vehicleID, SlotID: slotID, StartTime: start,
	}
	return id, nil
}

func (m *mockRepo) EndSession(_ context.Context, sessionID int64, endTime time.Time) (*repository.Session, error) {
	s, ok := m.sessions[sessionID]
	if !ok {
		return nil, errors.New("not found")
	}
	if s.EndTime != nil {
		return nil, errors.New("already closed")
	}
	s.EndTime = &endTime
	return s, nil
}

func (m *mockRepo) GetSession(_ context.Context, sessionID int64) (*repository.Session, error) {
	return m.sessions[sessionID], nil
}

func (m *mockRepo) CreatePayment(_ context.Context, sessionID int64, amount float64) error {
	m.payments[sessionID] = amount
	return nil
}

func (m *mockRepo) ListActiveByUser(_ context.Context, userID int64) ([]repository.Session, error) {
	var out []repository.Session
	for _, s := range m.sessions {
		if s.UserID == userID && s.EndTime == nil {
			out = append(out, *s)
		}
	}
	return out, nil
}

type mockParking struct {
	released []int64
	nextSlot int64
}

func (p *mockParking) AssignSlot(_ context.Context, _ string) (int64, error) {
	p.nextSlot++
	return p.nextSlot, nil
}
func (p *mockParking) ReleaseSlot(_ context.Context, slotID int64) error {
	p.released = append(p.released, slotID)
	return nil
}

type mockCache struct{}

func (c *mockCache) SetActiveSession(_ context.Context, _ int64) error { return nil }
func (c *mockCache) DeleteActive(_ context.Context, _ int64) error     { return nil }

type mockPublisher struct {
	started   []nats.ParkingStartedEvent
	completed []nats.PaymentCompletedEvent
}

func (p *mockPublisher) PublishParkingStarted(e nats.ParkingStartedEvent) error {
	p.started = append(p.started, e)
	return nil
}
func (p *mockPublisher) PublishPaymentCompleted(e nats.PaymentCompletedEvent) error {
	p.completed = append(p.completed, e)
	return nil
}

func newSvc() (*service.Service, *mockRepo, *mockParking, *mockPublisher) {
	r := newMockRepo()
	p := &mockParking{}
	c := &mockCache{}
	pub := &mockPublisher{}
	svc := service.New(r, p, c, pub, zap.NewNop(), 500.0)
	return svc, r, p, pub
}

// ---- tests ----

func TestStartSession_Success(t *testing.T) {
	svc, r, _, pub := newSvc()
	s, err := svc.StartSession(context.Background(), 1, 2, "req-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID == 0 || s.SlotID == 0 {
		t.Fatalf("expected non-zero ids, got %+v", s)
	}
	if len(r.sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(r.sessions))
	}
	if len(pub.started) != 1 {
		t.Fatalf("expected parking.started event, got %d", len(pub.started))
	}
}

func TestStartSession_Validates(t *testing.T) {
	svc, _, _, _ := newSvc()
	if _, err := svc.StartSession(context.Background(), 0, 0, ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestStartSession_DBFails_ReleasesSlot(t *testing.T) {
	svc, r, p, _ := newSvc()
	r.failNext = true
	_, err := svc.StartSession(context.Background(), 1, 2, "")
	if err == nil {
		t.Fatal("expected error")
	}
	if len(p.released) != 1 {
		t.Fatalf("expected slot to be released on db failure, got %v", p.released)
	}
}

func TestEndSession_ChargesAndReleases(t *testing.T) {
	svc, r, p, pub := newSvc()
	s, _ := svc.StartSession(context.Background(), 1, 2, "")
	// rewind start to 90 minutes ago
	r.sessions[s.ID].StartTime = time.Now().Add(-90 * time.Minute)

	_, amount, err := svc.EndSession(context.Background(), s.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if amount != 1000 { // 1.5h → ceil to 2h → 2 * 500
		t.Fatalf("expected amount 1000, got %v", amount)
	}
	if len(p.released) != 1 {
		t.Fatalf("expected slot to be released, got %v", p.released)
	}
	if len(pub.completed) != 1 {
		t.Fatalf("expected payment.completed event, got %d", len(pub.completed))
	}
	if r.payments[s.ID] != 1000 {
		t.Fatalf("expected payment row 1000, got %v", r.payments[s.ID])
	}
}

func TestCalculatePrice_MinimumOneHour(t *testing.T) {
	svc, _, _, _ := newSvc()
	start := time.Now()
	end := start.Add(15 * time.Minute)
	if got := svc.CalculatePrice(start, end); got != 500 {
		t.Fatalf("expected 500 for <1h, got %v", got)
	}
}

func TestCalculatePrice_RoundsUp(t *testing.T) {
	svc, _, _, _ := newSvc()
	start := time.Now()
	end := start.Add(2*time.Hour + 1*time.Minute)
	if got := svc.CalculatePrice(start, end); got != 1500 { // ceil(2.016) = 3
		t.Fatalf("expected 1500 (3h * 500), got %v", got)
	}
}
