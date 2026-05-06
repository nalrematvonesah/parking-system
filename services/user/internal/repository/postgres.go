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
	var id int64
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		email, hashedPassword,
	).Scan(&id)
	return id, err
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

func (r *PostgresRepo) AddVehicle(ctx context.Context, userID int64, plate string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO vehicles (user_id, plate_number) VALUES ($1, $2)`,
		userID, plate,
	)
	return err
}

func (r *PostgresRepo) GetVehiclesByUser(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT plate_number FROM vehicles WHERE user_id = $1 ORDER BY id`, userID,
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
