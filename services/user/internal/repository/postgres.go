package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int64
	Email    string
	Password string
}

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateUser(ctx context.Context, email, hashedPassword string) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		email, hashedPassword,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit(ctx)
}

func (r *PostgresRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *PostgresRepo) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *PostgresRepo) AddVehicle(ctx context.Context, userID int64, plate string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO vehicles (user_id, plate_number) VALUES ($1, $2)`,
		userID, plate,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) DeleteVehicle(ctx context.Context, userID int64, plate string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx,
		`DELETE FROM vehicles WHERE user_id = $1 AND plate_number = $2`,
		userID, plate,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("vehicle not found")
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) GetVehiclesByUser(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT plate_number FROM vehicles WHERE user_id = $1 ORDER BY created_at`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plates []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		plates = append(plates, p)
	}
	return plates, rows.Err()
}
