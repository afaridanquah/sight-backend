-- +goose Up
-- +goose StatementBegin
ALTER TABLE businesses
ADD COLUMN registration_number VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS businesses;
-- +goose StatementEnd
