
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user_detail
(
    id TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    password TEXT,
    is_email_verified BOOLEAN,
    dob TIMESTAMP WITHOUT TIME ZONE,
    last_login_at TIMESTAMP WITHOUT TIME ZONE,
    active BOOLEAN,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT user_email_uk UNIQUE (email)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

