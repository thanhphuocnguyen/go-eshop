CREATE TABLE verify_emails (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    verify_code VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at timestamptz NOT NULL DEFAULT (now ()),
    expired_at timestamptz NOT NULL DEFAULT (now () + interval '1 day')
);