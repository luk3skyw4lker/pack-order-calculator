-- +goose Up
-- +goose StatementBegin
CREATE TABLE pack_sizes (
    id UUID NOT NULL PRIMARY KEY,
    size INT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pack_sizes;
-- +goose StatementEnd
