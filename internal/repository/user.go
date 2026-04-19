package repository

import (
	"context"
	"database/sql"
	"errors"
	"expense-tracker/internal/model"
	"expense-tracker/internal/service"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*service.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email= $1
	`

	var u service.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser inserts a new user into the database
func (r *AuthRepo) CreateUser(email, passwordHash string) (*model.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, is_verified, created_at, updated_at
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail
func (r *AuthRepo) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, password_hash, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
