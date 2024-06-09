-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS offer (
   id SERIAL NOT NULL PRIMARY KEY,
   name VARCHAR(100) NOT NULL,
   description VARCHAR,
   slug VARCHAR(50) NOT NULL,
   price INTEGER NOT NULL DEFAULT 0,
   currency VARCHAR(30) NOT NULL DEFAULT 'RUB',
   is_free BOOLEAN DEFAULT true,
   send_order_created BOOLEAN DEFAULT true,
   send_order_completed BOOLEAN DEFAULT true,
   send_registration_email BOOLEAN DEFAULT false,
   registration_email VARCHAR,
   add_to_newsletter BOOLEAN DEFAULT false,
   type VARCHAR(30) NOT NULL DEFAULT 'one_time',
   settings JSON,
   project_id INTEGER REFERENCES project(id) ON DELETE CASCADE,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   registration_email_theme VARCHAR,
   success_message VARCHAR,
   redirect_url VARCHAR,
   ask_for_phone BOOLEAN DEFAULT false,
   ask_for_comment BOOLEAN DEFAULT false,
   oferta_url VARCHAR,
   agreement_url VARCHAR,
   privacy_url VARCHAR,
   send_welcome_email BOOLEAN DEFAULT true,
   ask_for_telegram BOOLEAN DEFAULT false,
   ask_for_instagram BOOLEAN DEFAULT false,
   can_use_promocode BOOLEAN DEFAULT false,
   is_donate BOOLEAN DEFAULT false,
   min_donate_price INTEGER DEFAULT 0,
   send_to_salebot BOOLEAN DEFAULT true,
   salebot_callback_text VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE offer;
-- +goose StatementEnd
