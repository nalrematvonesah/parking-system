package service_test

import (
	"context"
	"testing"

	"user-service/internal/email"
	"user-service/internal/repository"
	"user-service/internal/service"

	"go.uber.org/zap"
)

type mockRepo struct {
	users    map[string]*repository.User
	vehicles map[int64][]string
	nextID   int64
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users:    make(map[string]*repository.User),
		vehicles: make(map[int64][]string),
		nextID:   1,
	}
}

func (m *mockRepo) CreateUser(_ context.Context, eml, hash string) (int64, error) {
	id := m.nextID
	m.nextID++
	m.users[eml] = &repository.User{ID: id, Email: eml, Password: hash}
	return id, nil
}

func (m *mockRepo) GetUserByEmail(_ context.Context, eml string) (*repository.User, error) {
	if u, ok := m.users[eml]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockRepo) GetUserByID(_ context.Context, id int64) (*repository.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) AddVehicle(_ context.Context, userID int64, plate string) error {
	m.vehicles[userID] = append(m.vehicles[userID], plate)
	return nil
}

func (m *mockRepo) DeleteVehicle(_ context.Context, userID int64, plate string) error {
	plates := m.vehicles[userID]
	for i, p := range plates {
		if p == plate {
			m.vehicles[userID] = append(plates[:i], plates[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockRepo) GetVehiclesByUser(_ context.Context, userID int64) ([]string, error) {
	return m.vehicles[userID], nil
}

type mockPublisher struct{}

func (p *mockPublisher) PublishUserRegistered(_ int64, _ string) error { return nil }

func newService() *service.Service {
	return service.New(
		newMockRepo(),
		&mockPublisher{},
		email.New("", "", "", "", ""),
		zap.NewNop(),
	)
}

func TestRegister_Success(t *testing.T) {
	svc := newService()
	id, err := svc.Register(context.Background(), "askhat@test.com", "secret123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero user ID")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newService()
	_, _ = svc.Register(context.Background(), "askhat@test.com", "secret123")
	_, err := svc.Register(context.Background(), "askhat@test.com", "other")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestRegister_EmptyFields(t *testing.T) {
	svc := newService()
	_, err := svc.Register(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLogin_Success(t *testing.T) {
	svc := newService()
	_, _ = svc.Register(context.Background(), "askhat@test.com", "secret123")
	id, err := svc.Login(context.Background(), "askhat@test.com", "secret123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero user ID")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newService()
	_, _ = svc.Register(context.Background(), "askhat@test.com", "secret123")
	_, err := svc.Login(context.Background(), "askhat@test.com", "wrong")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	svc := newService()
	_, err := svc.Login(context.Background(), "nobody@test.com", "pass")
	if err == nil {
		t.Fatal("expected error for unknown email")
	}
}

func TestAddVehicle_Success(t *testing.T) {
	svc := newService()
	id, _ := svc.Register(context.Background(), "askhat@test.com", "secret123")
	if err := svc.AddVehicle(context.Background(), id, "ABC-123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGetVehicles(t *testing.T) {
	svc := newService()
	id, _ := svc.Register(context.Background(), "askhat@test.com", "secret123")
	_ = svc.AddVehicle(context.Background(), id, "ABC-123")
	_ = svc.AddVehicle(context.Background(), id, "XYZ-999")

	plates, err := svc.GetVehicles(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plates) != 2 {
		t.Fatalf("expected 2 vehicles, got %d", len(plates))
	}
}

func TestDeleteVehicle(t *testing.T) {
	svc := newService()
	id, _ := svc.Register(context.Background(), "askhat@test.com", "secret123")
	_ = svc.AddVehicle(context.Background(), id, "ABC-123")
	_ = svc.DeleteVehicle(context.Background(), id, "ABC-123")

	plates, _ := svc.GetVehicles(context.Background(), id)
	if len(plates) != 0 {
		t.Fatalf("expected empty list after delete, got %v", plates)
	}
}
