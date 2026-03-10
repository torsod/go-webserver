package store

import (
	"context"

	"github.com/torsod/go-webserver/internal/domain"
)

// OfferingStore defines the data access interface for offerings
type OfferingStore interface {
	FindAll(ctx context.Context) ([]*domain.Offering, error)
	FindByID(ctx context.Context, id string) (*domain.Offering, error)
	FindBySymbol(ctx context.Context, symbol string) (*domain.Offering, error)
	FindByStates(ctx context.Context, states []domain.OfferingState) ([]*domain.Offering, error)
	Insert(ctx context.Context, offering *domain.Offering) (string, error)
	Update(ctx context.Context, id string, offering *domain.Offering) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// OrderStore defines the data access interface for orders
type OrderStore interface {
	FindAll(ctx context.Context) ([]*domain.Order, error)
	FindBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error)
	FindActive(ctx context.Context) ([]*domain.Order, error)
	FindActiveBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error)
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	FindForAllocation(ctx context.Context, symbol string, minPrice float64) ([]*domain.Order, error)
	Insert(ctx context.Context, order *domain.Order) (string, error)
	Update(ctx context.Context, id string, order *domain.Order) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	NextSequence(ctx context.Context) (int64, error)
}

// UserStore defines the data access interface for users
type UserStore interface {
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	Insert(ctx context.Context, user *domain.User) (string, error)
	Update(ctx context.Context, id string, user *domain.User) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// TradeStore defines the data access interface for trades
type TradeStore interface {
	FindBySymbol(ctx context.Context, symbol string) ([]*domain.Trade, error)
	FindByAllocationID(ctx context.Context, allocationID string) ([]*domain.Trade, error)
	Insert(ctx context.Context, trade *domain.Trade) (string, error)
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// AllocationSessionStore defines the data access interface for allocation sessions
type AllocationSessionStore interface {
	FindAll(ctx context.Context) ([]*domain.AllocationSession, error)
	FindByID(ctx context.Context, id string) (*domain.AllocationSession, error)
	FindBySymbol(ctx context.Context, symbol string) ([]*domain.AllocationSession, error)
	Insert(ctx context.Context, session *domain.AllocationSession) (string, error)
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// SettingsStore defines the data access interface for system settings
type SettingsStore interface {
	Get(ctx context.Context) (*domain.SystemSettings, error)
	Update(ctx context.Context, settings *domain.SystemSettings) error
}

// FIXSessionStore defines the data access interface for FIX sessions
type FIXSessionStore interface {
	FindByID(ctx context.Context, id string) (*domain.FIXSession, error)
	FindActive(ctx context.Context) (*domain.FIXSession, error)
	Insert(ctx context.Context, session *domain.FIXSession) (string, error)
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// FIXOrderStore defines the data access interface for FIX orders
type FIXOrderStore interface {
	Insert(ctx context.Context, order *domain.FIXOrder) (string, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]*domain.FIXOrder, error)
	FindByClOrdID(ctx context.Context, clOrdID string) (*domain.FIXOrder, error)
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
}

// FIXLogStore defines the data access interface for FIX logs
type FIXLogStore interface {
	Insert(ctx context.Context, log *domain.FIXLog) (string, error)
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
	FindBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.FIXLog, error)
	DeleteBySessionID(ctx context.Context, sessionID string) (int64, error)
}
