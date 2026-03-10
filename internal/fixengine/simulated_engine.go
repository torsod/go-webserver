package fixengine

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
)

// SimulatedEngine provides a FIX engine that simulates responses for development
type SimulatedEngine struct {
	connected atomic.Bool
	sessionID string
	execCb    ExecutionReportCallback
	mu        sync.Mutex
	stopCh    chan struct{}
	orderSeq  atomic.Int64
}

// NewSimulatedEngine creates a new simulated FIX engine
func NewSimulatedEngine() *SimulatedEngine {
	return &SimulatedEngine{}
}

func (e *SimulatedEngine) Start(ctx context.Context, settings domain.FIXConnectionSettings) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.connected.Load() {
		return fmt.Errorf("engine already connected")
	}

	e.sessionID = fmt.Sprintf("FIX.4.2:%s->%s", settings.SenderCompID, settings.TargetCompID)
	e.stopCh = make(chan struct{})

	slog.Info("[SIM] connecting to FIX engine",
		"host", settings.Host,
		"port", settings.Port,
		"sender", settings.SenderCompID,
		"target", settings.TargetCompID,
	)

	// Simulate connection delay
	time.Sleep(200 * time.Millisecond)

	e.connected.Store(true)
	slog.Info("[SIM] FIX session established", "sessionId", e.sessionID)

	// Start heartbeat goroutine
	go e.heartbeatLoop(settings.HeartbeatInterval)

	return nil
}

func (e *SimulatedEngine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.connected.Load() {
		return fmt.Errorf("engine not connected")
	}

	slog.Info("[SIM] disconnecting FIX session", "sessionId", e.sessionID)

	close(e.stopCh)
	e.connected.Store(false)

	// Simulate logout delay
	time.Sleep(100 * time.Millisecond)

	slog.Info("[SIM] FIX session disconnected")
	return nil
}

func (e *SimulatedEngine) IsConnected() bool {
	return e.connected.Load()
}

func (e *SimulatedEngine) SendNewOrderSingle(params SendOrderParams) error {
	if !e.connected.Load() {
		return fmt.Errorf("engine not connected")
	}

	slog.Info("[SIM] sending New Order Single",
		"clOrdId", params.ClOrdID,
		"symbol", params.Symbol,
		"side", params.Side,
		"qty", params.Quantity,
		"price", params.Price,
	)

	// Simulate async execution report after random delay (100-300ms)
	go func() {
		delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
		select {
		case <-time.After(delay):
		case <-e.stopCh:
			return
		}

		if !e.connected.Load() {
			return
		}

		seq := e.orderSeq.Add(1)
		report := &ExecutionReport{
			ClOrdID:   params.ClOrdID,
			OrderID:   fmt.Sprintf("SIM-ORD-%d", seq),
			ExecID:    fmt.Sprintf("SIM-EXEC-%d-%d", seq, time.Now().UnixMilli()),
			ExecType:  "0", // New
			OrdStatus: "0", // New
			Symbol:    params.Symbol,
			Side:      params.Side,
			LeavesQty: params.Quantity,
			CumQty:    0,
			AvgPx:     0,
		}

		if e.execCb != nil {
			e.execCb(context.Background(), report)
		}
	}()

	return nil
}

func (e *SimulatedEngine) SendOrderCancelRequest(origClOrdID, clOrdID, symbol, side string, quantity int64) error {
	if !e.connected.Load() {
		return fmt.Errorf("engine not connected")
	}

	slog.Info("[SIM] sending Order Cancel Request",
		"origClOrdId", origClOrdID,
		"clOrdId", clOrdID,
		"symbol", symbol,
	)

	// Simulate async cancel acknowledgment after random delay
	go func() {
		delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
		select {
		case <-time.After(delay):
		case <-e.stopCh:
			return
		}

		if !e.connected.Load() {
			return
		}

		seq := e.orderSeq.Add(1)
		report := &ExecutionReport{
			ClOrdID:     clOrdID,
			OrigClOrdID: origClOrdID,
			OrderID:     fmt.Sprintf("SIM-ORD-%d", seq),
			ExecID:      fmt.Sprintf("SIM-EXEC-%d-%d", seq, time.Now().UnixMilli()),
			ExecType:    "4", // Canceled
			OrdStatus:   "4", // Canceled
			Symbol:      symbol,
			Side:        side,
			LeavesQty:   0,
			CumQty:      0,
			AvgPx:       0,
		}

		if e.execCb != nil {
			e.execCb(context.Background(), report)
		}
	}()

	return nil
}

func (e *SimulatedEngine) SetExecutionReportCallback(cb ExecutionReportCallback) {
	e.execCb = cb
}

func (e *SimulatedEngine) SessionID() string {
	return e.sessionID
}

func (e *SimulatedEngine) heartbeatLoop(intervalSec int) {
	if intervalSec <= 0 {
		intervalSec = 30
	}
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if e.connected.Load() {
				slog.Debug("[SIM] heartbeat sent")
			}
		case <-e.stopCh:
			return
		}
	}
}
