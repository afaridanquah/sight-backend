-- +goose Up
-- +goose StatementBegin
CREATE TABLE stored_events (
    id UUID,
    type VARCHAR(50),
    aggregate_id UUID,
    aggregate_type TEXT,
    aggregate_version BIGINT,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    meta_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    occured_at TIMESTAMP,
    registered_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stored_events;
-- +goose StatementEnd
