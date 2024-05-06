-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS project (
     id SERIAL NOT NULL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     domain VARCHAR(100) NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     owner_id INTEGER REFERENCES "user"(id) ON DELETE SET NULL,
     UNIQUE(domain)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project;
-- +goose StatementEnd
