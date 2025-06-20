-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('ACTIVE', 'INACTIVE', 'BLOCKED', 'SUSPENDED');

CREATE TABLE businesses (
    id UUID,
    name TEXT NOT NULL,
    owner_id UUID NOT NULL,
    status status,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE businesses;

DROP TYPE status;

-- +goose StatementEnd
