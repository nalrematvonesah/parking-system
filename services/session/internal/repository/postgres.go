package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r *PostgresRepo) CreateSession(
	ctx context.Context,
	vehicleID string,
	slotID int64,
	start time.Time,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions(vehicle_id, slot_id, start_time)
		VALUES($1, $2, $3)
	`, vehicleID, slotID, start)

	return err
}

func (r *PostgresRepo) CreatePayment(
	ctx context.Context,
	sessionID int64,
	amount float64,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO payments(session_id, amount)
		VALUES($1, $2)
	`, sessionID, amount)

	return err
}

func (r *PostgresRepo) EndSession(
	ctx context.Context,
	sessionID int64,
	endTime time.Time,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE sessions
		SET end_time = $1
		WHERE id = $2
		`,
		endTime,
		sessionID,
	)

	return err
}
