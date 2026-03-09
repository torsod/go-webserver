CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    market_open VARCHAR(10) NOT NULL DEFAULT '13:00',
    market_close VARCHAR(10) NOT NULL DEFAULT '21:00',
    closing_window1 INTEGER NOT NULL DEFAULT 15,
    closing_window2 INTEGER NOT NULL DEFAULT 15,
    quiet_duration INTEGER NOT NULL DEFAULT 120,
    volume_change_threshold NUMERIC(18,4) NOT NULL DEFAULT 1000000,
    closing_start_time VARCHAR(10) NOT NULL DEFAULT '20:30',
    default_seasoning_period_minutes INTEGER NOT NULL DEFAULT 5,
    default_mod_seasoning_period_minutes INTEGER NOT NULL DEFAULT 3,
    default_price_collar_pct NUMERIC(10,4) NOT NULL DEFAULT 0.10,
    holidays JSONB NOT NULL DEFAULT '[]',
    asset_type_schedules JSONB NOT NULL DEFAULT '[]',
    default_allocation_method VARCHAR(30) NOT NULL DEFAULT 'PRICE_TIME',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default settings
INSERT INTO settings (market_open, market_close, closing_window1, closing_window2,
    quiet_duration, volume_change_threshold, closing_start_time,
    default_seasoning_period_minutes, default_mod_seasoning_period_minutes,
    default_price_collar_pct, default_allocation_method)
VALUES ('13:00', '21:00', 15, 15, 120, 1000000, '20:30', 5, 3, 0.10, 'PRICE_TIME');
