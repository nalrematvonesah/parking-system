package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Session struct {
	ID        int64
	UserID    int64
	VehicleID int64
	SlotID    int64
	StartTime time.Time
	EndTime   *time.Time
}

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r *PostgresRepo) CreateSession(
	ctx context.Context,
	userID, vehicleID, slotID int64,
	start time.Time,
) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO sessions (user_id, vehicle_id, slot_id, start_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, vehicleID, slotID, start).Scan(&id)
	return id, err
}

func (r *PostgresRepo) EndSession(
	ctx context.Context,
	sessionID int64,
	endTime time.Time,
) (*Session, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE sessions
		SET end_time = $1
		WHERE id = $2 AND end_time IS NULL
		RETURNING id, user_id, vehicle_id, slot_id, start_time, end_time
	`, endTime, sessionID)

	s := &Session{}
	err := row.Scan(&s.ID, &s.UserID, &s.VehicleID, &s.SlotID, &s.StartTime, &s.EndTime)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("session not found or already closed")
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepo) GetSession(ctx context.Context, sessionID int64) (*Session, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, vehicle_id, slot_id, start_time, end_time
		FROM sessions WHERE id = $1
	`, sessionID)

	s := &Session{}
	err := row.Scan(&s.ID, &s.UserID, &s.VehicleID, &s.SlotID, &s.StartTime, &s.EndTime)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepo) CreatePayment(
	ctx context.Context,
	sessionID int64,
	amount float64,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO payments (session_id, amount)
		VALUES ($1, $2)
	`, sessionID, amount)
	return err
}

func (r *PostgresRepo) ListActiveByUser(ctx context.Context, userID int64) ([]Session, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, vehicle_id, slot_id, start_time, end_time
		FROM sessions
		WHERE user_id = $1 AND end_time IS NULL
		ORDER BY start_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.VehicleID, &s.SlotID, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PostgresRepo) ListByUser(ctx context.Context, userID int64) ([]Session, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, vehicle_id, slot_id, start_time, end_time
		FROM sessions
		WHERE user_id = $1
		ORDER BY start_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.VehicleID, &s.SlotID, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
