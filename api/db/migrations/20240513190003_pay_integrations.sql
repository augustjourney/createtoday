-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pay_integration (
     id SERIAL NOT NULL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     type VARCHAR(100) NOT NULL,
     login VARCHAR  NOT NULL,
     password VARCHAR  NOT NULL,
     is_active BOOLEAN DEFAULT FALSE,
     send_receipt BOOLEAN DEFAULT FALSE,
     receipt_settings JSON,
     project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pay_integration;
-- +goose StatementEnd
