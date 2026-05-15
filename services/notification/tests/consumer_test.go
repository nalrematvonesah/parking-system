package notification_test

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ---- email sender mock ----

type capturedEmail struct {
	to      string
	subject string
	body    string
}

type mockSender struct {
	mu     sync.Mutex
	emails []capturedEmail
}

func (m *mockSender) Send(to, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emails = append(m.emails, capturedEmail{to: to, subject: subject, body: body})
	return nil
}

func (m *mockSender) last() capturedEmail {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.emails) == 0 {
		return capturedEmail{}
	}
	return m.emails[len(m.emails)-1]
}

func (m *mockSender) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.emails)
}

// ---- event payload types (mirror consumer package) ----

type UserRegisteredEvent struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

type ParkingStartedEvent struct {
	SessionID int64 `json:"session_id"`
	UserID    int64 `json:"user_id"`
	VehicleID int64 `json:"vehicle_id"`
	SlotID    int64 `json:"slot_id"`
	StartUnix int64 `json:"start_unix"`
}

type PaymentCompletedEvent struct {
	SessionID int64   `json:"session_id"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	EndUnix   int64   `json:"end_unix"`
}

// ---- local handler logic mirrors consumer package exactly ----

func derivedEmail(userID int64) string {
	return "user-" + itoa(userID) + "@parking.local"
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 20)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

type handler struct {
	sender *mockSender
	log    *zap.Logger
}

func (h *handler) handleUserRegistered(data []byte) {
	var e UserRegisteredEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return
	}
	_ = h.sender.Send(
		e.Email,
		"Welcome to Smart Parking!",
		"Hi,\n\nYour Smart Parking account (id="+itoa(e.UserID)+") has been created.\n\nEnjoy!",
	)
}

func (h *handler) handleParkingStarted(data []byte) {
	var e ParkingStartedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return
	}
	to := derivedEmail(e.UserID)
	_ = h.sender.Send(
		to,
		"Parking session #"+itoa(e.SessionID)+" started",
		"Your vehicle "+itoa(e.VehicleID)+" was parked in slot "+itoa(e.SlotID)+
			" at "+time.Unix(e.StartUnix, 0).UTC().Format(time.RFC1123)+
			".\n\nThanks for using Smart Parking!",
	)
}

func (h *handler) handlePaymentCompleted(data []byte) {
	var e PaymentCompletedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return
	}
	to := derivedEmail(e.UserID)
	_ = h.sender.Send(
		to,
		"Receipt for parking session #"+itoa(e.SessionID),
		"Thank you for using Smart Parking!\n\nSession: "+itoa(e.SessionID)+"\nAmount charged: 750.00\n",
	)
}

// ---- tests ----

func newHandler() (*handler, *mockSender) {
	s := &mockSender{}
	log, _ := zap.NewDevelopment()
	return &handler{sender: s, log: log}, s
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestHandleUserRegistered_SendsEmail(t *testing.T) {
	h, s := newHandler()
	payload := mustMarshal(UserRegisteredEvent{UserID: 42, Email: "test@example.com"})

	h.handleUserRegistered(payload)

	if s.count() != 1 {
		t.Fatalf("expected 1 email, got %d", s.count())
	}
	e := s.last()
	if e.to != "test@example.com" {
		t.Errorf("expected to=test@example.com, got %q", e.to)
	}
	if !strings.Contains(e.subject, "Welcome") {
		t.Errorf("expected welcome subject, got %q", e.subject)
	}
	if !strings.Contains(e.body, "42") {
		t.Errorf("expected user id in body, got %q", e.body)
	}
}

func TestHandleUserRegistered_BadJSON(t *testing.T) {
	h, s := newHandler()
	h.handleUserRegistered([]byte(`{bad json`))
	if s.count() != 0 {
		t.Fatal("expected no email on bad payload")
	}
}

func TestHandleParkingStarted_SendsEmail(t *testing.T) {
	h, s := newHandler()
	now := time.Now().Unix()
	payload := mustMarshal(ParkingStartedEvent{
		SessionID: 7, UserID: 3, VehicleID: 5, SlotID: 12, StartUnix: now,
	})

	h.handleParkingStarted(payload)

	if s.count() != 1 {
		t.Fatalf("expected 1 email, got %d", s.count())
	}
	e := s.last()
	if e.to != "user-3@parking.local" {
		t.Errorf("expected derived email, got %q", e.to)
	}
	if !strings.Contains(e.subject, "7") {
		t.Errorf("expected session id in subject, got %q", e.subject)
	}
	if !strings.Contains(e.body, "slot 12") {
		t.Errorf("expected slot id in body, got %q", e.body)
	}
}

func TestHandleParkingStarted_BadJSON(t *testing.T) {
	h, s := newHandler()
	h.handleParkingStarted([]byte(`not-json`))
	if s.count() != 0 {
		t.Fatal("expected no email on bad payload")
	}
}

func TestHandlePaymentCompleted_SendsEmail(t *testing.T) {
	h, s := newHandler()
	payload := mustMarshal(PaymentCompletedEvent{
		SessionID: 99, UserID: 8, Amount: 750.00, EndUnix: time.Now().Unix(),
	})

	h.handlePaymentCompleted(payload)

	if s.count() != 1 {
		t.Fatalf("expected 1 email, got %d", s.count())
	}
	e := s.last()
	if e.to != "user-8@parking.local" {
		t.Errorf("expected derived email, got %q", e.to)
	}
	if !strings.Contains(e.subject, "99") {
		t.Errorf("expected session id in subject, got %q", e.subject)
	}
}

func TestHandlePaymentCompleted_BadJSON(t *testing.T) {
	h, s := newHandler()
	h.handlePaymentCompleted([]byte(`{"bad`))
	if s.count() != 0 {
		t.Fatal("expected no email on bad payload")
	}
}

func TestDerivedEmail(t *testing.T) {
	cases := []struct {
		userID int64
		want   string
	}{
		{1, "user-1@parking.local"},
		{42, "user-42@parking.local"},
		{100, "user-100@parking.local"},
	}
	for _, tc := range cases {
		got := derivedEmail(tc.userID)
		if got != tc.want {
			t.Errorf("derivedEmail(%d) = %q, want %q", tc.userID, got, tc.want)
		}
	}
}

func TestHandleMultipleEvents_InOrder(t *testing.T) {
	h, s := newHandler()

	h.handleUserRegistered(mustMarshal(UserRegisteredEvent{UserID: 1, Email: "a@test.com"}))
	h.handleParkingStarted(mustMarshal(ParkingStartedEvent{SessionID: 1, UserID: 1, SlotID: 5, StartUnix: time.Now().Unix()}))
	h.handlePaymentCompleted(mustMarshal(PaymentCompletedEvent{SessionID: 1, UserID: 1, Amount: 500}))

	if s.count() != 3 {
		t.Fatalf("expected 3 emails for full flow, got %d", s.count())
	}
}
