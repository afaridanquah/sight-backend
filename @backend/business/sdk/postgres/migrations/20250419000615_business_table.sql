-- +goose Up
-- +goose StatementBegin
-- CREATE TYPE status AS ENUM ('ACTIVE', 'INACTIVE', 'BLOCKED', 'SUSPENDED');
CREATE TYPE status AS ENUM ('APPROVED', 'PENDING', 'REJECTED', 'DRAFT');
CREATE TYPE entity AS ENUM ('ESTATE', 'SOLE_PROPRIETOR', 'CORPORATION', 'EXEMPT_ORGANIZATION');


CREATE TABLE businesses (
    id VARCHAR(50) PRIMARY KEY,
    org_id VARCHAR(50),
    legal_name TEXT NOT NULL,
    entity entity NOT NULL,
    tax_id TEXT,
    dba TEXT NOT NULL,
    jurisdiction VARCHAR(2) NOT NULL,
    admin_id VARCHAR(50) NOT NULL,
    owners JSONB,
    address JSONB,
    website VARCHAR(200),
    phone_numbers JSONB,
    email_addresses JSONB,
    documents JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE businesses;

DROP TYPE status;

DROP TYPE entity;

-- +goose StatementEnd
