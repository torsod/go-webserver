package handler

import (
	"encoding/json"
	"net/http"

	"github.com/torsod/go-webserver/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type logoutRequest struct {
	Username string `json:"username"`
}

// Login authenticates a user and returns a token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	user, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not-authorized", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token":    user.Username, // Simple token = username for now
		"user":     user,
	})
}

// Logout logs out a user
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.auth.Logout(r.Context(), req.Username); err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
