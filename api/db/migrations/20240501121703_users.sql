-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL NOT NULL PRIMARY KEY,
    email VARCHAR(80) NOT NULL,
    password VARCHAR(80) NOT NULL,
    first_name VARCHAR(80),
    last_name VARCHAR(80),
    avatar VARCHAR(80),
    phone VARCHAR(30) DEFAULT null,
    about TEXT DEFAULT null,
    telegram VARCHAR(100) DEFAULT null,
    instagram VARCHAR(100) DEFAULT null,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             UNIQUE(email)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
