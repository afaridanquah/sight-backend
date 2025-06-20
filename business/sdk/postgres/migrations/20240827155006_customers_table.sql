-- +goose Up
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    middle_name VARCHAR(30),
    business_id UUID,
    user_id UUID,
    email TEXT,
    country VARCHAR(2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now ()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now ()),
    deleted_at TIMESTAMP
);

-- +goose Down
DROP TABLE customers;
