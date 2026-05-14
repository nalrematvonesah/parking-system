package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

const subjectUserRegistered = "user.registered"

type Publisher struct {
	nc *nats.Conn
}

func New(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats: connect to %s: %w", url, err)
	}
	return &Publisher{nc: nc}, nil
}

type UserRegisteredEvent struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (p *Publisher) PublishUserRegistered(userID int64, email string) error {
	payload, err := json.Marshal(UserRegisteredEvent{UserID: userID, Email: email})
	if err != nil {
		return err
	}
	return p.nc.Publish(subjectUserRegistered, payload)
}

func (p *Publisher) Close() {
	_ = p.nc.Drain()
}
