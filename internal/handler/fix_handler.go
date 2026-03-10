package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/service"
)

// FIXHandler handles FIX protocol REST endpoints
type FIXHandler struct {
	fix *service.FIXService
}

// NewFIXHandler creates a new FIX handler
func NewFIXHandler(fix *service.FIXService) *FIXHandler {
	return &FIXHandler{fix: fix}
}

// StartSession starts a new FIX session
// POST /api/fix/session/start
func (h *FIXHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	var settings domain.FIXConnectionSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		// Use defaults if no body or invalid JSON
		settings = domain.DefaultFIXSettings()
	}

	// Apply defaults for zero values
	if settings.Port == 0 {
		defaults := domain.DefaultFIXSettings()
		settings.Port = defaults.Port
	}
	if settings.SenderCompID == "" {
		defaults := domain.DefaultFIXSettings()
		settings.SenderCompID = defaults.SenderCompID
	}
	if settings.TargetCompID == "" {
		defaults := domain.DefaultFIXSettings()
		settings.TargetCompID = defaults.TargetCompID
	}
	if settings.HeartbeatInterval == 0 {
		settings.HeartbeatInterval = 30
	}

	result, err := h.fix.StartSession(r.Context(), settings)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if !result.Success {
		writeError(w, http.StatusConflict, "fix-error", result.Message)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// StopSession stops the active FIX session
// POST /api/fix/session/stop
func (h *FIXHandler) StopSession(w http.ResponseWriter, r *http.Request) {
	result, err := h.fix.StopSession(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if !result.Success {
		writeError(w, http.StatusBadRequest, "fix-error", result.Message)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetStatus returns the current FIX session status
// GET /api/fix/session/status
func (h *FIXHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	session, err := h.fix.GetStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if session == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"connected": false,
			"message":   "No FIX session",
		})
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// SendOrder sends a single order via FIX
// POST /api/fix/orders
func (h *FIXHandler) SendOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.CSVOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if order.Symbol == "" || order.Quantity <= 0 || order.Price <= 0 {
		writeError(w, http.StatusBadRequest, "invalid-data", "Symbol, quantity, and price are required")
		return
	}

	result, err := h.fix.SendOrder(r.Context(), order)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if !result.Success {
		writeError(w, http.StatusBadRequest, "fix-error", result.Message)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// CancelOrder cancels an order by ClOrdID
// POST /api/fix/orders/cancel
func (h *FIXHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClOrdID string `json:"clOrdId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ClOrdID == "" {
		writeError(w, http.StatusBadRequest, "invalid-data", "clOrdId is required")
		return
	}

	result, err := h.fix.CancelOrder(r.Context(), req.ClOrdID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if !result.Success {
		writeError(w, http.StatusBadRequest, "fix-error", result.Message)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// SendOrdersFromCSV processes CSV content and sends orders with throttling
// POST /api/fix/orders/csv
func (h *FIXHandler) SendOrdersFromCSV(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Failed to read request body")
		return
	}

	csvContent := string(body)
	if csvContent == "" {
		writeError(w, http.StatusBadRequest, "invalid-data", "CSV content is required")
		return
	}

	result, err := h.fix.SendOrdersFromCSV(r.Context(), csvContent)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetLogs returns FIX logs
// GET /api/fix/logs?sessionId=&limit=
func (h *FIXHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	logs, err := h.fix.GetLogs(r.Context(), sessionID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if logs == nil {
		logs = []*domain.FIXLog{}
	}

	writeJSON(w, http.StatusOK, logs)
}

// GetOrders returns FIX orders for a session
// GET /api/fix/orders/{sessionId}
func (h *FIXHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	orders, err := h.fix.GetOrders(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	if orders == nil {
		orders = []*domain.FIXOrder{}
	}

	writeJSON(w, http.StatusOK, orders)
}

// ClearLogs deletes FIX logs
// DELETE /api/fix/logs?sessionId=
func (h *FIXHandler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")

	result, err := h.fix.ClearLogs(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fix-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
