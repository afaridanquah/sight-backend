-- +goose Up
-- +goose StatementBegin
CREATE TABLE customer_identifications (
    id UUID,
    customer_id UUID PRIMARY KEY NOT NULL,
    identification_type identification_type NOT NULL,
    issued_country VARCHAR(2),
    pin VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE customer_identifications;

-- +goose StatementEnd
