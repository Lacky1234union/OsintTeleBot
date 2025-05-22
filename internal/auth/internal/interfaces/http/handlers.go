package http

import (
	"encoding/json"
	"net/http"

	"github.com/russunion/OsintTeleBot/internal/auth/internal/domain"
)

type AuthHandler struct {
	userService  domain.UserService
	tokenService domain.TokenService
}

func NewAuthHandler(userService domain.UserService, tokenService domain.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := authResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := authResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	claims, err := h.tokenService.ValidateToken(token)
	if err != nil {
		if err == domain.ErrExpiredToken {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
