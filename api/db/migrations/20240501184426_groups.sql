-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "group" (
     id SERIAL NOT NULL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     description VARCHAR,
     settings JSON,
     project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "group";
-- +goose StatementEnd
