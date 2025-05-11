CREATE TABLE
    "users" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    "user_addresses" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
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
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    verify_emails (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        verify_code VARCHAR(255) NOT NULL,
        is_used BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        expired_at TIMESTAMPTZ NOT NULL DEFAULT (NOW () + interval '1 day')
    );

CREATE TABLE
    "sessions" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "user_id" UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        "refresh_token" VARCHAR NOT NULL,
        "user_agent" VARCHAR NOT NULL,
        "client_ip" VARCHAR NOT NULL,
        "blocked" boolean NOT NULL DEFAULT FALSE,
        "expired_at" TIMESTAMPTZ NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    user_payment_infos (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "user_id" UUID REFERENCES users (id) ON DELETE CASCADE,
        "card_number" VARCHAR(16) NOT NULL,
        "cardholder_name" VARCHAR(100) NOT NULL,
        "expiration_date" DATE NOT NULL,
        "billing_address" TEXT NOT NULL,
        "default" BOOLEAN DEFAULT FALSE,
        "created_at" TIMESTAMPTZ DEFAULT NOW (),
        "updated_at" TIMESTAMPTZ DEFAULT NOW (),
        UNIQUE ("user_id", "card_number")
    );

CREATE INDEX ON "user_payment_infos" ("user_id");

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "user_addresses" ("user_id", "default");

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");