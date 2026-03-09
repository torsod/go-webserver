package domain

import "time"

// Holiday entry that closes the market
type Holiday struct {
	Date       string `json:"date"`       // YYYY-MM-DD format
	Name       string `json:"name"`
	EarlyClose string `json:"earlyClose,omitempty"` // Optional HH:mm UTC
}

// AssetTypeSchedule per-asset-type schedule override
type AssetTypeSchedule struct {
	AssetType       string   `json:"assetType"`
	MarketOpen      string   `json:"marketOpen,omitempty"`
	MarketClose     string   `json:"marketClose,omitempty"`
	ExcludeHolidays []string `json:"excludeHolidays,omitempty"`
}

// SystemSettings system-wide configuration
type SystemSettings struct {
	ID string `json:"id" db:"id"`

	// Market Hours (UTC, HH:mm)
	MarketOpen  string `json:"marketOpen" db:"market_open"`
	MarketClose string `json:"marketClose" db:"market_close"`

	// Closing Windows defaults
	ClosingWindow1        int     `json:"closingWindow1" db:"closing_window1"`
	ClosingWindow2        int     `json:"closingWindow2" db:"closing_window2"`
	QuietDuration         int     `json:"quietDuration" db:"quiet_duration"`
	VolumeChangeThreshold float64 `json:"volumeChangeThreshold" db:"volume_change_threshold"`
	ClosingStartTime      string  `json:"closingStartTime" db:"closing_start_time"`

	// Seasoning Period defaults
	DefaultSeasoningPeriodMinutes    int `json:"defaultSeasoningPeriodMinutes" db:"default_seasoning_period_minutes"`
	DefaultModSeasoningPeriodMinutes int `json:"defaultModSeasoningPeriodMinutes" db:"default_mod_seasoning_period_minutes"`

	// Price Collar defaults
	DefaultPriceCollarPct float64 `json:"defaultPriceCollarPct" db:"default_price_collar_pct"`

	// Trading Calendar - stored as JSONB
	Holidays           []Holiday           `json:"holidays" db:"holidays"`
	AssetTypeSchedules []AssetTypeSchedule `json:"assetTypeSchedules" db:"asset_type_schedules"`

	// Allocation defaults
	DefaultAllocationMethod string `json:"defaultAllocationMethod" db:"default_allocation_method"`

	// Timestamps
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// DefaultSystemSettings returns factory defaults
func DefaultSystemSettings() SystemSettings {
	return SystemSettings{
		MarketOpen:                       "13:00",
		MarketClose:                      "21:00",
		ClosingWindow1:                   15,
		ClosingWindow2:                   15,
		QuietDuration:                    120,
		VolumeChangeThreshold:            1000000,
		ClosingStartTime:                 "20:30",
		DefaultSeasoningPeriodMinutes:    5,
		DefaultModSeasoningPeriodMinutes: 3,
		DefaultPriceCollarPct:            0.10,
		Holidays:                         []Holiday{},
		AssetTypeSchedules:               []AssetTypeSchedule{},
		DefaultAllocationMethod:          "PRICE_TIME",
		UpdatedAt:                        time.Now(),
	}
}
