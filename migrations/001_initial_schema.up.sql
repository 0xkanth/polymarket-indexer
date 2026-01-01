-- Polymarket Indexer - Initial Database Schema
-- TimescaleDB migration for event storage and aggregations

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- =============================================================================
-- EVENTS TABLE (Primary hypertable for all blockchain events)
-- =============================================================================

CREATE TABLE events (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    block_hash TEXT NOT NULL,
    tx_hash TEXT NOT NULL,
    tx_index INTEGER NOT NULL,
    log_index INTEGER NOT NULL,
    contract_address TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_signature TEXT NOT NULL,
    event_data JSONB NOT NULL,
    raw_data BYTEA,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Unique constraint on (tx_hash, log_index) for deduplication
    CONSTRAINT events_tx_log_unique UNIQUE (tx_hash, log_index)
);

-- Convert to hypertable partitioned by time
SELECT create_hypertable('events', 'time', 
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Create indexes for common queries
CREATE INDEX idx_events_block_number ON events (block_number DESC);
CREATE INDEX idx_events_contract ON events (contract_address, time DESC);
CREATE INDEX idx_events_event_name ON events (event_name, time DESC);
CREATE INDEX idx_events_tx_hash ON events (tx_hash);
CREATE INDEX idx_events_event_data ON events USING GIN (event_data);

-- =============================================================================
-- CTF EXCHANGE SPECIFIC TABLES
-- =============================================================================

-- Order filled events (trades)
CREATE TABLE order_fills (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    order_hash TEXT NOT NULL,
    maker TEXT NOT NULL,
    taker TEXT NOT NULL,
    maker_asset_id NUMERIC(78, 0) NOT NULL,
    taker_asset_id NUMERIC(78, 0) NOT NULL,
    maker_amount_filled NUMERIC(78, 0) NOT NULL,
    taker_amount_filled NUMERIC(78, 0) NOT NULL,
    fee NUMERIC(78, 0) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

SELECT create_hypertable('order_fills', 'time',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

CREATE INDEX idx_order_fills_maker ON order_fills (maker, time DESC);
CREATE INDEX idx_order_fills_taker ON order_fills (taker, time DESC);
CREATE INDEX idx_order_fills_order_hash ON order_fills (order_hash);
CREATE INDEX idx_order_fills_assets ON order_fills (maker_asset_id, taker_asset_id);

-- Token registrations (new markets)
CREATE TABLE token_registrations (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    token0 NUMERIC(78, 0) NOT NULL,
    token1 NUMERIC(78, 0) NOT NULL,
    condition_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT token_registrations_unique UNIQUE (token0, token1, condition_id)
);

SELECT create_hypertable('token_registrations', 'time',
    chunk_time_interval => INTERVAL '7 days',
    if_not_exists => TRUE
);

CREATE INDEX idx_token_registrations_condition ON token_registrations (condition_id);

-- =============================================================================
-- CONDITIONAL TOKENS TABLES
-- =============================================================================

-- Token transfers
CREATE TABLE token_transfers (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    log_index INTEGER NOT NULL,
    operator TEXT NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    token_id NUMERIC(78, 0) NOT NULL,
    amount NUMERIC(78, 0) NOT NULL,
    is_batch BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT token_transfers_unique UNIQUE (tx_hash, log_index)
);

SELECT create_hypertable('token_transfers', 'time',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

CREATE INDEX idx_token_transfers_from ON token_transfers (from_address, time DESC);
CREATE INDEX idx_token_transfers_to ON token_transfers (to_address, time DESC);
CREATE INDEX idx_token_transfers_token_id ON token_transfers (token_id, time DESC);

-- Condition preparations (market creation)
CREATE TABLE conditions (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    condition_id TEXT NOT NULL PRIMARY KEY,
    oracle TEXT NOT NULL,
    question_id TEXT NOT NULL,
    outcome_slot_count INTEGER NOT NULL,
    resolved BOOLEAN DEFAULT FALSE,
    resolution_block BIGINT,
    resolution_time TIMESTAMPTZ,
    payout_numerators INTEGER[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Not a hypertable - lookup table
CREATE INDEX idx_conditions_oracle ON conditions (oracle);
CREATE INDEX idx_conditions_question_id ON conditions (question_id);
CREATE INDEX idx_conditions_time ON conditions (time DESC);

-- =============================================================================
-- CHECKPOINT TABLE (for resume)
-- =============================================================================

CREATE TABLE checkpoints (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL UNIQUE,
    last_block BIGINT NOT NULL,
    last_block_hash TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Initialize checkpoint for indexer
INSERT INTO checkpoints (service_name, last_block, last_block_hash, updated_at)
VALUES ('indexer', 20558323, '0x0000000000000000000000000000000000000000000000000000000000000000', NOW())
ON CONFLICT (service_name) DO NOTHING;

-- =============================================================================
-- CONTINUOUS AGGREGATES (Materialized views for analytics)
-- =============================================================================

-- Hourly trading volume per market
CREATE MATERIALIZED VIEW trading_volume_hourly
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 hour', time) AS hour,
    maker_asset_id,
    taker_asset_id,
    COUNT(*) AS trade_count,
    SUM(maker_amount_filled) AS maker_volume,
    SUM(taker_amount_filled) AS taker_volume,
    SUM(fee) AS total_fees
FROM order_fills
GROUP BY hour, maker_asset_id, taker_asset_id
WITH NO DATA;

-- Refresh policy: update every hour, cover last 2 days
SELECT add_continuous_aggregate_policy('trading_volume_hourly',
    start_offset => INTERVAL '2 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

-- Daily active users (traders)
CREATE MATERIALIZED VIEW active_traders_daily
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 day', time) AS day,
    COUNT(DISTINCT maker) AS unique_makers,
    COUNT(DISTINCT taker) AS unique_takers,
    COUNT(DISTINCT maker) + COUNT(DISTINCT taker) AS unique_traders,
    COUNT(*) AS total_trades
FROM order_fills
GROUP BY day
WITH NO DATA;

SELECT add_continuous_aggregate_policy('active_traders_daily',
    start_offset => INTERVAL '7 days',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Market activity summary (per condition)
CREATE MATERIALIZED VIEW market_activity_daily
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 day', t.time) AS day,
    tr.condition_id,
    COUNT(DISTINCT t.tx_hash) AS transaction_count,
    COUNT(*) AS transfer_count,
    SUM(t.amount) AS volume
FROM token_transfers t
JOIN token_registrations tr ON (t.token_id = tr.token0 OR t.token_id = tr.token1)
GROUP BY day, tr.condition_id
WITH NO DATA;

SELECT add_continuous_aggregate_policy('market_activity_daily',
    start_offset => INTERVAL '7 days',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- =============================================================================
-- HELPER FUNCTIONS
-- =============================================================================

-- Function to get recent order fills for a market
CREATE OR REPLACE FUNCTION get_recent_fills(
    p_maker_asset NUMERIC,
    p_taker_asset NUMERIC,
    p_hours INTEGER DEFAULT 24
)
RETURNS TABLE (
    time TIMESTAMPTZ,
    maker TEXT,
    taker TEXT,
    maker_amount NUMERIC,
    taker_amount NUMERIC,
    price NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        of.time,
        of.maker,
        of.taker,
        of.maker_amount_filled,
        of.taker_amount_filled,
        CASE 
            WHEN of.maker_amount_filled > 0 
            THEN of.taker_amount_filled::NUMERIC / of.maker_amount_filled::NUMERIC
            ELSE 0
        END AS price
    FROM order_fills of
    WHERE of.maker_asset_id = p_maker_asset
      AND of.taker_asset_id = p_taker_asset
      AND of.time > NOW() - (p_hours || ' hours')::INTERVAL
    ORDER BY of.time DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get user trading stats
CREATE OR REPLACE FUNCTION get_user_stats(
    p_user_address TEXT,
    p_days INTEGER DEFAULT 30
)
RETURNS TABLE (
    total_trades BIGINT,
    total_volume NUMERIC,
    unique_markets BIGINT,
    first_trade TIMESTAMPTZ,
    last_trade TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*)::BIGINT AS total_trades,
        SUM(maker_amount_filled + taker_amount_filled) AS total_volume,
        COUNT(DISTINCT maker_asset_id || '-' || taker_asset_id)::BIGINT AS unique_markets,
        MIN(time) AS first_trade,
        MAX(time) AS last_trade
    FROM order_fills
    WHERE (maker = p_user_address OR taker = p_user_address)
      AND time > NOW() - (p_days || ' days')::INTERVAL;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- GRANTS (adjust for your users)
-- =============================================================================

-- Grant permissions to polymarket user
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO polymarket;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO polymarket;

-- =============================================================================
-- COMMENTS
-- =============================================================================

COMMENT ON TABLE events IS 'All blockchain events from Polymarket contracts';
COMMENT ON TABLE order_fills IS 'CTF Exchange order fill events (trades)';
COMMENT ON TABLE token_registrations IS 'New market token registrations';
COMMENT ON TABLE token_transfers IS 'Conditional token transfers (ERC-1155)';
COMMENT ON TABLE conditions IS 'Market conditions (questions) and their resolutions';
COMMENT ON TABLE checkpoints IS 'Block processing checkpoints for resume capability';

COMMENT ON MATERIALIZED VIEW trading_volume_hourly IS 'Hourly trading volume aggregates per market';
COMMENT ON MATERIALIZED VIEW active_traders_daily IS 'Daily active trader counts';
COMMENT ON MATERIALIZED VIEW market_activity_daily IS 'Daily market activity per condition';
