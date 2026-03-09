package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// AllocationService handles allocation business logic
type AllocationService struct {
	orders    store.OrderStore
	offerings store.OfferingStore
	trades    store.TradeStore
	sessions  store.AllocationSessionStore
}

func NewAllocationService(
	orders store.OrderStore,
	offerings store.OfferingStore,
	trades store.TradeStore,
	sessions store.AllocationSessionStore,
) *AllocationService {
	return &AllocationService{
		orders: orders, offerings: offerings,
		trades: trades, sessions: sessions,
	}
}

// TempClose performs a temporary close (allocation)
func (s *AllocationService) TempClose(ctx context.Context, symbol string, offeringPrice float64, algorithm string, lmUser string) (*domain.AllocationSession, error) {
	offering, err := s.offerings.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("offering not found: %s", symbol)
	}

	allocationSize := domain.GetAllocationSize(offering)
	if allocationSize <= 0 {
		return nil, fmt.Errorf("offering has no allocation size")
	}

	// Get eligible orders (bids at or above offering price)
	orders, err := s.orders.FindForAllocation(ctx, symbol, offeringPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no eligible orders found")
	}

	// Run allocation
	var result *domain.AllocationResult
	switch algorithm {
	case "PRICE_TIME":
		result = allocatePriceTime(orders, allocationSize)
	case "PRO_RATA":
		result = allocateProRata(orders, allocationSize)
	case "PRIORITY_GROUP_PRO_RATA":
		result = allocatePriorityGroupProRata(orders, allocationSize)
	default:
		return nil, fmt.Errorf("unknown allocation algorithm: %s", algorithm)
	}

	// Create allocation session
	now := time.Now()
	session := &domain.AllocationSession{
		Symbol:           symbol,
		OfferingPrice:    offeringPrice,
		OfferingSize:     allocationSize,
		TotalAllocated:   result.TotalAllocated,
		AllocationMethod: algorithm,
		LMUser:           lmUser,
		Status:           domain.SessionStatusTemporary,
		CreatedAt:        now,
	}

	sessionID, err := s.sessions.Insert(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	session.ID = sessionID

	// Create trades
	// LM Sell trade
	_, err = s.trades.Insert(ctx, &domain.Trade{
		Symbol:       symbol,
		TradeType:    domain.TradeTypeSell,
		Quantity:     result.TotalAllocated,
		LeavesQty:    allocationSize - result.TotalAllocated,
		Price:        offeringPrice,
		UserID:       lmUser,
		AllocationID: sessionID,
		Status:       domain.TradeStatusFilled,
		Timestamp:    now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sell trade: %w", err)
	}

	// BD Buy trades
	for _, alloc := range result.Allocations {
		order, err := s.orders.FindByID(ctx, alloc.OrderID)
		if err != nil {
			continue
		}

		// Calculate commissions
		netAmount := float64(alloc.AllocatedQuantity) * offeringPrice
		grossSpread := netAmount * (offering.GrossUnderwritingSpread / 100)
		sellingConcession := grossSpread * (offering.SellingConcession / 100)

		isTracked := order.ExecInst == domain.ExecInstDTCTracking

		bidPrice := order.Price
		bidQty := order.Quantity
		pg := order.PriorityGroup

		_, err = s.trades.Insert(ctx, &domain.Trade{
			Symbol:                  symbol,
			TradeType:               domain.TradeTypeBuy,
			Quantity:                alloc.AllocatedQuantity,
			LeavesQty:              order.Quantity - alloc.AllocatedQuantity,
			Price:                   offeringPrice,
			BidPrice:                &bidPrice,
			BidQuantity:             &bidQty,
			PriorityGroup:           &pg,
			OrderType:               string(order.OrderType),
			Account:                 order.Account,
			ExecInst:                string(order.ExecInst),
			UserID:                  order.UserID,
			BDFirmID:                order.BDFirmID,
			OrderID:                 order.ID,
			AllocationID:            sessionID,
			IsDtcTracked:            &isTracked,
			SellingConcessionAmount: &sellingConcession,
			GrossSpreadAmount:       &grossSpread,
			Status:                  domain.TradeStatusFilled,
			Timestamp:               now,
		})
		if err != nil {
			continue
		}

		// Update order with allocation
		allocQty := alloc.AllocatedQuantity
		status := domain.OrderStatusFilled
		if alloc.IsPartial {
			status = domain.OrderStatusPartialFill
		}
		s.orders.UpdateFields(ctx, order.ID, map[string]interface{}{
			"allocated_quantity":    allocQty,
			"allocation_session_id": sessionID,
			"status":               string(status),
			"updated_at":           now,
		})
	}

	return session, nil
}

// Bust reverses an allocation session
func (s *AllocationService) Bust(ctx context.Context, sessionID string) error {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	if session.Status != domain.SessionStatusTemporary {
		return fmt.Errorf("can only bust TEMPORARY sessions, current: %s", session.Status)
	}

	now := time.Now()

	// Mark trades as busted
	trades, _ := s.trades.FindByAllocationID(ctx, sessionID)
	for _, t := range trades {
		s.trades.UpdateFields(ctx, t.ID, map[string]interface{}{
			"status":    string(domain.TradeStatusBusted),
			"busted_at": now,
		})
	}

	// Reset order allocations
	orders, _ := s.orders.FindAll(ctx)
	for _, o := range orders {
		if o.AllocationSessionID == sessionID {
			s.orders.UpdateFields(ctx, o.ID, map[string]interface{}{
				"allocated_quantity":    nil,
				"allocation_session_id": nil,
				"status":               string(domain.OrderStatusActive),
				"updated_at":           now,
			})
		}
	}

	// Update session
	return s.sessions.UpdateFields(ctx, sessionID, map[string]interface{}{
		"status":    string(domain.SessionStatusBusted),
		"busted_at": now,
	})
}

// Confirm confirms an allocation session
func (s *AllocationService) Confirm(ctx context.Context, sessionID string) error {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	if session.Status != domain.SessionStatusTemporary {
		return fmt.Errorf("can only confirm TEMPORARY sessions, current: %s", session.Status)
	}

	now := time.Now()
	return s.sessions.UpdateFields(ctx, sessionID, map[string]interface{}{
		"status":       string(domain.SessionStatusConfirmed),
		"confirmed_at": now,
	})
}

// CancelLeaves cancels unfilled order quantities
func (s *AllocationService) CancelLeaves(ctx context.Context, sessionID string) error {
	trades, err := s.trades.FindByAllocationID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to find trades: %w", err)
	}

	now := time.Now()
	for _, t := range trades {
		if t.TradeType == domain.TradeTypeBuy && t.LeavesQty > 0 {
			s.trades.UpdateFields(ctx, t.ID, map[string]interface{}{
				"leaves_qty":  0,
				"canceled_at": now,
			})
		}
	}
	return nil
}

// GetSessions returns allocation sessions
func (s *AllocationService) GetSessions(ctx context.Context, symbol string) ([]*domain.AllocationSession, error) {
	if symbol != "" {
		return s.sessions.FindBySymbol(ctx, symbol)
	}
	return s.sessions.FindAll(ctx)
}

// GetTrades returns trades for a session
func (s *AllocationService) GetTrades(ctx context.Context, sessionID string) ([]*domain.Trade, error) {
	return s.trades.FindByAllocationID(ctx, sessionID)
}

// Allocation Algorithms

func allocatePriceTime(orders []*domain.Order, allocationSize int64) *domain.AllocationResult {
	// Sort: price descending, then orderSequence ascending
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].Price != orders[j].Price {
			return orders[i].Price > orders[j].Price
		}
		return orders[i].OrderSequence < orders[j].OrderSequence
	})

	result := &domain.AllocationResult{}
	remaining := allocationSize

	for _, o := range orders {
		if remaining <= 0 {
			result.UnallocatedOrders++
			result.UnallocatedList = append(result.UnallocatedList, domain.UnallocatedOrder{
				OrderID: o.ID, Quantity: o.Quantity, Price: o.Price,
				PriorityGroup: o.PriorityGroup, UserID: o.UserID,
			})
			continue
		}

		qty := o.Quantity
		isPartial := false
		if qty > remaining {
			qty = remaining
			isPartial = true
		}

		result.Allocations = append(result.Allocations, domain.OrderAllocation{
			OrderID: o.ID, AllocatedQuantity: qty, IsPartial: isPartial,
		})
		result.TotalAllocated += qty
		remaining -= qty

		if isPartial {
			result.PartialAllocations++
		}
	}

	return result
}

func allocateProRata(orders []*domain.Order, allocationSize int64) *domain.AllocationResult {
	result := &domain.AllocationResult{}

	totalDemand := int64(0)
	for _, o := range orders {
		totalDemand += o.Quantity
	}

	if totalDemand <= allocationSize {
		// Everyone gets fully filled
		for _, o := range orders {
			result.Allocations = append(result.Allocations, domain.OrderAllocation{
				OrderID: o.ID, AllocatedQuantity: o.Quantity, IsPartial: false,
			})
			result.TotalAllocated += o.Quantity
		}
		return result
	}

	// Pro-rata allocation
	remaining := allocationSize
	for _, o := range orders {
		qty := (o.Quantity * allocationSize) / totalDemand
		if qty > remaining {
			qty = remaining
		}
		if qty <= 0 {
			result.UnallocatedOrders++
			result.UnallocatedList = append(result.UnallocatedList, domain.UnallocatedOrder{
				OrderID: o.ID, Quantity: o.Quantity, Price: o.Price,
				PriorityGroup: o.PriorityGroup, UserID: o.UserID,
			})
			continue
		}

		isPartial := qty < o.Quantity
		result.Allocations = append(result.Allocations, domain.OrderAllocation{
			OrderID: o.ID, AllocatedQuantity: qty, IsPartial: isPartial,
		})
		result.TotalAllocated += qty
		remaining -= qty

		if isPartial {
			result.PartialAllocations++
		}
	}

	// Distribute remainder (1 share each to highest priority/largest orders)
	if remaining > 0 {
		for i := range result.Allocations {
			if remaining <= 0 {
				break
			}
			result.Allocations[i].AllocatedQuantity++
			result.TotalAllocated++
			remaining--
		}
	}

	return result
}

func allocatePriorityGroupProRata(orders []*domain.Order, allocationSize int64) *domain.AllocationResult {
	// Group orders by priority group
	groups := make(map[int][]*domain.Order)
	for _, o := range orders {
		groups[o.PriorityGroup] = append(groups[o.PriorityGroup], o)
	}

	// Get sorted group numbers
	var groupNums []int
	for g := range groups {
		groupNums = append(groupNums, g)
	}
	sort.Ints(groupNums)

	result := &domain.AllocationResult{}
	remaining := allocationSize

	for _, g := range groupNums {
		groupOrders := groups[g]
		groupDemand := int64(0)
		for _, o := range groupOrders {
			groupDemand += o.Quantity
		}

		if groupDemand <= remaining {
			// Full fill for this group
			for _, o := range groupOrders {
				result.Allocations = append(result.Allocations, domain.OrderAllocation{
					OrderID: o.ID, AllocatedQuantity: o.Quantity, IsPartial: false,
				})
				result.TotalAllocated += o.Quantity
			}
			remaining -= groupDemand
		} else {
			// Pro-rata within group
			groupResult := allocateProRata(groupOrders, remaining)
			result.Allocations = append(result.Allocations, groupResult.Allocations...)
			result.TotalAllocated += groupResult.TotalAllocated
			result.PartialAllocations += groupResult.PartialAllocations
			result.UnallocatedOrders += groupResult.UnallocatedOrders
			result.UnallocatedList = append(result.UnallocatedList, groupResult.UnallocatedList...)
			remaining = 0
		}
	}

	return result
}
