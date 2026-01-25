-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses(
    id CHAR(31) PRIMARY KEY,
    entity_id CHAR(31),
    line_1 VARCHAR(255),
    line_2 VARCHAR(255),
    state VARCHAR(100),
    city VARCHAR(100),
    zip VARCHAR(70)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd
