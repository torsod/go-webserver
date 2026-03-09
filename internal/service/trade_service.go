package service

import (
	"context"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// TradeService handles trade queries
type TradeService struct {
	trades store.TradeStore
}

func NewTradeService(trades store.TradeStore) *TradeService {
	return &TradeService{trades: trades}
}

func (s *TradeService) FindBySymbol(ctx context.Context, symbol string) ([]*domain.Trade, error) {
	return s.trades.FindBySymbol(ctx, symbol)
}

func (s *TradeService) FindByAllocationID(ctx context.Context, allocationID string) ([]*domain.Trade, error) {
	return s.trades.FindByAllocationID(ctx, allocationID)
}
