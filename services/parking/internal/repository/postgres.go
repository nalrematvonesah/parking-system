package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) AssignSlot(ctx context.Context) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var slotID int64
	err = tx.QueryRow(ctx, `
		SELECT id FROM parking_slots
		WHERE is_occupied = false
		ORDER BY id
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`).Scan(&slotID)

	if err != nil {
		return 0, errors.New("no available slots")
	}

	_, err = tx.Exec(ctx, `
		UPDATE parking_slots
		SET is_occupied = true, updated_at = NOW()
		WHERE id = $1
	`, slotID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return slotID, nil
}

func (r *PostgresRepo) ReleaseSlot(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE parking_slots
		SET is_occupied = false, updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *PostgresRepo) CountFree(ctx context.Context) (int32, error) {
	var c int32
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM parking_slots WHERE is_occupied = false
	`).Scan(&c)
	return c, err
}

type Slot struct {
	ID         int64
	IsOccupied bool
}

func (r *PostgresRepo) GetSlot(ctx context.Context, id int64) (*Slot, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, is_occupied
		FROM parking_slots
		WHERE id = $1
	`, id)

	var s Slot
	if err := row.Scan(&s.ID, &s.IsOccupied); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *PostgresRepo) ListAllSlots(ctx context.Context) ([]Slot, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, is_occupied
		FROM parking_slots
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Slot
	for rows.Next() {
		var s Slot
		if err := rows.Scan(&s.ID, &s.IsOccupied); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

func (r *PostgresRepo) AddSlots(ctx context.Context, count int32) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO parking_slots (is_occupied)
		SELECT false FROM generate_series(1, $1)
	`, count)

	return err
}
