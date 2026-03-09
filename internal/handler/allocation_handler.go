package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/service"
)

// AllocationHandler handles allocation endpoints
type AllocationHandler struct {
	allocation *service.AllocationService
}

func NewAllocationHandler(allocation *service.AllocationService) *AllocationHandler {
	return &AllocationHandler{allocation: allocation}
}

type tempCloseRequest struct {
	Symbol         string  `json:"symbol"`
	OfferingPrice  float64 `json:"offeringPrice"`
	Algorithm      string  `json:"algorithm"`
	LMUser         string  `json:"lmUser"`
}

// TempClose performs a temporary close (allocation)
func (h *AllocationHandler) TempClose(w http.ResponseWriter, r *http.Request) {
	var req tempCloseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	result, err := h.allocation.TempClose(r.Context(), req.Symbol, req.OfferingPrice, req.Algorithm, req.LMUser)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// BustAllocation busts an allocation session
func (h *AllocationHandler) BustAllocation(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	if err := h.allocation.Bust(r.Context(), sessionID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "busted", "sessionId": sessionID})
}

// ConfirmAllocation confirms an allocation session
func (h *AllocationHandler) ConfirmAllocation(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	if err := h.allocation.Confirm(r.Context(), sessionID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed", "sessionId": sessionID})
}

// CancelLeaves cancels unfilled quantities
func (h *AllocationHandler) CancelLeaves(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	if err := h.allocation.CancelLeaves(r.Context(), sessionID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "leaves canceled", "sessionId": sessionID})
}

// GetSessions returns allocation sessions
func (h *AllocationHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	sessions, err := h.allocation.GetSessions(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// GetSessionTrades returns trades for an allocation session
func (h *AllocationHandler) GetSessionTrades(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	trades, err := h.allocation.GetTrades(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"trades":    trades,
		"count":     len(trades),
		"sessionId": sessionID,
	})
}
