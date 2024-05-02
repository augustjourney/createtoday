-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_group (
   id SERIAL NOT NULL PRIMARY KEY,
   group_id INTEGER NOT NULL REFERENCES "group"(id) ON DELETE CASCADE,
   product_id INTEGER NOT NULL REFERENCES product(id) ON DELETE CASCADE,
   all_lessons_access BOOLEAN DEFAULT TRUE,
   no_access_content JSON,
   UNIQUE(group_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_group;
-- +goose StatementEnd
