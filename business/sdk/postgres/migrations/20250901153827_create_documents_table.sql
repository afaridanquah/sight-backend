-- +goose Up
-- +goose StatementBegin
CREATE TYPE document_status AS ENUM ('PENDING', 'REJECTED', 'APPROVED', 'DRAFT');

CREATE TABLE documents(
    id CHAR(31) PRIMARY KEY,
    customer_id CHAR(31),
    business_id CHAR(31),
    original_name TEXT,
    filename VARCHAR(255),
    mimetype VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    status document_status,
    creator_id CHAR(31),
    metadata JSON,
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents;
DROP TYPE document_status;
-- +goose StatementEnd
