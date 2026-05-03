CREATE TABLE IF NOT EXISTS parking_slots (
    id SERIAL PRIMARY KEY,
    is_occupied BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_slots_free ON parking_slots(is_occupied);

INSERT INTO parking_slots (is_occupied)
SELECT false FROM generate_series(1, 50)
WHERE NOT EXISTS (SELECT 1 FROM parking_slots);