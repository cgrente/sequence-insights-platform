-- Create sequences table.
-- This file is mounted into Postgres' /docker-entrypoint-initdb.d for local dev.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    values_json JSONB NOT NULL,

    count INT NOT NULL,
    sum_fourth_powers_non_positive BIGINT NOT NULL,
    min_value BIGINT NOT NULL,
    max_value BIGINT NOT NULL,

    processed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_sequences_created_at ON sequences (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sequences_processed ON sequences (processed);
