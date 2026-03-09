package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/service"
)

// PublicAPIHandler handles public API endpoints
type PublicAPIHandler struct {
	offerings *service.OfferingService
	orderbook *service.OrderbookService
}

func NewPublicAPIHandler(offerings *service.OfferingService, orderbook *service.OrderbookService) *PublicAPIHandler {
	return &PublicAPIHandler{offerings: offerings, orderbook: orderbook}
}

// GetPublicOfferings returns all offerings with public fields
func (h *PublicAPIHandler) GetPublicOfferings(w http.ResponseWriter, r *http.Request) {
	offerings, err := h.offerings.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"offerings": offerings,
		"count":     len(offerings),
	})
}

// GetPublicOffering returns a single offering
func (h *PublicAPIHandler) GetPublicOffering(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	offering, err := h.offerings.FindBySymbol(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusNotFound, "not-found", "Offering not found: "+symbol)
		return
	}
	writeJSON(w, http.StatusOK, offering)
}

// GetPublicOrderbook returns the orderbook for a symbol
func (h *PublicAPIHandler) GetPublicOrderbook(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	summary, err := h.orderbook.GetSummary(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusNotFound, "not-found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// GetPublicSnapshot returns full snapshot of offerings + orderbooks
func (h *PublicAPIHandler) GetPublicSnapshot(w http.ResponseWriter, r *http.Request) {
	offerings, err := h.offerings.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	orderbooks := make(map[string]interface{})
	for _, o := range offerings {
		summary, err := h.orderbook.GetSummary(r.Context(), o.Symbol)
		if err == nil {
			orderbooks[o.Symbol] = summary
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"type":       "snapshot",
		"offerings":  offerings,
		"orderbooks": orderbooks,
	})
}
