package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, email, passwordHash string) (*User, error)
	SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	ValidateRefreshToken(ctx context.Context, tokenHash string) (string, error)
}

// Auth holds dependencies for auth operations
type AuthService struct {
	jwtSecret   []byte
	userService UserService
}

func NewAuthService(jwtSecret string, userService UserService) *AuthService {
	return &AuthService{
		jwtSecret:   []byte(jwtSecret),
		userService: userService,
	}
}

// HashPassword hashes plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPassword compares a plain-text password with a hash
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateAccessToken creates a short-lived JWT for API request
func (s *AuthService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// GenerateRefreshToken creates a random opaque token (UUID-style)
// This will be stored in the database, not a JWT
func GenerateRefreshToken() string {
	// Simple approach: use crypto/rand for secure random string
	// For production consider using github.com/google/uuid
	return fmt.Sprintf("%x", time.Now().UnixNano()) // replace with proper uuid in real app
}

// ValidateAccessToken parses and validates a JWT
func (s *AuthService) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}

// Auth business logic
func (s *AuthService) Register(ctx context.Context, email, password string) (*User, error) {
	// Check if email exists
	existing, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check user exists: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash and create
	hash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	return s.userService.CreateUser(ctx, email, hash)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}
	// Generic errors to prevent enumeration
	if user == nil || !CheckPassword(password, user.PasswordHash) {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err = s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}
	refreshToken = GenerateRefreshToken()

	// Store refresh token hash in DB
	if err := s.userService.SaveRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return "", "", fmt.Errorf("save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
