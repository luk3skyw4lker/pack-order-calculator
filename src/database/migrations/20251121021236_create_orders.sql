-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id UUID NOT NULL PRIMARY KEY,
    items_count INT NOT NULL,
    pack_setup VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
