CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    is_logged_in BOOLEAN NOT NULL DEFAULT false,
    disabled BOOLEAN NOT NULL DEFAULT false,
    bd_role VARCHAR(20),
    bd_firm_id VARCHAR(50),
    bd_firm_name VARCHAR(100),
    qsr_active BOOLEAN NOT NULL DEFAULT false,
    firm_accounts JSONB NOT NULL DEFAULT '[]',
    assigned_accounts JSONB NOT NULL DEFAULT '[]',
    read_only_accounts JSONB NOT NULL DEFAULT '[]',
    can_cancel_group_orders BOOLEAN NOT NULL DEFAULT false,
    can_write_group_orders BOOLEAN NOT NULL DEFAULT false,
    lm_firm_id VARCHAR(50),
    email VARCHAR(255),
    display_name VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_users_user_type ON users(user_type);
CREATE INDEX idx_users_bd_firm_id ON users(bd_firm_id);
