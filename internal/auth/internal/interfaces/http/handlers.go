package http

import (
	"encoding/json"
	"net/http"

	"github.com/russunion/OsintTeleBot/internal/auth/internal/domain"
	"github.com/russunion/OsintTeleBot/internal/auth/internal/share/errs"
	"github.com/russunion/OsintTeleBot/internal/auth/internal/share/logger"
)

type AuthHandler struct {
	userService  domain.UserService
	tokenService domain.TokenService
	logger       *logger.Logger
}

func NewAuthHandler(userService domain.UserService, tokenService domain.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
		logger:       logger.New(logger.LevelInfo),
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
	ctx := r.Context()
	h.logger.Info(ctx, "Handling register request")

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if errs.IsAlreadyExists(err) {
			h.logger.Warn(ctx, "User already exists: %s", req.Username)
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		h.logger.Error(ctx, "Failed to register user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		h.logger.Error(ctx, "Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := authResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	h.logger.Info(ctx, "Successfully registered user: %s", user.Username)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Handling login request")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		if errs.IsUnauthorized(err) {
			h.logger.Warn(ctx, "Invalid credentials for user: %s", req.Username)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		h.logger.Error(ctx, "Failed to login user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		h.logger.Error(ctx, "Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := authResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	h.logger.Info(ctx, "Successfully logged in user: %s", user.Username)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Handling token validation request")

	token := r.Header.Get("Authorization")
	if token == "" {
		h.logger.Warn(ctx, "Missing authorization token")
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	claims, err := h.tokenService.ValidateToken(token)
	if err != nil {
		if errs.IsUnauthorized(err) {
			h.logger.Warn(ctx, "Token validation failed: %v", err)
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}
		h.logger.Error(ctx, "Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
	h.logger.Info(ctx, "Successfully validated token for user: %s", claims.Username)
}
