CREATE TABLE offerings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    issuer VARCHAR(255) NOT NULL,
    market VARCHAR(50) NOT NULL,
    cusip VARCHAR(20),

    -- Asset classification
    asset_type VARCHAR(20) NOT NULL,
    security_type VARCHAR(20),

    -- Bond-specific
    yield_type VARCHAR(30),
    face_value NUMERIC(18,4),
    maturity_date TIMESTAMPTZ,
    offering_coupon NUMERIC(10,4),

    -- Fund/ETF
    nav NUMERIC(18,4),
    distribution_fee NUMERIC(10,4),

    -- State
    state VARCHAR(20) NOT NULL DEFAULT 'NEW',
    previous_state VARCHAR(20),

    -- Dates
    announcement_date TIMESTAMPTZ NOT NULL,
    bid_period_start_date TIMESTAMPTZ NOT NULL,
    scheduled_close_date TIMESTAMPTZ NOT NULL,
    settlement_date TIMESTAMPTZ,

    -- SEC Effectiveness
    sec_effectiveness_date TIMESTAMPTZ,
    sec_effectiveness_delay_minutes INTEGER NOT NULL DEFAULT 60,

    -- Price Collars
    max_price_allowed NUMERIC(18,4) NOT NULL,
    min_price_allowed NUMERIC(18,4) NOT NULL,

    -- Price Range
    high_price_range NUMERIC(18,4) NOT NULL,
    low_price_range NUMERIC(18,4) NOT NULL,

    -- Minimum Acceptable IPO Price
    min_acceptable_ipo_price NUMERIC(18,4) NOT NULL,

    -- Dividend
    dividend NUMERIC(18,4),

    -- Offer Quantities
    primary_quantity BIGINT NOT NULL DEFAULT 0,
    upsize_quantity BIGINT NOT NULL DEFAULT 0,
    committed_secondary_quantity BIGINT NOT NULL DEFAULT 0,
    lm_short_quantity BIGINT NOT NULL DEFAULT 0,

    -- Order Constraints
    min_bid_quantity BIGINT NOT NULL DEFAULT 10,
    max_bid_quantity BIGINT NOT NULL DEFAULT 1000000,
    qty_increment BIGINT NOT NULL DEFAULT 1,
    price_increment NUMERIC(10,4) NOT NULL DEFAULT 0.01,

    -- MinQty Configuration
    min_order_size_for_min_qty BIGINT NOT NULL DEFAULT 100000,
    max_min_qty_percentage NUMERIC(10,4) NOT NULL DEFAULT 50,
    min_qty_deadline TIMESTAMPTZ,

    -- Allocation
    allocation_method VARCHAR(30) NOT NULL DEFAULT 'PRO_RATA',
    min_allocation_per_account BIGINT NOT NULL DEFAULT 0,

    -- Preferential Bids
    preferential_bids_allowed VARCHAR(30) NOT NULL DEFAULT 'NO',

    -- Secondary Sell Offers
    secondary_offers_allowed VARCHAR(30) NOT NULL DEFAULT 'NO',

    -- Cprice publishing
    cprice_publish_mode VARCHAR(30) NOT NULL DEFAULT 'MIN_IPO_PRICE',
    cprice_publishing_active BOOLEAN NOT NULL DEFAULT false,

    -- Fees
    gross_underwriting_spread NUMERIC(10,4) NOT NULL DEFAULT 5,
    selling_concession NUMERIC(10,4) NOT NULL DEFAULT 40,

    -- Priority Group Test Mode
    prio_group_test BOOLEAN NOT NULL DEFAULT true,

    -- JSONB fields
    time_windows JSONB NOT NULL DEFAULT '[]',
    closing_windows_config JSONB NOT NULL DEFAULT '{}',
    closing_windows_data JSONB,
    listing_minimums JSONB,
    excluded_broker_dealers JSONB NOT NULL DEFAULT '[]',

    -- Lead Manager BD
    lead_manager_bd_id VARCHAR(100),

    -- Clearing/Offering Price
    clearing_price NUMERIC(18,4),
    offering_price NUMERIC(18,4),
    bond_offering_price NUMERIC(18,4),

    -- Computed metrics
    indicative_clearing_price NUMERIC(18,4),
    cprice NUMERIC(18,4),
    total_demand BIGINT,
    total_orders BIGINT,

    -- Change Log
    change_log JSONB NOT NULL DEFAULT '[]',

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_offerings_state ON offerings(state);
CREATE INDEX idx_offerings_symbol ON offerings(symbol);
