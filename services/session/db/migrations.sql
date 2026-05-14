CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    vehicle_id TEXT NOT NULL,
    slot_id BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL,
    amount NUMERIC NOT NULL
);
