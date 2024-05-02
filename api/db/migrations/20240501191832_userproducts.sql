-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW _userproducts AS (
    SELECT
        ug.user_id, p.id, p.name, p.slug, p.description, p.settings,
        p.parent_id, p.cover, p.layout, p.show_lessons_without_access, p.project_id, p.position
    FROM product_group AS pg
    JOIN user_group AS ug ON ug.group_id = pg.group_id
    JOIN product AS p ON p.id = pg.product_id AND p.is_published IS TRUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW _userproducts;
-- +goose StatementEnd
