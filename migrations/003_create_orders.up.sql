CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    quantity BIGINT NOT NULL,
    price NUMERIC(18,4) NOT NULL,

    -- MinQty
    min_qty BIGINT,

    -- Account & Tracking
    account VARCHAR(100),
    exec_inst VARCHAR(5) NOT NULL DEFAULT '',

    -- Priority & Time
    priority_group INTEGER NOT NULL DEFAULT 1,
    "timestamp" TIMESTAMPTZ NOT NULL,
    original_entry_time TIMESTAMPTZ NOT NULL,
    order_sequence BIGINT NOT NULL,

    -- User/Ownership
    user_id VARCHAR(100) NOT NULL,
    bd_firm_id VARCHAR(50),
    entered_by VARCHAR(100),

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Seasoning Period tracking
    seasoning_expires_at TIMESTAMPTZ,
    time_window_at_entry INTEGER NOT NULL DEFAULT 1,

    -- Halted Order tracking
    halt_reason TEXT,
    halted_at TIMESTAMPTZ,
    pre_halt_timestamp TIMESTAMPTZ,

    -- Allocation results
    allocated_quantity BIGINT,
    allocation_session_id UUID,

    -- Cancellation
    canceled_at TIMESTAMPTZ,
    cancel_reason TEXT,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    modification_history JSONB NOT NULL DEFAULT '[]'
);

CREATE INDEX idx_orders_symbol_status ON orders(symbol, status);
CREATE INDEX idx_orders_symbol_side_price ON orders(symbol, side, status, price);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_bd_firm_id ON orders(bd_firm_id);
CREATE INDEX idx_orders_order_sequence ON orders(order_sequence);
CREATE INDEX idx_orders_allocation_session ON orders(allocation_session_id);

-- Sequence for order sequence counter
CREATE SEQUENCE order_sequence_seq START 1;
