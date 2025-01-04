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

CREATE INDEX ON "user_addresses" ("user_id", "default");

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");