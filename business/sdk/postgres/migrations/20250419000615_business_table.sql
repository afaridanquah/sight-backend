-- +goose Up
-- +goose StatementBegin
-- CREATE TYPE status AS ENUM ('ACTIVE', 'INACTIVE', 'BLOCKED', 'SUSPENDED');
CREATE TYPE status AS ENUM ('APPROVED', 'PENDING', 'REJECTED', 'DRAFT');
CREATE TYPE entity AS ENUM ('ESTATE', 'SOLE_PROPRIETOR', 'CORPORATION', 'EXEMPT_ORGANIZATION');


CREATE TABLE businesses (
    id UUID PRIMARY KEY,
    org_id UUID,
    legal_name TEXT NOT NULL,
    entity entity NOT NULL,
    tax_id TEXT,
    dba TEXT NOT NULL,
    jurisdiction VARCHAR(2) NOT NULL,
    admin_id UUID NOT NULL,
    owners JSONB,
    address JSONB,
    website VARCHAR(200),
    phone_numbers JSONB DEFAULT '[]'::jsonb,
    email_addresses JSONB DEFAULT '[]'::jsonb,
    documents JSONB DEFAULT '[]'::jsonb,
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
