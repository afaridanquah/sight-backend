-- +goose Up
-- +goose StatementBegin
CREATE TYPE identification_type AS ENUM (
    'PASSPORT',
    'DRIVERS_LICENSE',
    'NATIONAL_ID',
    'RESIDENT_PERMIT',
    'SSN'
);

CREATE TABLE customer_identifications (
    id CHAR(31) PRIMARY KEY NOT NULL,
    customer_id CHAR(31),
    first_name VARCHAR(45),
    last_name VARCHAR(45),
    middle_name VARCHAR(45),
    other_names VARCHAR(45),
    pin VARCHAR(100) NOT NULL,
    identification_type identification_type NOT NULL,
    issued_country VARCHAR(2),
    issued_date DATE,
    place_of_birth VARCHAR(45),
    place_of_issue VARCHAR(45),
    date_of_birth DATE,
    address_1 TEXT,
    address_2 TEXT,
    city VARCHAR(150),
    state_region VARCHAR(150),
    zip_code VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX customer_identifications_customer_id
ON customer_identifications(customer_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS customer_identifications;
DROP INDEX IF EXISTS customer_identifications_customer_id;
DROP TYPE IF EXISTS identification_type CASCADE;

-- +goose StatementEnd
