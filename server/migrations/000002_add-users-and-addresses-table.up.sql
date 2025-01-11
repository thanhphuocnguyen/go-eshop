CREATE TABLE
    "users" (
        "user_id" bigserial PRIMARY KEY,
        "role" user_role NOT NULL DEFAULT 'user',
        "username" varchar UNIQUE NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "phone" varchar(20) NOT NULL CHECK (
            char_length(phone) >= 10
            AND char_length(phone) <= 20
        ),
        "fullname" varchar NOT NULL,
        "hashed_password" varchar NOT NULL,
        "verified_email" bool NOT NULL DEFAULT false,
        "verified_phone" bool NOT NULL DEFAULT false,
        "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "updated_at" timestamptz NOT NULL DEFAULT now (),
        "created_at" timestamptz NOT NULL DEFAULT now ()
    );

CREATE TABLE
    "user_addresses" (
        "user_address_id" bigserial PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES "users" ("user_id") ON DELETE CASCADE,
        "phone" varchar(20) NOT NULL CHECK (
            char_length(phone) >= 10
            AND char_length(phone) <= 20
        ),
        "street" varchar NOT NULL,
        "ward" varchar(100),
        "district" varchar(100) NOT NULL,
        "city" varchar(100) NOT NULL,
        "default" BOOLEAN NOT NULL DEFAULT false,
        "deleted" BOOLEAN NOT NULL DEFAULT false,
        "created_at" timestamptz NOT NULL DEFAULT now (),
        "updated_at" timestamptz NOT NULL DEFAULT now ()
    );

CREATE TABLE
    verify_emails (
        id SERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        verify_code VARCHAR(255) NOT NULL,
        is_used BOOLEAN NOT NULL DEFAULT FALSE,
        created_at timestamptz NOT NULL DEFAULT (now ()),
        expired_at timestamptz NOT NULL DEFAULT (now () + interval '1 day')
    );

CREATE TABLE
    "sessions" (
        "session_id" uuid PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
        "refresh_token" varchar NOT NULL,
        "user_agent" varchar NOT NULL,
        "client_ip" varchar NOT NULL,
        "blocked" boolean NOT NULL DEFAULT false,
        "expired_at" timestamptz NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT now ()
    );

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "user_addresses" ("user_id", "default");

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");