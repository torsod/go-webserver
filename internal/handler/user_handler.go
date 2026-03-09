package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/service"
)

// UserHandler handles user CRUD endpoints
type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	id, err := h.users.Create(r.Context(), &user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	user.ID = id
	writeJSON(w, http.StatusCreated, user)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.users.Update(r.Context(), id, &user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated", "id": id})
}
