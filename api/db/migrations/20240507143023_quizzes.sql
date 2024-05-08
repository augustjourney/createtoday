-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quiz (
      id SERIAL NOT NULL PRIMARY KEY,
      name VARCHAR(100),
      slug VARCHAR(100) NOT NULL,
      content JSON,
      type VARCHAR(50) NOT NULL,
      settings JSON,
      show_others_answers BOOLEAN DEFAULT false,
      lesson_id INTEGER REFERENCES lesson(id) ON DELETE CASCADE,
      created_by INTEGER REFERENCES "user"(id) ON DELETE SET NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
      product_id INTEGER REFERENCES product(id) ON DELETE CASCADE,
      UNIQUE(slug, project_id, lesson_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz;
-- +goose StatementEnd
