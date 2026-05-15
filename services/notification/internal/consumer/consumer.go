package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"notification-service/internal/email"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Event subjects published by other services.
const (
	SubjectUserRegistered   = "user.registered"
	SubjectParkingStarted   = "parking.started"
	SubjectPaymentCompleted = "payment.completed"
)

// ---- payload types (must match producers) ----

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

type Consumer struct {
	conn   *nats.Conn
	mailer *email.Sender
	log    *zap.Logger
	subs   []*nats.Subscription
}

func New(natsURL string, mailer *email.Sender, log *zap.Logger) (*Consumer, error) {
	conn, err := nats.Connect(natsURL, nats.Name("notification-service"))
	if err != nil {
		return nil, fmt.Errorf("nats: connect %s: %w", natsURL, err)
	}
	return &Consumer{conn: conn, mailer: mailer, log: log}, nil
}

// Start subscribes to all event subjects. Each handler is wrapped so that
// invalid payloads don't crash the consumer.
func (c *Consumer) Start() error {
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{SubjectUserRegistered, c.handleUserRegistered},
		{SubjectParkingStarted, c.handleParkingStarted},
		{SubjectPaymentCompleted, c.handlePaymentCompleted},
	}

	for _, s := range subs {
		sub, err := c.conn.Subscribe(s.subject, s.handler)
		if err != nil {
			return fmt.Errorf("subscribe %s: %w", s.subject, err)
		}
		c.subs = append(c.subs, sub)
		c.log.Info("subscribed", zap.String("subject", s.subject))
	}
	return nil
}

func (c *Consumer) Close() {
	for _, s := range c.subs {
		_ = s.Unsubscribe()
	}
	_ = c.conn.Drain()
}

// ---- handlers ----

func (c *Consumer) handleUserRegistered(m *nats.Msg) {
	var e UserRegisteredEvent
	if err := json.Unmarshal(m.Data, &e); err != nil {
		c.log.Warn("bad user.registered payload", zap.Error(err))
		return
	}
	c.log.Info("event received", zap.String("subject", m.Subject),
		zap.Int64("user_id", e.UserID), zap.String("email", e.Email))

	_ = c.mailer.Send(
		e.Email,
		"Welcome to Smart Parking!",
		fmt.Sprintf("Hi,\n\nYour Smart Parking account (id=%d) has been created.\n\nEnjoy!", e.UserID),
	)
}

func (c *Consumer) handleParkingStarted(m *nats.Msg) {
	var e ParkingStartedEvent
	if err := json.Unmarshal(m.Data, &e); err != nil {
		c.log.Warn("bad parking.started payload", zap.Error(err))
		return
	}
	c.log.Info("event received", zap.String("subject", m.Subject),
		zap.Int64("session_id", e.SessionID), zap.Int64("user_id", e.UserID),
		zap.Int64("slot_id", e.SlotID))

	to := derivedEmail(e.UserID)
	_ = c.mailer.Send(
		to,
		fmt.Sprintf("Parking session #%d started", e.SessionID),
		fmt.Sprintf(
			"Your vehicle %d was parked in slot %d at %s.\n\nThanks for using Smart Parking!",
			e.VehicleID, e.SlotID, time.Unix(e.StartUnix, 0).UTC().Format(time.RFC1123),
		),
	)
}

func (c *Consumer) handlePaymentCompleted(m *nats.Msg) {
	var e PaymentCompletedEvent
	if err := json.Unmarshal(m.Data, &e); err != nil {
		c.log.Warn("bad payment.completed payload", zap.Error(err))
		return
	}
	c.log.Info("event received", zap.String("subject", m.Subject),
		zap.Int64("session_id", e.SessionID), zap.Float64("amount", e.Amount))

	to := derivedEmail(e.UserID)
	_ = c.mailer.Send(
		to,
		fmt.Sprintf("Receipt for parking session #%d", e.SessionID),
		fmt.Sprintf(
			"Thank you for using Smart Parking!\n\nSession: %d\nClosed at: %s\nAmount charged: %.2f\n",
			e.SessionID,
			time.Unix(e.EndUnix, 0).UTC().Format(time.RFC1123),
			e.Amount,
		),
	)
}

// derivedEmail builds a placeholder address when we don't carry the email
// inside the event payload (parking.started, payment.completed).
//
// In production this would be replaced by a lookup against the user service.
func derivedEmail(userID int64) string {
	return fmt.Sprintf("user-%d@parking.local", userID)
}
