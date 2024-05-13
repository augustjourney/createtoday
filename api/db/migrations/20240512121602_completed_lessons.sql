-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS completed_lessons (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES product(id) ON DELETE CASCADE,
    lesson_id INTEGER REFERENCES lesson(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id, lesson_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE completed_lessons;
-- +goose StatementEnd
