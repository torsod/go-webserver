package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/service"
)

// RESTHandler handles read-only REST API endpoints
type RESTHandler struct {
	offerings  *service.OfferingService
	orders     *service.OrderService
	orderbook  *service.OrderbookService
	trades     *service.TradeService
	settlement *service.SettlementService
	settings   *service.SettingsService
}

// NewRESTHandler creates a new REST handler
func NewRESTHandler(
	offerings *service.OfferingService,
	orders *service.OrderService,
	orderbook *service.OrderbookService,
	trades *service.TradeService,
	settlement *service.SettlementService,
	settings *service.SettingsService,
) *RESTHandler {
	return &RESTHandler{
		offerings:  offerings,
		orders:     orders,
		orderbook:  orderbook,
		trades:     trades,
		settlement: settlement,
		settings:   settings,
	}
}

// GetAPIDoc returns API documentation
func (h *RESTHandler) GetAPIDoc(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":    "PDP Go Webserver API",
		"version": "1.0.0",
		"endpoints": []map[string]string{
			{"method": "GET", "path": "/api/health", "description": "Health check"},
			{"method": "GET", "path": "/api/offerings", "description": "All offerings"},
			{"method": "GET", "path": "/api/orders", "description": "All orders"},
			{"method": "GET", "path": "/api/orders/symbol/{symbol}", "description": "Orders by symbol"},
			{"method": "GET", "path": "/api/orderbook/{symbol}", "description": "Orderbook for symbol"},
			{"method": "GET", "path": "/api/cprice/{symbol}", "description": "Clearing price for symbol"},
			{"method": "GET", "path": "/api/trades/{symbol}", "description": "Trades by symbol"},
			{"method": "GET", "path": "/api/settlement/{sessionId}", "description": "Settlement files"},
		},
	})
}

// GetOfferings returns all offerings with ETag support
func (h *RESTHandler) GetOfferings(w http.ResponseWriter, r *http.Request) {
	offerings, err := h.offerings.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	if checkETag(w, r, offerings) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"offerings": offerings,
		"count":     len(offerings),
	})
}

// GetOrders returns all orders with ETag support
func (h *RESTHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orders.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	if checkETag(w, r, orders) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// GetOrdersBySymbol returns orders for a specific symbol
func (h *RESTHandler) GetOrdersBySymbol(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	orders, err := h.orders.FindBySymbol(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
		"symbol": symbol,
	})
}

// GetOrderbook returns aggregated orderbook for a symbol
func (h *RESTHandler) GetOrderbook(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	summary, err := h.orderbook.GetSummary(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusNotFound, "not-found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// GetCprice returns clearing price for a symbol
func (h *RESTHandler) GetCprice(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	cprice, err := h.orderbook.GetCprice(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusNotFound, "not-found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cprice)
}

// GetTrades returns trades for a symbol
func (h *RESTHandler) GetTrades(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	trades, err := h.trades.FindBySymbol(r.Context(), symbol)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"trades": trades,
		"count":  len(trades),
		"symbol": symbol,
	})
}

// GetSettlement returns settlement files for an allocation session
func (h *RESTHandler) GetSettlement(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	result, err := h.settlement.GenerateSettlementFiles(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not-found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetSettings returns system settings
func (h *RESTHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settings.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)
}
