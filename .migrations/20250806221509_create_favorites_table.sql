-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
    CREATE TABLE favorites (
        client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        product_id INT NOT NULL,
        registred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE (client_id, product_id)
    );

CREATE INDEX IF NOT EXISTS idx_client_id ON favorites (client_id);

CREATE INDEX IF NOT EXISTS idx_client_product ON favorites (client_id, product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
    DROP TABLE IF EXISTS favorites;
    DROP INDEX IF EXISTS idx_client_id;
    DROP INDEX IF EXISTS idx_client_product;
-- +goose StatementEnd
