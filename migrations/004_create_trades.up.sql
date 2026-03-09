CREATE TABLE allocation_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    offering_price NUMERIC(18,4) NOT NULL,
    offering_size BIGINT NOT NULL,
    total_allocated BIGINT NOT NULL DEFAULT 0,
    allocation_method VARCHAR(30) NOT NULL,
    lm_user VARCHAR(100) NOT NULL,

    -- Preferential allocation summary
    preferential_allocated BIGINT,
    preferential_bidders INTEGER,

    -- Priority group summaries
    group_summaries JSONB,

    -- MinQty impact
    min_qty_excluded INTEGER,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'TEMPORARY',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    busted_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ
);

CREATE INDEX idx_allocation_sessions_symbol ON allocation_sessions(symbol);
CREATE INDEX idx_allocation_sessions_status ON allocation_sessions(status);

CREATE TABLE trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    trade_type VARCHAR(10) NOT NULL,
    quantity BIGINT NOT NULL,
    leaves_qty BIGINT NOT NULL DEFAULT 0,
    price NUMERIC(18,4) NOT NULL,

    -- Original order reference
    bid_price NUMERIC(18,4),
    bid_quantity BIGINT,
    priority_group INTEGER,
    order_type VARCHAR(20),
    account VARCHAR(100),
    exec_inst VARCHAR(5),

    -- User/Ownership
    user_id VARCHAR(100) NOT NULL,
    bd_firm_id VARCHAR(50),
    order_id UUID,
    allocation_id UUID NOT NULL REFERENCES allocation_sessions(id),

    -- Settlement
    is_dtc_tracked BOOLEAN,
    selling_concession_amount NUMERIC(18,4),
    gross_spread_amount NUMERIC(18,4),

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'NEW',
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    busted_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ
);

CREATE INDEX idx_trades_symbol ON trades(symbol);
CREATE INDEX idx_trades_allocation_id ON trades(allocation_id);
CREATE INDEX idx_trades_user_id ON trades(user_id);
