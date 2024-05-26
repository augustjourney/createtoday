-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quiz_comment (
      id SERIAL PRIMARY KEY,
      author_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
      quiz_solved_id INTEGER NOT NULL REFERENCES quiz_solved(id) ON DELETE CASCADE,
      uuid VARCHAR NOT NULL,
      text VARCHAR,
      is_read BOOLEAN DEFAULT FALSE,
      is_edited BOOLEAN DEFAULT FALSE,
      is_from_moderator BOOLEAN DEFAULT FALSE,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz_comment;
-- +goose StatementEnd
