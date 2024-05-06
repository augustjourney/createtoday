-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product (
     id SERIAL NOT NULL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     slug VARCHAR(100) NOT NULL,
     description VARCHAR,
     layout VARCHAR(30) NOT NULL DEFAULT 'all_published',
     position INTEGER NOT NULL DEFAULT 0,
     access VARCHAR(30) DEFAULT 'nobody',
     is_published BOOLEAN DEFAULT FALSE,
     cover JSON,
     parent_id INTEGER REFERENCES product(id) ON DELETE SET NULL,
     project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
     settings JSON,
     show_lessons_without_access BOOLEAN DEFAULT FALSE,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     created_by INTEGER REFERENCES "user"(id) ON DELETE SET NULL,
     UNIQUE(slug, project_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product;
-- +goose StatementEnd
