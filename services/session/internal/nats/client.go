package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// Subjects used by the session service.
const (
	SubjectParkingStarted   = "parking.started"
	SubjectPaymentCompleted = "payment.completed"
)

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

type Publisher struct {
	conn *nats.Conn
}

func New(url string) (*Publisher, error) {
	conn, err := nats.Connect(url, nats.Name("session-service"))
	if err != nil {
		return nil, fmt.Errorf("nats: connect %s: %w", url, err)
	}
	return &Publisher{conn: conn}, nil
}

func (p *Publisher) PublishParkingStarted(e ParkingStartedEvent) error {
	return p.publish(SubjectParkingStarted, e)
}

func (p *Publisher) PublishPaymentCompleted(e PaymentCompletedEvent) error {
	return p.publish(SubjectPaymentCompleted, e)
}

func (p *Publisher) publish(subject string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, data)
}

func (p *Publisher) Close() {
	_ = p.conn.Drain()
}
