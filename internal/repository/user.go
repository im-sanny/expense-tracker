package repository

import (
	"database/sql"
	"expense-tracker/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(email, passwordHash string) (*model.User, error) {
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
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
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
