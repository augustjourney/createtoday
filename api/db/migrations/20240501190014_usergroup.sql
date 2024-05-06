-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_group (
    id SERIAL NOT NULL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES "group"(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    status VARCHAR(30) DEFAULT 'active',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    remove_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(group_id, user_id, status)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_group;
-- +goose StatementEnd
