package service

import (
	"context"
	"fmt"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// OfferingService handles offering business logic
type OfferingService struct {
	store store.OfferingStore
}

func NewOfferingService(s store.OfferingStore) *OfferingService {
	return &OfferingService{store: s}
}

func (s *OfferingService) FindAll(ctx context.Context) ([]*domain.Offering, error) {
	return s.store.FindAll(ctx)
}

func (s *OfferingService) FindByID(ctx context.Context, id string) (*domain.Offering, error) {
	return s.store.FindByID(ctx, id)
}

func (s *OfferingService) FindBySymbol(ctx context.Context, symbol string) (*domain.Offering, error) {
	return s.store.FindBySymbol(ctx, symbol)
}

func (s *OfferingService) Create(ctx context.Context, offering *domain.Offering) (string, error) {
	if offering.Symbol == "" {
		return "", fmt.Errorf("symbol is required")
	}
	if offering.Market == "" {
		return "", fmt.Errorf("market is required")
	}
	if offering.HighPriceRange < offering.LowPriceRange {
		return "", fmt.Errorf("high price range must be >= low price range")
	}

	if offering.State == "" {
		offering.State = domain.OfferingStateNew
	}
	offering.CreatedAt = time.Now()

	if offering.TimeWindows == nil {
		offering.TimeWindows = []domain.TimeWindow{}
	}
	if offering.ExcludedBrokerDealers == nil {
		offering.ExcludedBrokerDealers = []string{}
	}
	if offering.ChangeLog == nil {
		offering.ChangeLog = []domain.OfferingChangeLog{}
	}

	return s.store.Insert(ctx, offering)
}

func (s *OfferingService) Update(ctx context.Context, id string, offering *domain.Offering) error {
	existing, err := s.store.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offering not found: %w", err)
	}

	// Preserve immutable fields
	offering.ID = existing.ID
	offering.CreatedAt = existing.CreatedAt
	now := time.Now()
	offering.UpdatedAt = &now

	return s.store.Update(ctx, id, offering)
}
