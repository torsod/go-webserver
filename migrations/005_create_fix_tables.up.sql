CREATE TABLE fix_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(100) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    sender_comp_id VARCHAR(50) NOT NULL,
    target_comp_id VARCHAR(50) NOT NULL,
    connected BOOLEAN NOT NULL DEFAULT false,
    simulated BOOLEAN NOT NULL DEFAULT false,
    messages_sent INTEGER NOT NULL DEFAULT 0,
    messages_received INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE fix_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cl_ord_id VARCHAR(100) NOT NULL,
    session_id VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(5) NOT NULL,
    quantity BIGINT NOT NULL,
    price NUMERIC(18,4) NOT NULL,
    ord_type VARCHAR(5) NOT NULL DEFAULT '2',
    time_in_force VARCHAR(5) NOT NULL DEFAULT '1',
    account VARCHAR(100),
    main_order_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'NEW',
    transact_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fix_orders_session_id ON fix_orders(session_id);
CREATE INDEX idx_fix_orders_main_order_id ON fix_orders(main_order_id);

CREATE TABLE fix_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(100) NOT NULL,
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    level VARCHAR(20) NOT NULL DEFAULT 'INFO',
    message TEXT NOT NULL,
    direction VARCHAR(5),
    raw_data TEXT
);

CREATE INDEX idx_fix_logs_session_id ON fix_logs(session_id);
CREATE INDEX idx_fix_logs_timestamp ON fix_logs("timestamp");
