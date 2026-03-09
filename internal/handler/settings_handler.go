package handler

import (
	"encoding/json"
	"net/http"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/service"
)

// SettingsHandler handles settings endpoints
type SettingsHandler struct {
	settings *service.SettingsService
}

func NewSettingsHandler(settings *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{settings: settings}
}

// UpdateSettings updates system settings
func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.SystemSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.settings.Update(r.Context(), &settings); err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
