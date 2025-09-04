-- +goose Up
-- +goose StatementBegin
CREATE TYPE channel AS ENUM ('SMS', 'EMAIL');

CREATE TABLE otps (
    id VARCHAR(50),
    customer_id VARCHAR(50),
    hashed_code TEXT,
    code TEXT,
    channel channel,
    expires_at TIMESTAMP,
    verified_at TIMESTAMP,
    destination TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE otps;

DROP TYPE channel;

-- +goose StatementEnd
