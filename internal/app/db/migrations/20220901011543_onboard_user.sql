-- +goose Up
-- +goose StatementBegin
CREATE TABLE onboard_user
(
    id bigserial,
    email text,
    verified boolean,
    created_at timestamp without time zone DEFAULT NOW(),
    verified_at timestamp without time zone,
    PRIMARY KEY (id),
    CONSTRAINT onboard_email_uk UNIQUE (email)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
