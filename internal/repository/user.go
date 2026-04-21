package repository

import (
	"context"
	"database/sql"
	"errors"
	"expense-tracker/internal/model"
	"time"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email= $1
	`

	var u model.User
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
func (r *AuthRepo) CreateUser(ctx context.Context, email, passwordHash string) (*model.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, is_verified, created_at, updated_at
	`

	var u model.User
	err := r.db.QueryRowContext(ctx, query, email, passwordHash).Scan(
		&u.ID,
		&u.Email,
		&u.IsVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *AuthRepo) SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, tokenHash, expiresAt)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *AuthRepo) DeleteRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	query := `
		SELECT user_id FROM refresh_tokens
		WHERE token_hash =$1 AND expiresAt > NOW()
	`
	var userID string
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("token not found or expired")
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}
