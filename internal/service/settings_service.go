package service

import (
	"context"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// SettingsService manages system settings
type SettingsService struct {
	settings store.SettingsStore
}

func NewSettingsService(settings store.SettingsStore) *SettingsService {
	return &SettingsService{settings: settings}
}

func (s *SettingsService) Get(ctx context.Context) (*domain.SystemSettings, error) {
	return s.settings.Get(ctx)
}

func (s *SettingsService) Update(ctx context.Context, settings *domain.SystemSettings) error {
	return s.settings.Update(ctx, settings)
}
