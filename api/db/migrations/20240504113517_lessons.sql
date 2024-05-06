-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lesson
(
    id             SERIAL       NOT NULL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    slug           VARCHAR(100) NOT NULL,
    description    VARCHAR,
    content        JSON,
    position       INTEGER      NOT NULL    DEFAULT 0,
    is_published   BOOLEAN                  DEFAULT FALSE,
    is_stop_lesson BOOLEAN                  DEFAULT FALSE,
    can_complete   BOOLEAN                  DEFAULT FALSE,
    settings       JSON,
    is_public      BOOLEAN                  DEFAULT FALSE,
    is_deleted     BOOLEAN                  DEFAULT FALSE,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by     INTEGER      REFERENCES "user" (id) ON DELETE SET NULL,
    project_id     INTEGER REFERENCES project (id) ON DELETE CASCADE,
    product_id     INTEGER REFERENCES product (id) ON DELETE CASCADE,
    publish_at     TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE (slug, project_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lesson;
-- +goose StatementEnd
