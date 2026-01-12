-- +goose Up
-- +goose StatementBegin
CREATE TABLE gov_identifications(
    id CHAR(31) PRIMARY KEY NOT NULL,
    first_name VARCHAR(45),
    last_name VARCHAR(45),
    middle_name VARCHAR(45),
    other_names VARCHAR(45),
    pin VARCHAR(100) NOT NULL,
    card_number VARCHAR(255),
    identification_type identification_type NOT NULL,
    occupation TEXT,
    issued_country VARCHAR(2),
    issued_date DATE,
    nationality VARCHAR(2),
    place_of_birth VARCHAR(45),
    place_of_issue VARCHAR(45),
    digital_address VARCHAR(255),
    mother_maiden_name VARCHAR(255),
    date_of_birth DATE,
    phone_number_1 VARCHAR(20),
    phone_number_2 VARCHAR(20),
    address_1 TEXT,
    address_2 TEXT,
    city VARCHAR(150),
    state_region VARCHAR(150),
    zip_code VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(card_number, pin)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gov_identifications;
-- +goose StatementEnd
