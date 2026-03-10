package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/fixengine"
	"github.com/torsod/go-webserver/internal/store"
)

// FIXSessionResult is the response for session operations
type FIXSessionResult struct {
	SessionID string `json:"sessionId,omitempty"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// FIXOrderResult is the response for order operations
type FIXOrderResult struct {
	Success     bool   `json:"success"`
	ClOrdID     string `json:"clOrdId,omitempty"`
	MainOrderID string `json:"mainOrderId,omitempty"`
	Message     string `json:"message"`
}

// FIXCSVResult is the response for CSV bulk upload
type FIXCSVResult struct {
	Success         bool   `json:"success"`
	OrdersSubmitted int    `json:"ordersSubmitted"`
	Message         string `json:"message"`
}

// FIXService handles FIX session and order management
type FIXService struct {
	engine      fixengine.Engine
	sessions    store.FIXSessionStore
	fixOrders   store.FIXOrderStore
	fixLogs     store.FIXLogStore
	mainOrders  store.OrderStore
	offerings   store.OfferingStore

	currentDBSessionID string // database row ID for the active session
	mu                 sync.RWMutex
}

// NewFIXService creates a new FIX service
func NewFIXService(
	engine fixengine.Engine,
	sessions store.FIXSessionStore,
	fixOrders store.FIXOrderStore,
	fixLogs store.FIXLogStore,
	mainOrders store.OrderStore,
	offerings store.OfferingStore,
) *FIXService {
	return &FIXService{
		engine:     engine,
		sessions:   sessions,
		fixOrders:  fixOrders,
		fixLogs:    fixLogs,
		mainOrders: mainOrders,
		offerings:  offerings,
	}
}

// CleanupOrphanedSessions marks any sessions still showing connected=true as disconnected
func (s *FIXService) CleanupOrphanedSessions(ctx context.Context) {
	session, err := s.sessions.FindActive(ctx)
	if err != nil {
		return // no active sessions found
	}
	if session != nil {
		slog.Info("cleaning up orphaned FIX session", "sessionId", session.SessionID)
		now := time.Now()
		_ = s.sessions.UpdateFields(ctx, session.ID, map[string]interface{}{
			"connected":  false,
			"updated_at": now,
		})
	}
}

// StartSession starts a new FIX session
func (s *FIXService) StartSession(ctx context.Context, settings domain.FIXConnectionSettings) (*FIXSessionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.engine.IsConnected() {
		return &FIXSessionResult{Success: false, Message: "A FIX session is already active"}, nil
	}

	// Generate session ID
	sessionID := fmt.Sprintf("FIX-%d-%04d", time.Now().UnixMilli(), time.Now().Nanosecond()%10000)

	// Create session record in DB
	session := &domain.FIXSession{
		SessionID:    sessionID,
		Host:         settings.Host,
		Port:         settings.Port,
		SenderCompID: settings.SenderCompID,
		TargetCompID: settings.TargetCompID,
		Connected:    false,
		Simulated:    settings.Simulated,
	}

	dbID, err := s.sessions.Insert(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("insert FIX session: %w", err)
	}
	s.currentDBSessionID = dbID

	// Log session initialization
	s.logMessage(ctx, sessionID, "INFO", "OUT",
		fmt.Sprintf("Initializing FIX session: %s -> %s @ %s:%d",
			settings.SenderCompID, settings.TargetCompID, settings.Host, settings.Port),
		"")

	// Set the execution report callback before starting
	s.engine.SetExecutionReportCallback(s.handleExecutionReport)

	// Start the engine
	if err := s.engine.Start(ctx, settings); err != nil {
		s.logMessage(ctx, sessionID, "ERROR", "",
			fmt.Sprintf("Failed to start FIX session: %v", err), "")
		return nil, fmt.Errorf("start FIX engine: %w", err)
	}

	// Update session as connected
	now := time.Now()
	_ = s.sessions.UpdateFields(ctx, dbID, map[string]interface{}{
		"connected":  true,
		"updated_at": now,
	})

	s.logMessage(ctx, sessionID, "INFO", "IN",
		"FIX session established (Logon acknowledged)", "")

	slog.Info("FIX session started", "sessionId", sessionID, "simulated", settings.Simulated)

	return &FIXSessionResult{
		SessionID: sessionID,
		Success:   true,
		Message:   "FIX session started successfully",
	}, nil
}

// StopSession stops the active FIX session
func (s *FIXService) StopSession(ctx context.Context) (*FIXSessionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentDBSessionID == "" {
		return &FIXSessionResult{Success: false, Message: "No active FIX session"}, nil
	}

	// Get current session for logging
	session, err := s.sessions.FindByID(ctx, s.currentDBSessionID)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}

	s.logMessage(ctx, session.SessionID, "INFO", "OUT",
		"Sending Logout (MsgType=5)", "")

	// Stop the engine
	if err := s.engine.Stop(); err != nil {
		slog.Warn("error stopping FIX engine", "error", err)
	}

	// Update session as disconnected
	now := time.Now()
	_ = s.sessions.UpdateFields(ctx, s.currentDBSessionID, map[string]interface{}{
		"connected":  false,
		"updated_at": now,
	})

	s.logMessage(ctx, session.SessionID, "INFO", "IN",
		"FIX session disconnected (Logout complete)", "")

	slog.Info("FIX session stopped", "sessionId", session.SessionID)

	s.currentDBSessionID = ""

	return &FIXSessionResult{
		SessionID: session.SessionID,
		Success:   true,
		Message:   "FIX session stopped successfully",
	}, nil
}

// GetStatus returns the current session status
func (s *FIXService) GetStatus(ctx context.Context) (*domain.FIXSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentDBSessionID == "" {
		// Check for most recent session (may be disconnected)
		session, err := s.sessions.FindActive(ctx)
		if err != nil {
			return nil, nil // no sessions at all
		}
		return session, nil
	}

	session, err := s.sessions.FindByID(ctx, s.currentDBSessionID)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}

	// Update connected status from engine
	session.Connected = s.engine.IsConnected()
	return session, nil
}

// SendOrder sends a single order via FIX
func (s *FIXService) SendOrder(ctx context.Context, order domain.CSVOrder) (*FIXOrderResult, error) {
	s.mu.RLock()
	sessionID := s.currentDBSessionID
	s.mu.RUnlock()

	if sessionID == "" || !s.engine.IsConnected() {
		return &FIXOrderResult{Success: false, Message: "No active FIX session"}, nil
	}

	// Get current session for logging
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}

	// Generate ClOrdID
	clOrdID := fmt.Sprintf("ORD-%d-%06d", time.Now().UnixMilli(), time.Now().Nanosecond()%1000000)

	// Determine side
	side := domain.OrderSideBid
	fixSide := "1" // Buy
	if order.Side == "OFFER" || order.Side == "2" {
		side = domain.OrderSideOffer
		fixSide = "2" // Sell
	}

	// Determine order type
	orderType := domain.OrderTypeCompetitive
	if order.OrderType != "" {
		orderType = domain.OrderType(order.OrderType)
	}

	// Create the main order record
	now := time.Now()
	seq, _ := s.mainOrders.NextSequence(ctx)
	mainOrder := &domain.Order{
		Symbol:            order.Symbol,
		Side:              side,
		OrderType:         orderType,
		Quantity:          order.Quantity,
		Price:             order.Price,
		PriorityGroup:     order.PriorityGroup,
		Account:           order.Account,
		MinQty:            order.MinQty,
		UserID:            order.BDUser,
		Status:            domain.OrderStatusActive,
		Timestamp:         now,
		OriginalEntryTime: now,
		OrderSequence:     seq,
	}
	if order.ExecInst != "" {
		mainOrder.ExecInst = domain.ExecInst(order.ExecInst)
	}

	mainOrderID, err := s.mainOrders.Insert(ctx, mainOrder)
	if err != nil {
		return nil, fmt.Errorf("insert main order: %w", err)
	}

	// Create FIX order record
	fixOrder := &domain.FIXOrder{
		ClOrdID:      clOrdID,
		SessionID:    session.SessionID,
		Symbol:       order.Symbol,
		Side:         fixSide,
		Quantity:     order.Quantity,
		Price:        order.Price,
		OrdType:      "Limit",
		TimeInForce:  "GTC",
		Account:      order.Account,
		MainOrderID:  mainOrderID,
		Status:       domain.FIXOrderStatusPending,
		TransactTime: now,
	}

	_, err = s.fixOrders.Insert(ctx, fixOrder)
	if err != nil {
		return nil, fmt.Errorf("insert fix order: %w", err)
	}

	// Build priority text
	priorityText := ""
	if order.PriorityGroup > 0 {
		priorityText = fmt.Sprintf("PriorityGroup=%d", order.PriorityGroup)
	}

	// Log outgoing message
	s.logMessage(ctx, session.SessionID, "INFO", "OUT",
		fmt.Sprintf("New Order Single (D): ClOrdID=%s Symbol=%s Side=%s Qty=%d Price=%.4f Account=%s",
			clOrdID, order.Symbol, fixSide, order.Quantity, order.Price, order.Account),
		fmt.Sprintf("35=D|11=%s|21=2|55=%s|54=%s|38=%d|40=2|44=%.4f|59=1|1=%s|58=%s",
			clOrdID, order.Symbol, fixSide, order.Quantity, order.Price, order.Account, priorityText),
	)

	// Increment messages sent
	_ = s.sessions.UpdateFields(ctx, sessionID, map[string]interface{}{
		"messages_sent": session.MessagesSent + 1,
		"updated_at":    now,
	})

	// Send via engine
	params := fixengine.SendOrderParams{
		ClOrdID:       clOrdID,
		Symbol:        order.Symbol,
		Side:          fixSide,
		Quantity:      order.Quantity,
		Price:         order.Price,
		Account:       order.Account,
		PriorityGroup: order.PriorityGroup,
		Text:          priorityText,
	}

	if err := s.engine.SendNewOrderSingle(params); err != nil {
		return nil, fmt.Errorf("send order: %w", err)
	}

	return &FIXOrderResult{
		Success:     true,
		ClOrdID:     clOrdID,
		MainOrderID: mainOrderID,
		Message:     "Order submitted successfully",
	}, nil
}

// CancelOrder cancels an order by ClOrdID
func (s *FIXService) CancelOrder(ctx context.Context, clOrdID string) (*FIXOrderResult, error) {
	s.mu.RLock()
	sessionID := s.currentDBSessionID
	s.mu.RUnlock()

	if sessionID == "" || !s.engine.IsConnected() {
		return &FIXOrderResult{Success: false, Message: "No active FIX session"}, nil
	}

	// Find the FIX order
	fixOrder, err := s.fixOrders.FindByClOrdID(ctx, clOrdID)
	if err != nil {
		return &FIXOrderResult{Success: false, Message: fmt.Sprintf("Order not found: %s", clOrdID)}, nil
	}

	// Validate status is cancelable
	if fixOrder.Status == domain.FIXOrderStatusFilled ||
		fixOrder.Status == domain.FIXOrderStatusCanceled ||
		fixOrder.Status == domain.FIXOrderStatusRejected {
		return &FIXOrderResult{
			Success: false,
			Message: fmt.Sprintf("Cannot cancel order in status: %s", fixOrder.Status),
		}, nil
	}

	// Get session for logging
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}

	// Generate cancel ClOrdID
	cancelClOrdID := fmt.Sprintf("CXL-%d-%04d", time.Now().UnixMilli(), time.Now().Nanosecond()%10000)

	// Cancel the main order if linked
	if fixOrder.MainOrderID != "" {
		now := time.Now()
		_ = s.mainOrders.UpdateFields(ctx, fixOrder.MainOrderID, map[string]interface{}{
			"status":      string(domain.OrderStatusCanceled),
			"canceled_at": now,
			"updated_at":  now,
		})
	}

	// Log outgoing cancel
	s.logMessage(ctx, session.SessionID, "INFO", "OUT",
		fmt.Sprintf("Order Cancel Request (F): OrigClOrdID=%s ClOrdID=%s Symbol=%s",
			clOrdID, cancelClOrdID, fixOrder.Symbol),
		fmt.Sprintf("35=F|41=%s|11=%s|55=%s|54=%s|38=%d",
			clOrdID, cancelClOrdID, fixOrder.Symbol, fixOrder.Side, fixOrder.Quantity),
	)

	// Send cancel via engine
	if err := s.engine.SendOrderCancelRequest(
		clOrdID, cancelClOrdID, fixOrder.Symbol, fixOrder.Side, fixOrder.Quantity,
	); err != nil {
		return nil, fmt.Errorf("send cancel: %w", err)
	}

	return &FIXOrderResult{
		Success: true,
		ClOrdID: cancelClOrdID,
		Message: "Cancel request submitted",
	}, nil
}

// SendOrdersFromCSV parses CSV content and sends orders with throttling
func (s *FIXService) SendOrdersFromCSV(ctx context.Context, csvContent string) (*FIXCSVResult, error) {
	orders, err := ParseOrdersCSV(csvContent)
	if err != nil {
		return &FIXCSVResult{Success: false, Message: err.Error()}, nil
	}

	if len(orders) == 0 {
		return &FIXCSVResult{Success: false, Message: "No valid orders found in CSV"}, nil
	}

	slog.Info("processing CSV orders", "count", len(orders))

	submitted := 0
	for i, order := range orders {
		result, err := s.SendOrder(ctx, order)
		if err != nil {
			slog.Warn("CSV order failed", "index", i, "error", err)
			continue
		}
		if result.Success {
			submitted++
		}

		// Throttle: 50ms delay between orders (max ~20 orders/sec)
		if i < len(orders)-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return &FIXCSVResult{
		Success:         true,
		OrdersSubmitted: submitted,
		Message:         fmt.Sprintf("Submitted %d of %d orders", submitted, len(orders)),
	}, nil
}

// GetLogs returns FIX logs for a session
func (s *FIXService) GetLogs(ctx context.Context, sessionID string, limit int) ([]*domain.FIXLog, error) {
	if limit <= 0 {
		limit = 100
	}
	if sessionID == "" {
		// Use current session
		s.mu.RLock()
		dbID := s.currentDBSessionID
		s.mu.RUnlock()
		if dbID != "" {
			session, err := s.sessions.FindByID(ctx, dbID)
			if err == nil {
				sessionID = session.SessionID
			}
		}
	}
	if sessionID == "" {
		return nil, nil
	}
	return s.fixLogs.FindBySessionID(ctx, sessionID, limit)
}

// GetOrders returns FIX orders for a session
func (s *FIXService) GetOrders(ctx context.Context, sessionID string) ([]*domain.FIXOrder, error) {
	if sessionID == "" {
		s.mu.RLock()
		dbID := s.currentDBSessionID
		s.mu.RUnlock()
		if dbID != "" {
			session, err := s.sessions.FindByID(ctx, dbID)
			if err == nil {
				sessionID = session.SessionID
			}
		}
	}
	if sessionID == "" {
		return nil, nil
	}
	return s.fixOrders.FindBySessionID(ctx, sessionID)
}

// ClearLogs deletes logs for a session
func (s *FIXService) ClearLogs(ctx context.Context, sessionID string) (*FIXSessionResult, error) {
	if sessionID == "" {
		s.mu.RLock()
		dbID := s.currentDBSessionID
		s.mu.RUnlock()
		if dbID != "" {
			session, err := s.sessions.FindByID(ctx, dbID)
			if err == nil {
				sessionID = session.SessionID
			}
		}
	}
	if sessionID == "" {
		return &FIXSessionResult{Success: false, Message: "No session specified"}, nil
	}

	deleted, err := s.fixLogs.DeleteBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("clear logs: %w", err)
	}

	return &FIXSessionResult{
		Success: true,
		Message: fmt.Sprintf("Cleared %d log entries", deleted),
	}, nil
}

// handleExecutionReport is the callback for incoming execution reports from the engine
func (s *FIXService) handleExecutionReport(ctx context.Context, report *fixengine.ExecutionReport) {
	slog.Info("handling execution report",
		"clOrdId", report.ClOrdID,
		"execType", report.ExecType,
	)

	// Determine the original ClOrdID to look up
	lookupClOrdID := report.ClOrdID
	if report.OrigClOrdID != "" {
		lookupClOrdID = report.OrigClOrdID
	}

	// Find the FIX order
	fixOrder, err := s.fixOrders.FindByClOrdID(ctx, lookupClOrdID)
	if err != nil {
		slog.Warn("execution report for unknown order", "clOrdId", lookupClOrdID, "error", err)
		return
	}

	// Map ExecType to status
	var newStatus string
	switch report.ExecType {
	case "0": // New
		newStatus = domain.FIXOrderStatusNew
	case "1": // Partial Fill
		newStatus = domain.FIXOrderStatusPartial
	case "2": // Fill
		newStatus = domain.FIXOrderStatusFilled
	case "4": // Canceled
		newStatus = domain.FIXOrderStatusCanceled
	case "8": // Rejected
		newStatus = domain.FIXOrderStatusRejected
	default:
		slog.Warn("unknown ExecType", "execType", report.ExecType)
		return
	}

	// Update FIX order status
	_ = s.fixOrders.UpdateFields(ctx, fixOrder.ID, map[string]interface{}{
		"status": newStatus,
	})

	// Log the incoming execution report
	s.mu.RLock()
	dbID := s.currentDBSessionID
	s.mu.RUnlock()

	if dbID != "" {
		session, err := s.sessions.FindByID(ctx, dbID)
		if err == nil {
			s.logMessage(ctx, session.SessionID, "INFO", "IN",
				fmt.Sprintf("Execution Report (8): ClOrdID=%s ExecType=%s OrdStatus=%s Symbol=%s CumQty=%d AvgPx=%.4f",
					report.ClOrdID, report.ExecType, report.OrdStatus, report.Symbol, report.CumQty, report.AvgPx),
				fmt.Sprintf("35=8|11=%s|17=%s|150=%s|39=%s|55=%s|151=%d|14=%d|6=%.4f",
					report.ClOrdID, report.ExecID, report.ExecType, report.OrdStatus,
					report.Symbol, report.LeavesQty, report.CumQty, report.AvgPx),
			)

			// Increment messages received
			_ = s.sessions.UpdateFields(ctx, dbID, map[string]interface{}{
				"messages_received": session.MessagesRecv + 1,
				"updated_at":        time.Now(),
			})
		}
	}
}

// logMessage inserts a log entry for a FIX session
func (s *FIXService) logMessage(ctx context.Context, sessionID, level, direction, message, rawData string) {
	log := &domain.FIXLog{
		SessionID: sessionID,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Direction: direction,
		RawData:   rawData,
	}
	if _, err := s.fixLogs.Insert(ctx, log); err != nil {
		slog.Warn("failed to insert FIX log", "error", err)
	}
}

// ParseOrdersCSV parses CSV content in format: Symbol,Quantity,Price,Priority,BD User
func ParseOrdersCSV(csvContent string) ([]domain.CSVOrder, error) {
	lines := strings.Split(csvContent, "\n")
	var orders []domain.CSVOrder

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip header row
		if i == 0 {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "symbol") || strings.Contains(lower, "quantity") {
				continue
			}
		}

		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}

		qty, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			continue
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if err != nil {
			continue
		}

		priority, err := strconv.Atoi(strings.TrimSpace(fields[3]))
		if err != nil {
			priority = 1
		}

		order := domain.CSVOrder{
			Symbol:        strings.TrimSpace(fields[0]),
			Quantity:      qty,
			Price:         price,
			PriorityGroup: priority,
			BDUser:        strings.TrimSpace(fields[4]),
			Side:          "BID",
		}

		// Optional fields
		if len(fields) > 5 {
			order.Side = strings.TrimSpace(fields[5])
		}
		if len(fields) > 6 {
			order.Account = strings.TrimSpace(fields[6])
		}

		orders = append(orders, order)
	}

	return orders, nil
}
