package fixengine

import (
	"context"

	"github.com/torsod/go-webserver/internal/domain"
)

// ExecutionReportCallback is called when the engine receives an execution report
type ExecutionReportCallback func(ctx context.Context, report *ExecutionReport)

// ExecutionReport holds parsed execution report data
type ExecutionReport struct {
	ClOrdID      string  `json:"clOrdId"`
	OrigClOrdID  string  `json:"origClOrdId,omitempty"`
	OrderID      string  `json:"orderId"`
	ExecID       string  `json:"execId"`
	ExecType     string  `json:"execType"`   // "0"=New, "1"=Partial, "2"=Fill, "4"=Canceled, "8"=Rejected
	OrdStatus    string  `json:"ordStatus"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	LeavesQty    int64   `json:"leavesQty"`
	CumQty       int64   `json:"cumQty"`
	AvgPx        float64 `json:"avgPx"`
	LastShares   int64   `json:"lastShares"`
	LastPx       float64 `json:"lastPx"`
	OrdRejReason int     `json:"ordRejReason,omitempty"`
	Text         string  `json:"text,omitempty"`
}

// SendOrderParams contains parameters for sending a new order
type SendOrderParams struct {
	ClOrdID       string
	Symbol        string
	Side          string // "1"=Buy, "2"=Sell
	Quantity      int64
	Price         float64
	Account       string
	PriorityGroup int
	Text          string
}

// Engine defines the FIX engine interface (real or simulated)
type Engine interface {
	// Start initializes and connects the engine
	Start(ctx context.Context, settings domain.FIXConnectionSettings) error
	// Stop gracefully disconnects
	Stop() error
	// IsConnected returns connection status
	IsConnected() bool
	// SendNewOrderSingle sends an order via FIX
	SendNewOrderSingle(params SendOrderParams) error
	// SendOrderCancelRequest sends a cancel request via FIX
	SendOrderCancelRequest(origClOrdID, clOrdID, symbol, side string, quantity int64) error
	// SetExecutionReportCallback registers the callback for incoming execution reports
	SetExecutionReportCallback(cb ExecutionReportCallback)
	// SessionID returns the FIX session identifier string
	SessionID() string
}
