-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS related_media (
       id SERIAL PRIMARY KEY NOT NULL,
       media_id INTEGER REFERENCES media(id) ON DELETE CASCADE,
       related_id INTEGER NOT NULL,
       related_type VARCHAR(30) NOT NULL,
       project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       UNIQUE(media_id, related_id, related_type, project_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE related_media;
-- +goose StatementEnd
