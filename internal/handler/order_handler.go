package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/service"
)

// OrderHandler handles order CRUD endpoints
type OrderHandler struct {
	orders *service.OrderService
}

func NewOrderHandler(orders *service.OrderService) *OrderHandler {
	return &OrderHandler{orders: orders}
}

type insertOrderRequest struct {
	Symbol        string           `json:"symbol"`
	Side          domain.OrderSide `json:"side"`
	OrderType     domain.OrderType `json:"orderType"`
	Quantity      int64            `json:"quantity"`
	Price         float64          `json:"price"`
	MinQty        *int64           `json:"minQty,omitempty"`
	Account       string           `json:"account,omitempty"`
	ExecInst      domain.ExecInst  `json:"execInst"`
	PriorityGroup int              `json:"priorityGroup"`
	UserID        string           `json:"userId"`
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req insertOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	order, err := h.orders.Insert(r.Context(), &req.Symbol, req.Side, req.OrderType,
		req.Quantity, req.Price, req.MinQty, req.Account, req.ExecInst,
		req.PriorityGroup, req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

type updateOrderRequest struct {
	Quantity *int64   `json:"quantity,omitempty"`
	Price    *float64 `json:"price,omitempty"`
	MinQty   *int64   `json:"minQty,omitempty"`
	UserID   string   `json:"userId"`
}

// UpdateOrder modifies an existing order
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	if err := h.orders.Update(r.Context(), id, req.Quantity, req.Price, req.MinQty, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated", "id": id})
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		UserID string `json:"userId"`
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := h.orders.Cancel(r.Context(), id, body.UserID, body.Reason); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "canceled", "id": id})
}

// CancelAllOrders cancels all orders for a user
func (h *OrderHandler) CancelAllOrders(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string `json:"userId"`
		Symbol string `json:"symbol"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid-data", "Invalid request body")
		return
	}

	count, err := h.orders.CancelAll(r.Context(), body.UserID, body.Symbol)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server-error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":   "canceled",
		"canceled": count,
	})
}
