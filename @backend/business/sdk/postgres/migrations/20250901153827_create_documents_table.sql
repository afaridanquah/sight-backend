-- +goose Up
-- +goose StatementBegin
CREATE TYPE document_status AS ENUM ('PENDING', 'REJECTED', 'APPROVED', 'DRAFT');

CREATE TABLE documents(
    id UUID PRIMARY KEY,
    customer_id UUID,
    business_id UUID,
    original_name TEXT,
    filename VARCHAR(255),
    mimetype VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    status document_status,
    creator_id UUID,
    metadata JSON,
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents;
DROP TYPE document_status;
-- +goose StatementEnd
