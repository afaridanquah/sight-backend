-- +goose Up
-- +goose StatementBegin
CREATE TYPE outcome AS ENUM ('CLEARED', 'ATTENTION_NEEDED');

CREATE TYPE decision AS ENUM ('APPROVED', 'DECLINED', 'UNKNOWN');

CREATE TABLE verifications (
    id UUID PRIMARY KEY,
    verification_type VARCHAR(30),
    customer_id UUID NOT NULL,
    customer JSONB,
    business JSONB,
    business_id VARCHAR(50),
    org_id VARCHAR(50),
    creator_id UUID NOT NULL,
    outcome outcome,
    aml_insight JSONB,
    phone_insight JSONB,
    decision decision,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE verifications;

DROP TYPE outcome;

DROP TYPE decision;

-- +goose StatementEnd
