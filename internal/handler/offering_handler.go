package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/service"
)

// OfferingHandler handles offering CRUD endpoints
type OfferingHandler struct {
	offerings *service.OfferingService
	state     *service.OfferingStateService
}

func NewOfferingHandler(offerings *service.OfferingService, state *service.OfferingStateService) *OfferingHandler {
	return &OfferingHandler{offerings: offerings, state: state}
}

// CreateOffering creates a new offering
func (h *OfferingHandler) CreateOffering(w http.ResponseWriter, r *http.Request) {
	var offering domain.Offering
	if err := json.NewDecoder(r.Body).Decode(&offering); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	id, err := h.offerings.Create(r.Context(), &offering)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	offering.ID = id
	writeJSON(w, http.StatusCreated, offering)
}

// UpdateOffering updates an existing offering
func (h *OfferingHandler) UpdateOffering(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var offering domain.Offering
	if err := json.NewDecoder(r.Body).Decode(&offering); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.offerings.Update(r.Context(), id, &offering); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated", "id": id})
}

type changeStateRequest struct {
	TargetState domain.OfferingState `json:"targetState"`
	Reason      string               `json:"reason"`
	UserID      string               `json:"userId"`
}

// ChangeState changes the state of an offering
func (h *OfferingHandler) ChangeState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req changeStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.state.ChangeState(r.Context(), id, req.TargetState, req.Reason, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-state", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "state changed", "id": id, "newState": string(req.TargetState)})
}
