CREATE TABLE
    "users" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        "role" user_role NOT NULL DEFAULT 'user',
        "username" VARCHAR UNIQUE NOT NULL,
        "email" VARCHAR UNIQUE NOT NULL,
        "phone_number" VARCHAR(20) NOT NULL CHECK (
            char_length(phone_number) >= 10
            AND char_length(phone_number) <= 20
        ),
        "first_name" VARCHAR NOT NULL,
        "last_name" VARCHAR NOT NULL,
        "avatar_url" VARCHAR,
        "avatar_image_id" VARCHAR,
        "hashed_password" VARCHAR NOT NULL,
        "verified_email" bool NOT NULL DEFAULT FALSE,
        "verified_phone" bool NOT NULL DEFAULT FALSE,
        "locked" BOOLEAN NOT NULL DEFAULT FALSE,
        "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    "user_addresses" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
        "phone_number" VARCHAR(20) NOT NULL CHECK (
            char_length(phone_number) >= 10
            AND char_length(phone_number) <= 20
        ),
        "street" VARCHAR NOT NULL,
        "ward" VARCHAR(100),
        "district" VARCHAR(100) NOT NULL,
        "city" VARCHAR(100) NOT NULL,
        "default" BOOLEAN NOT NULL DEFAULT FALSE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    "email_verifications" (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        verify_code VARCHAR(255) NOT NULL,
        is_used BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        expired_at TIMESTAMPTZ NOT NULL DEFAULT (NOW () + interval '1 day')
    );

CREATE INDEX idx_email_verifications_expired_at ON email_verifications (expired_at);
CREATE TABLE
    "user_sessions" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        "user_id" UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        "refresh_token" VARCHAR NOT NULL,
        "user_agent" VARCHAR(512) NOT NULL,
        "client_ip" INET NOT NULL,
        "blocked" boolean NOT NULL DEFAULT FALSE,
        "expired_at" TIMESTAMPTZ NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    user_payment_infos (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        "user_id" UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        "card_number" VARCHAR(19) NOT NULL,
        "card_last4" VARCHAR(4) NOT NULL,
        "payment_method_token" VARCHAR(255) NOT NULL,
        "expiration_date" DATE NOT NULL,
        "billing_address" TEXT NOT NULL,
        "default" BOOLEAN DEFAULT FALSE,
        "created_at" TIMESTAMPTZ DEFAULT NOW (),
        "updated_at" TIMESTAMPTZ DEFAULT NOW (),
        UNIQUE ("user_id", "payment_method_token")
    );

CREATE INDEX ON "user_payment_infos" ("user_id");

CREATE INDEX ON "user_sessions" ("user_id");

CREATE INDEX ON "user_addresses" ("user_id", "default");
