CREATE TABLE IF NOT EXISTS sessions (
    id          SERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    vehicle_id  BIGINT      NOT NULL,
    slot_id     BIGINT      NOT NULL,
    start_time  TIMESTAMPTZ NOT NULL,
    end_time    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id   ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_active    ON sessions(id) WHERE end_time IS NULL;

CREATE TABLE IF NOT EXISTS payments (
    id          SERIAL PRIMARY KEY,
    session_id  BIGINT      NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    amount      NUMERIC     NOT NULL,
    paid_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_session_id ON payments(session_id);
