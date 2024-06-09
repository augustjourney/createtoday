-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order" (
     id SERIAL NOT NULL PRIMARY KEY,
     description VARCHAR(250),
     comment VARCHAR(250),
     price INTEGER NOT NULL,
     currency VARCHAR(30) NOT NULL DEFAULT 'RUB',
     status VARCHAR(50) NOT NULL DEFAULT 'pending',
     error JSON,
     card_info JSON,
     payment_id VARCHAR(150),
     integration_id INTEGER REFERENCES pay_integration(id) ON DELETE SET NULL,
     offer_id INTEGER REFERENCES offer(id) ON DELETE SET NULL,
     project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
     user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
     note VARCHAR,
     is_gift BOOLEAN DEFAULT false,
     gifted_to JSON,
     salebot_client_id VARCHAR,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "order";
-- +goose StatementEnd
