-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id UUID NOT NULL PRIMARY KEY,
    items_count INT NOT NULL,
    pack_count INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
