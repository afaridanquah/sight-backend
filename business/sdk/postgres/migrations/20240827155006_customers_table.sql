-- +goose Up
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    middle_name VARCHAR(30),
    org_id UUID,
    creator_id UUID,
    email TEXT,
    phone_number VARCHAR(20),
    city_of_birth VARCHAR(50),
    birth_country VARCHAR(50),
    date_of_birth DATE,
    identifications JSONB,
    addresses JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMP
);

-- +goose Down
DROP TABLE customers;
