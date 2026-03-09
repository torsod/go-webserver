package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

type settingsStore struct {
	pool *pgxpool.Pool
}

func NewSettingsStore(pool *pgxpool.Pool) SettingsStore {
	return &settingsStore{pool: pool}
}

func (s *settingsStore) Get(ctx context.Context) (*domain.SystemSettings, error) {
	settings := &domain.SystemSettings{}
	var holJSON, atsJSON []byte

	err := s.pool.QueryRow(ctx, `
		SELECT id, market_open, market_close, closing_window1, closing_window2,
			quiet_duration, volume_change_threshold, closing_start_time,
			default_seasoning_period_minutes, default_mod_seasoning_period_minutes,
			default_price_collar_pct, holidays, asset_type_schedules,
			default_allocation_method, updated_at
		FROM settings LIMIT 1
	`).Scan(
		&settings.ID, &settings.MarketOpen, &settings.MarketClose,
		&settings.ClosingWindow1, &settings.ClosingWindow2,
		&settings.QuietDuration, &settings.VolumeChangeThreshold, &settings.ClosingStartTime,
		&settings.DefaultSeasoningPeriodMinutes, &settings.DefaultModSeasoningPeriodMinutes,
		&settings.DefaultPriceCollarPct, &holJSON, &atsJSON,
		&settings.DefaultAllocationMethod, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	if holJSON != nil {
		json.Unmarshal(holJSON, &settings.Holidays)
	}
	if atsJSON != nil {
		json.Unmarshal(atsJSON, &settings.AssetTypeSchedules)
	}

	return settings, nil
}

func (s *settingsStore) Update(ctx context.Context, settings *domain.SystemSettings) error {
	holJSON, _ := json.Marshal(settings.Holidays)
	atsJSON, _ := json.Marshal(settings.AssetTypeSchedules)

	_, err := s.pool.Exec(ctx, `
		UPDATE settings SET
			market_open=$1, market_close=$2,
			closing_window1=$3, closing_window2=$4,
			quiet_duration=$5, volume_change_threshold=$6, closing_start_time=$7,
			default_seasoning_period_minutes=$8, default_mod_seasoning_period_minutes=$9,
			default_price_collar_pct=$10, holidays=$11, asset_type_schedules=$12,
			default_allocation_method=$13, updated_at=NOW()
		WHERE id=$14`,
		settings.MarketOpen, settings.MarketClose,
		settings.ClosingWindow1, settings.ClosingWindow2,
		settings.QuietDuration, settings.VolumeChangeThreshold, settings.ClosingStartTime,
		settings.DefaultSeasoningPeriodMinutes, settings.DefaultModSeasoningPeriodMinutes,
		settings.DefaultPriceCollarPct, holJSON, atsJSON,
		settings.DefaultAllocationMethod, settings.ID,
	)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}
