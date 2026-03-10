package fixengine

import (
	"context"
	"crypto/md5"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/config"
	"github.com/shopspring/decimal"
	"github.com/torsod/go-webserver/internal/domain"
)

// FIX 4.2 tag numbers
const (
	tagMsgType      quickfix.Tag = 35
	tagClOrdID      quickfix.Tag = 11
	tagOrigClOrdID  quickfix.Tag = 41
	tagOrderID      quickfix.Tag = 37
	tagExecID       quickfix.Tag = 17
	tagExecType     quickfix.Tag = 150
	tagOrdStatus    quickfix.Tag = 39
	tagSymbol       quickfix.Tag = 55
	tagSide         quickfix.Tag = 54
	tagOrderQty     quickfix.Tag = 38
	tagOrdType      quickfix.Tag = 40
	tagPrice        quickfix.Tag = 44
	tagHandlInst    quickfix.Tag = 21
	tagTransactTime quickfix.Tag = 60
	tagAccount      quickfix.Tag = 1
	tagTimeInForce  quickfix.Tag = 59
	tagText         quickfix.Tag = 58
	tagLeavesQty    quickfix.Tag = 151
	tagCumQty       quickfix.Tag = 14
	tagAvgPx        quickfix.Tag = 6
	tagLastShares   quickfix.Tag = 32
	tagLastPx       quickfix.Tag = 31
	tagOrdRejReason quickfix.Tag = 103
	tagUsername      quickfix.Tag = 553
	tagPassword     quickfix.Tag = 554
)

// RealEngine wraps quickfixgo Initiator for FIX 4.2 connections
type RealEngine struct {
	initiator *quickfix.Initiator
	sessionID quickfix.SessionID
	settings  domain.FIXConnectionSettings
	connected atomic.Bool
	execCb    ExecutionReportCallback
	mu        sync.Mutex
}

// NewRealEngine creates a new quickfixgo-based FIX engine
func NewRealEngine() *RealEngine {
	return &RealEngine{}
}

// quickfix.Application interface implementation

func (e *RealEngine) OnCreate(sessionID quickfix.SessionID) {
	slog.Info("FIX session created", "sessionId", sessionID.String())
	e.sessionID = sessionID
}

func (e *RealEngine) OnLogon(sessionID quickfix.SessionID) {
	slog.Info("FIX session logged on", "sessionId", sessionID.String())
	e.connected.Store(true)
}

func (e *RealEngine) OnLogout(sessionID quickfix.SessionID) {
	slog.Info("FIX session logged out", "sessionId", sessionID.String())
	e.connected.Store(false)
}

func (e *RealEngine) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, _ := msg.Header.GetString(tagMsgType)
	if msgType == "A" { // Logon
		if e.settings.Username != "" {
			msg.Body.SetString(tagUsername, e.settings.Username)
		}
		if e.settings.Password != "" {
			hash := md5.Sum([]byte(e.settings.Password))
			msg.Body.SetString(tagPassword, fmt.Sprintf("%x", hash))
		}
		slog.Info("FIX Logon message prepared",
			"sender", e.settings.SenderCompID,
			"target", e.settings.TargetCompID,
		)
	}
}

func (e *RealEngine) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	msgType, _ := msg.Header.GetString(tagMsgType)
	slog.Debug("FIX admin message received", "msgType", msgType)
	return nil
}

func (e *RealEngine) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	msgType, _ := msg.Header.GetString(tagMsgType)
	slog.Debug("FIX app message sending", "msgType", msgType)
	return nil
}

func (e *RealEngine) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	msgType, _ := msg.Header.GetString(tagMsgType)

	switch msgType {
	case "8": // Execution Report
		return e.onExecutionReport(msg)
	default:
		slog.Warn("FIX unexpected message type", "msgType", msgType)
	}

	return nil
}

func (e *RealEngine) onExecutionReport(msg *quickfix.Message) quickfix.MessageRejectError {
	report := &ExecutionReport{}

	if v, err := msg.Body.GetString(tagClOrdID); err == nil {
		report.ClOrdID = v
	}
	if v, err := msg.Body.GetString(tagOrigClOrdID); err == nil {
		report.OrigClOrdID = v
	}
	if v, err := msg.Body.GetString(tagOrderID); err == nil {
		report.OrderID = v
	}
	if v, err := msg.Body.GetString(tagExecID); err == nil {
		report.ExecID = v
	}
	if v, err := msg.Body.GetString(tagExecType); err == nil {
		report.ExecType = v
	}
	if v, err := msg.Body.GetString(tagOrdStatus); err == nil {
		report.OrdStatus = v
	}
	if v, err := msg.Body.GetString(tagSymbol); err == nil {
		report.Symbol = v
	}
	if v, err := msg.Body.GetString(tagSide); err == nil {
		report.Side = v
	}
	if v, err := msg.Body.GetString(tagLeavesQty); err == nil {
		if qty, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
			report.LeavesQty = qty
		}
	}
	if v, err := msg.Body.GetString(tagCumQty); err == nil {
		if qty, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
			report.CumQty = qty
		}
	}
	if v, err := msg.Body.GetString(tagAvgPx); err == nil {
		if px, parseErr := strconv.ParseFloat(v, 64); parseErr == nil {
			report.AvgPx = px
		}
	}
	if v, err := msg.Body.GetString(tagLastShares); err == nil {
		if qty, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
			report.LastShares = qty
		}
	}
	if v, err := msg.Body.GetString(tagLastPx); err == nil {
		if px, parseErr := strconv.ParseFloat(v, 64); parseErr == nil {
			report.LastPx = px
		}
	}
	if v, err := msg.Body.GetInt(tagOrdRejReason); err == nil {
		report.OrdRejReason = v
	}
	if v, err := msg.Body.GetString(tagText); err == nil {
		report.Text = v
	}

	slog.Info("FIX execution report received",
		"clOrdId", report.ClOrdID,
		"execType", report.ExecType,
		"ordStatus", report.OrdStatus,
		"symbol", report.Symbol,
	)

	if e.execCb != nil {
		e.execCb(context.Background(), report)
	}

	return nil
}

// Engine interface implementation

func (e *RealEngine) Start(ctx context.Context, settings domain.FIXConnectionSettings) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.connected.Load() {
		return fmt.Errorf("engine already connected")
	}

	e.settings = settings

	// Build quickfix settings programmatically
	qfSettings := quickfix.NewSettings()

	global := qfSettings.GlobalSettings()
	global.Set(config.BeginString, quickfix.BeginStringFIX42)
	global.Set(config.SenderCompID, settings.SenderCompID)
	global.Set(config.TargetCompID, settings.TargetCompID)
	global.Set(config.SocketConnectHost, settings.Host)
	global.Set(config.SocketConnectPort, fmt.Sprintf("%d", settings.Port))
	global.Set(config.HeartBtInt, fmt.Sprintf("%d", settings.HeartbeatInterval))
	global.Set(config.ResetOnLogon, "Y")
	global.Set(config.ReconnectInterval, "10")
	global.Set("ConnectionType", "initiator")
	global.Set(config.StartTime, "00:00:00")
	global.Set(config.EndTime, "00:00:00")

	// Add session
	sessionSettings := quickfix.NewSessionSettings()
	_, err := qfSettings.AddSession(sessionSettings)
	if err != nil {
		return fmt.Errorf("add FIX session: %w", err)
	}

	// Create initiator with in-memory store and screen logging
	initiator, err := quickfix.NewInitiator(
		e,
		quickfix.NewMemoryStoreFactory(),
		qfSettings,
		quickfix.NewNullLogFactory(),
	)
	if err != nil {
		return fmt.Errorf("create FIX initiator: %w", err)
	}
	e.initiator = initiator

	slog.Info("starting FIX initiator",
		"host", settings.Host,
		"port", settings.Port,
		"sender", settings.SenderCompID,
		"target", settings.TargetCompID,
	)

	if err := e.initiator.Start(); err != nil {
		return fmt.Errorf("start FIX initiator: %w", err)
	}

	return nil
}

func (e *RealEngine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.initiator != nil {
		slog.Info("stopping FIX initiator")
		e.initiator.Stop()
		e.initiator = nil
	}

	e.connected.Store(false)
	return nil
}

func (e *RealEngine) IsConnected() bool {
	return e.connected.Load()
}

func (e *RealEngine) SendNewOrderSingle(params SendOrderParams) error {
	if !e.connected.Load() {
		return fmt.Errorf("engine not connected")
	}

	msg := quickfix.NewMessage()
	msg.Header.SetString(tagMsgType, "D") // New Order Single

	msg.Body.SetString(tagClOrdID, params.ClOrdID)
	msg.Body.SetString(tagHandlInst, "2") // Automated execution, public
	msg.Body.SetString(tagSymbol, params.Symbol)
	msg.Body.SetString(tagSide, params.Side)
	msg.Body.SetString(tagTransactTime, time.Now().UTC().Format("20060102-15:04:05"))
	msg.Body.SetField(tagOrderQty, quickfix.FIXDecimal{Decimal: decimal.NewFromInt(params.Quantity), Scale: 0})
	msg.Body.SetString(tagOrdType, "2") // Limit
	msg.Body.SetField(tagPrice, quickfix.FIXDecimal{Decimal: decimal.NewFromFloat(params.Price), Scale: 4})
	msg.Body.SetString(tagTimeInForce, "1") // GTC

	if params.Account != "" {
		msg.Body.SetString(tagAccount, params.Account)
	}
	if params.Text != "" {
		msg.Body.SetString(tagText, params.Text)
	}

	slog.Info("sending New Order Single",
		"clOrdId", params.ClOrdID,
		"symbol", params.Symbol,
		"side", params.Side,
		"qty", params.Quantity,
		"price", params.Price,
	)

	return quickfix.SendToTarget(msg, e.sessionID)
}

func (e *RealEngine) SendOrderCancelRequest(origClOrdID, clOrdID, symbol, side string, quantity int64) error {
	if !e.connected.Load() {
		return fmt.Errorf("engine not connected")
	}

	msg := quickfix.NewMessage()
	msg.Header.SetString(tagMsgType, "F") // Order Cancel Request

	msg.Body.SetString(tagOrigClOrdID, origClOrdID)
	msg.Body.SetString(tagClOrdID, clOrdID)
	msg.Body.SetString(tagSymbol, symbol)
	msg.Body.SetString(tagSide, side)
	msg.Body.SetString(tagTransactTime, time.Now().UTC().Format("20060102-15:04:05"))
	msg.Body.SetField(tagOrderQty, quickfix.FIXDecimal{Decimal: decimal.NewFromInt(quantity), Scale: 0})

	slog.Info("sending Order Cancel Request",
		"origClOrdId", origClOrdID,
		"clOrdId", clOrdID,
		"symbol", symbol,
	)

	return quickfix.SendToTarget(msg, e.sessionID)
}

func (e *RealEngine) SetExecutionReportCallback(cb ExecutionReportCallback) {
	e.execCb = cb
}

func (e *RealEngine) SessionID() string {
	return e.sessionID.String()
}
