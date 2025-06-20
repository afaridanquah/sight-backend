-- +goose Up
-- +goose StatementBegin
CREATE TYPE identification_type AS ENUM (
    'PASSPORT',
    'DRIVERS_LICENSE',
    'NATIONAL_ID',
    'RESIDENT_PERMIT',
    'SSN'
);

CREATE TABLE identifications (
    id UUID NOT NULL,
    first_name VARCHAR(45),
    last_name VARCHAR(45),
    middle_name VARCHAR(45),
    other_names VARCHAR(45),
    pin TEXT NOT NULL,
    identification_type identification_type NOT NULL,
    issued_country VARCHAR(2),
    issued_date DATE,
    place_of_birth VARCHAR(45),
    place_of_issue VARCHAR(45),
    date_of_birth DATE,
    address_1 TEXT,
    address_2 TEXT,
    city TEXT,
    state_region TEXT,
    zip_code VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE identifications;

DROP TYPE identification_type;

-- +goose StatementEnd
