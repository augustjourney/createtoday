-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS offer_group (
    id SERIAL NOT NULL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES "group"(id) ON DELETE CASCADE,
    offer_id INTEGER NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    UNIQUE(group_id, offer_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE offer_group;
-- +goose StatementEnd
