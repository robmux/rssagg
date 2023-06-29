-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
        ENCODE(sha256(random()::text::bytea), 'hex')
        );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN api_key;
-- +goose StatementEnd
