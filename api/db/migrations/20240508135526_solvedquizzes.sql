-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quiz_solved
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER     NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    quiz_id     INTEGER     NOT NULL REFERENCES quiz (id) ON DELETE CASCADE,
    product_id  INTEGER     NOT NULL REFERENCES product (id) ON DELETE CASCADE,
    lesson_id   INTEGER     NOT NULL REFERENCES lesson (id) ON DELETE CASCADE,
    project_id  INTEGER     NOT NULL REFERENCES project (id) ON DELETE CASCADE,
    user_answer JSON,
    type        VARCHAR(30) NOT NULL,
    starred     BOOLEAN                  DEFAULT FALSE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, quiz_id, product_id, project_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz_solved;
-- +goose StatementEnd
