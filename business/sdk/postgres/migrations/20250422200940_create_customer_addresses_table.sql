-- +goose Up
-- +goose StatementBegin
CREATE TABLE customer_addresses (
    id UUID,
    customer_id UUID PRIMARY KEY NOT NULL,
    address_1 TEXT,
    address_2 TEXT,
    city TEXT,
    state VARCHAR(50),
    zip VARCHAR(25),
    country VARCHAR(2),
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW (),
    deleted_at TIMESTAMP DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE customer_addresses;

-- +goose StatementEnd
