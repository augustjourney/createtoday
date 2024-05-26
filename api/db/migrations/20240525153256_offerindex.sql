-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS offer_slug_idx ON public.offer (slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX offer_slug_idx;
-- +goose StatementEnd
