CREATE TABLE
    "users" (
        "user_id" UUID NOT NULL PRIMARY KEY,
        "role" user_role NOT NULL DEFAULT 'user',
        "username" VARCHAR UNIQUE NOT NULL,
        "email" VARCHAR UNIQUE NOT NULL,
        "phone" VARCHAR(20) NOT NULL CHECK (
            char_length(phone) >= 10
            AND char_length(phone) <= 20
        ),
        "fullname" VARCHAR NOT NULL,
        "hashed_password" VARCHAR NOT NULL,
        "verified_email" bool NOT NULL DEFAULT FALSE,
        "verified_phone" bool NOT NULL DEFAULT FALSE,
        "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    "user_addresses" (
        "user_address_id" BIGSERIAL PRIMARY KEY,
        "user_id" UUID NOT NULL REFERENCES "users" ("user_id") ON DELETE CASCADE,
        "phone" VARCHAR(20) NOT NULL CHECK (
            char_length(phone) >= 10
            AND char_length(phone) <= 20
        ),
        "street" VARCHAR NOT NULL,
        "ward" VARCHAR(100),
        "district" VARCHAR(100) NOT NULL,
        "city" VARCHAR(100) NOT NULL,
        "default" BOOLEAN NOT NULL DEFAULT FALSE,
        "deleted" BOOLEAN NOT NULL DEFAULT FALSE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    verify_emails (
        id SERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        verify_code VARCHAR(255) NOT NULL,
        is_used BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        expired_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + interval '1 day')
    );

CREATE TABLE
    "sessions" (
        "session_id" UUID NOT NULL PRIMARY KEY,
        "user_id" UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
        "refresh_token" VARCHAR NOT NULL,
        "user_agent" VARCHAR NOT NULL,
        "client_ip" VARCHAR NOT NULL,
        "blocked" boolean NOT NULL DEFAULT FALSE,
        "expired_at" TIMESTAMPTZ NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "user_addresses" ("user_id", "default");

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");