-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS media (
       id SERIAL NOT NULL PRIMARY KEY,
       name VARCHAR,
       type VARCHAR(30) DEFAULT 'image',
       slug VARCHAR(100) NOT NULL,
       mime VARCHAR(30),
       ext VARCHAR(30),
       size INTEGER,
       width INTEGER,
       height INTEGER,
       url VARCHAR,
       sources JSON,
       blurhash JSON,
       parent_id INTEGER REFERENCES media(id) ON DELETE CASCADE,
       storage VARCHAR (50) NOT NULL DEFAULT 'selectel',
       duration INTEGER,
       bucket VARCHAR (50),
       status VARCHAR(50),
       caption VARCHAR,
       original BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE media;
-- +goose StatementEnd
