package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// OrderbookService handles orderbook calculations
type OrderbookService struct {
	orders    store.OrderStore
	offerings store.OfferingStore
}

func NewOrderbookService(orders store.OrderStore, offerings store.OfferingStore) *OrderbookService {
	return &OrderbookService{orders: orders, offerings: offerings}
}

// GetSummary returns the full orderbook summary for a symbol
func (s *OrderbookService) GetSummary(ctx context.Context, symbol string) (*domain.OrderbookSummary, error) {
	offering, err := s.offerings.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("offering not found: %s", symbol)
	}

	orders, err := s.orders.FindActiveBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}

	summary := &domain.OrderbookSummary{
		Symbol: symbol,
	}

	// Separate bids and offers
	var bids, offers []*domain.Order
	for _, o := range orders {
		if o.Side == domain.OrderSideBid {
			bids = append(bids, o)
		} else {
			offers = append(offers, o)
		}
	}

	// Build bid entries (price descending)
	summary.BidEntries = buildEntries(bids, true)
	summary.OfferEntries = buildEntries(offers, false)

	// Calculate totals
	for _, e := range summary.BidEntries {
		summary.TotalBidQuantity += e.AccumulatedSize
		summary.TotalBidOrders += e.OrderCount
	}
	for _, e := range summary.OfferEntries {
		summary.TotalOfferQuantity += e.AccumulatedSize
		summary.TotalOfferOrders += e.OrderCount
	}

	// Calculate clearing price
	offeringSize := domain.GetOfferingSize(offering)
	clearingPrice := calculateClearingPrice(summary.BidEntries, offeringSize)
	summary.ClearingPrice = clearingPrice

	// Calculate Cprice
	if clearingPrice != nil {
		cprice := *clearingPrice
		if cprice < offering.MinAcceptableIPOPrice {
			cprice = offering.MinAcceptableIPOPrice
		}
		summary.Cprice = &cprice

		// Calculate aVWAP
		avwap := calculateAVWAP(summary.BidEntries, offeringSize)
		summary.AVWAP = avwap
	}

	// Calculate rVWAP (all bids)
	rvwap := calculateRVWAP(summary.BidEntries)
	summary.RVWAP = rvwap

	// Subscription ratio
	if offeringSize > 0 {
		ratio := float64(summary.TotalBidQuantity) / float64(offeringSize)
		summary.SubscriptionRatio = &ratio
	}

	return summary, nil
}

// GetCprice returns just the clearing price info
func (s *OrderbookService) GetCprice(ctx context.Context, symbol string) (map[string]interface{}, error) {
	offering, err := s.offerings.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("offering not found: %s", symbol)
	}

	// Check publish mode
	switch offering.CpricePublishMode {
	case domain.CpricePublishModeNotPublishedManualOn:
		if !offering.CpricePublishingActive {
			return map[string]interface{}{
				"symbol":  symbol,
				"cprice":  nil,
				"mode":    offering.CpricePublishMode,
				"message": "Cprice publishing not active",
			}, nil
		}
	case domain.CpricePublishModeNotPublishedAutoOn:
		if !offering.CpricePublishingActive {
			return map[string]interface{}{
				"symbol":  symbol,
				"cprice":  nil,
				"mode":    offering.CpricePublishMode,
				"message": "Cprice publishing not yet activated",
			}, nil
		}
	}

	summary, err := s.GetSummary(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"symbol":        symbol,
		"cprice":        summary.Cprice,
		"clearingPrice": summary.ClearingPrice,
		"mode":          offering.CpricePublishMode,
	}, nil
}

// buildEntries aggregates orders by price level
func buildEntries(orders []*domain.Order, descending bool) []domain.OrderbookEntry {
	// Group by price
	priceMap := make(map[float64]*domain.OrderbookEntry)
	for _, o := range orders {
		entry, exists := priceMap[o.Price]
		if !exists {
			entry = &domain.OrderbookEntry{Price: o.Price}
			priceMap[o.Price] = entry
		}
		entry.AccumulatedSize += o.Quantity
		entry.DollarValue += float64(o.Quantity) * o.Price
		entry.OrderCount++
	}

	// Convert to slice
	var entries []domain.OrderbookEntry
	for _, e := range priceMap {
		entries = append(entries, *e)
	}

	// Sort
	sort.Slice(entries, func(i, j int) bool {
		if descending {
			return entries[i].Price > entries[j].Price
		}
		return entries[i].Price < entries[j].Price
	})

	// Calculate running totals
	var totalSize int64
	var totalDollar float64
	for i := range entries {
		totalSize += entries[i].AccumulatedSize
		totalDollar += entries[i].DollarValue
		entries[i].TotalAccumulatedSize = totalSize
		entries[i].TotalDollarValue = totalDollar
	}

	return entries
}

// calculateClearingPrice finds the highest bid price where accumulated demand meets supply
func calculateClearingPrice(bidEntries []domain.OrderbookEntry, offeringSize int64) *float64 {
	if len(bidEntries) == 0 || offeringSize <= 0 {
		return nil
	}

	// Walk from highest to lowest price, finding where accumulated size >= offering size
	for i := range bidEntries {
		if bidEntries[i].TotalAccumulatedSize >= offeringSize {
			return &bidEntries[i].Price
		}
	}

	return nil
}

// calculateAVWAP computes accumulated VWAP up to the offering size
func calculateAVWAP(bidEntries []domain.OrderbookEntry, offeringSize int64) *float64 {
	if len(bidEntries) == 0 || offeringSize <= 0 {
		return nil
	}

	var totalValue float64
	var totalQty int64

	for _, e := range bidEntries {
		remaining := offeringSize - totalQty
		if remaining <= 0 {
			break
		}
		qty := e.AccumulatedSize
		if qty > remaining {
			qty = remaining
		}
		totalValue += float64(qty) * e.Price
		totalQty += qty
	}

	if totalQty == 0 {
		return nil
	}

	avwap := totalValue / float64(totalQty)
	return &avwap
}

// calculateRVWAP computes running VWAP of all bids
func calculateRVWAP(bidEntries []domain.OrderbookEntry) *float64 {
	if len(bidEntries) == 0 {
		return nil
	}

	var totalValue float64
	var totalQty int64

	for _, e := range bidEntries {
		totalValue += e.DollarValue
		totalQty += e.AccumulatedSize
	}

	if totalQty == 0 {
		return nil
	}

	rvwap := totalValue / float64(totalQty)
	return &rvwap
}
