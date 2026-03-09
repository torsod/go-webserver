package service

import (
	"context"
	"fmt"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// OrderService handles order business logic
type OrderService struct {
	orders    store.OrderStore
	offerings store.OfferingStore
}

func NewOrderService(orders store.OrderStore, offerings store.OfferingStore) *OrderService {
	return &OrderService{orders: orders, offerings: offerings}
}

func (s *OrderService) FindAll(ctx context.Context) ([]*domain.Order, error) {
	return s.orders.FindAll(ctx)
}

func (s *OrderService) FindBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error) {
	return s.orders.FindBySymbol(ctx, symbol)
}

func (s *OrderService) FindActiveBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error) {
	return s.orders.FindActiveBySymbol(ctx, symbol)
}

// Insert creates a new order with full validation
func (s *OrderService) Insert(ctx context.Context,
	symbol *string, side domain.OrderSide, orderType domain.OrderType,
	quantity int64, price float64, minQty *int64, account string,
	execInst domain.ExecInst, priorityGroup int, userID string,
) (*domain.Order, error) {

	if symbol == nil || *symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	// Find offering
	offering, err := s.offerings.FindBySymbol(ctx, *symbol)
	if err != nil {
		return nil, fmt.Errorf("offering not found: %s", *symbol)
	}

	// Check offering state
	if !IsOrderEntryAllowed(offering.State) {
		return nil, fmt.Errorf("order entry not allowed in state %s", offering.State)
	}

	// Validate quantity
	if quantity < offering.MinBidQuantity {
		return nil, fmt.Errorf("quantity %d below minimum %d", quantity, offering.MinBidQuantity)
	}
	if quantity > offering.MaxBidQuantity {
		return nil, fmt.Errorf("quantity %d exceeds maximum %d", quantity, offering.MaxBidQuantity)
	}
	if offering.QtyIncrement > 0 && quantity%offering.QtyIncrement != 0 {
		return nil, fmt.Errorf("quantity must be a multiple of %d", offering.QtyIncrement)
	}

	// Validate price
	if price < offering.MinPriceAllowed {
		return nil, fmt.Errorf("price %.4f below minimum %.4f", price, offering.MinPriceAllowed)
	}
	if price > offering.MaxPriceAllowed {
		return nil, fmt.Errorf("price %.4f exceeds maximum %.4f", price, offering.MaxPriceAllowed)
	}

	// Get next order sequence
	seq, err := s.orders.NextSequence(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order sequence: %w", err)
	}

	now := time.Now()

	// Determine time window and seasoning
	timeWindowAtEntry := 1
	var seasoningExpiresAt *time.Time
	for _, tw := range offering.TimeWindows {
		if now.After(tw.StartTime) || now.Equal(tw.StartTime) {
			timeWindowAtEntry = tw.WindowNumber
			if tw.SeasoningPeriodMinutes > 0 {
				expires := now.Add(time.Duration(tw.SeasoningPeriodMinutes) * time.Minute)
				seasoningExpiresAt = &expires
			}
			if !offering.PrioGroupTest {
				priorityGroup = tw.PriorityGroup
			}
		}
	}

	order := &domain.Order{
		Symbol:             *symbol,
		Side:               side,
		OrderType:          orderType,
		Quantity:           quantity,
		Price:              price,
		MinQty:             minQty,
		Account:            account,
		ExecInst:           execInst,
		PriorityGroup:      priorityGroup,
		Timestamp:          now,
		OriginalEntryTime:  now,
		OrderSequence:      seq,
		UserID:             userID,
		Status:             domain.OrderStatusActive,
		SeasoningExpiresAt: seasoningExpiresAt,
		TimeWindowAtEntry:  timeWindowAtEntry,
		CreatedAt:          now,
	}

	id, err := s.orders.Insert(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}
	order.ID = id

	return order, nil
}

// Update modifies an existing order
func (s *OrderService) Update(ctx context.Context, id string, quantity *int64, price *float64, minQty *int64, userID string) error {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status != domain.OrderStatusActive {
		return fmt.Errorf("can only modify ACTIVE orders, current status: %s", order.Status)
	}

	// Find offering for state check
	offering, err := s.offerings.FindBySymbol(ctx, order.Symbol)
	if err != nil {
		return fmt.Errorf("offering not found: %w", err)
	}

	isPriceUp := false
	if price != nil && *price > order.Price {
		isPriceUp = true
	}

	if !IsOrderModificationAllowed(offering.State, isPriceUp) {
		return fmt.Errorf("order modification not allowed in state %s", offering.State)
	}

	now := time.Now()
	var mods []domain.OrderModification

	if quantity != nil && *quantity != order.Quantity {
		mods = append(mods, domain.OrderModification{
			Field: "quantity", OldValue: order.Quantity, NewValue: *quantity,
			ModifiedBy: userID, ModifiedAt: now, ImpactsTimePriority: true,
		})
		order.Quantity = *quantity
	}

	if price != nil && *price != order.Price {
		mods = append(mods, domain.OrderModification{
			Field: "price", OldValue: order.Price, NewValue: *price,
			ModifiedBy: userID, ModifiedAt: now, ImpactsTimePriority: true,
		})
		order.Price = *price
	}

	if minQty != nil {
		oldMinQty := int64(0)
		if order.MinQty != nil {
			oldMinQty = *order.MinQty
		}
		if *minQty != oldMinQty {
			mods = append(mods, domain.OrderModification{
				Field: "minQty", OldValue: oldMinQty, NewValue: *minQty,
				ModifiedBy: userID, ModifiedAt: now, ImpactsTimePriority: true,
			})
			order.MinQty = minQty
		}
	}

	if len(mods) > 0 {
		// Reset time priority on significant modifications
		order.Timestamp = now
		order.ModificationHistory = append(order.ModificationHistory, mods...)
		order.UpdatedAt = &now

		// Recalculate seasoning
		for _, tw := range offering.TimeWindows {
			if now.After(tw.StartTime) || now.Equal(tw.StartTime) {
				if tw.SeasoningPeriodMinutes > 0 {
					expires := now.Add(time.Duration(tw.SeasoningPeriodMinutes) * time.Minute)
					order.SeasoningExpiresAt = &expires
				}
			}
		}
	}

	return s.orders.Update(ctx, id, order)
}

// Cancel cancels an order
func (s *OrderService) Cancel(ctx context.Context, id string, userID string, reason string) error {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status == domain.OrderStatusCanceled {
		return fmt.Errorf("order already canceled")
	}

	now := time.Now()
	order.Status = domain.OrderStatusCanceled
	order.CanceledAt = &now
	order.CancelReason = reason
	order.UpdatedAt = &now

	return s.orders.Update(ctx, id, order)
}

// CancelAll cancels all active orders for a user/symbol
func (s *OrderService) CancelAll(ctx context.Context, userID string, symbol string) (int, error) {
	orders, err := s.orders.FindActiveBySymbol(ctx, symbol)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, o := range orders {
		if userID != "" && o.UserID != userID {
			continue
		}
		if err := s.Cancel(ctx, o.ID, userID, "bulk cancel"); err == nil {
			count++
		}
	}
	return count, nil
}
