package handler

import (
	"encoding/json"
	"expense-tracker/internal/repository"
	"expense-tracker/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
}

func NewAuthHandler(auth *service.AuthService, userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		userRepo:    userRepo,
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	// check if user already exists
	existing, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// hash password
	hash, err := service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// create user
	user, err := h.userRepo.CreateUser(req.Email, hash)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// return success without password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user
	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generic error message to prevent user enumeration
	if user == nil || !service.CheckPassword(req.Password, user.PasswordHash) {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Generate tokens
	accessToken, err := h.authService.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken := service.GenerateRefreshToken()
	// save refreshToken hash to database

	// set secure cookies
	setSecureCookies(w, "access_token", accessToken, 15*60)
	setSecureCookies(w, "refresh_token", refreshToken, 7*24*60*60)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}

// Helper: Set secure HttpOnly cookie
func setSecureCookies(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production (requires https)
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
}
