CREATE TABLE
    "users" (
        "id" bigserial PRIMARY KEY,
        "role" user_role NOT NULL DEFAULT 'user',
        "username" varchar UNIQUE NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "phone" varchar(12) NOT NULL,
        "full_name" varchar NOT NULL,
        "hashed_password" varchar NOT NULL,
        "verified_email" bool NOT NULL DEFAULT false,
        "verified_phone" bool NOT NULL DEFAULT false,
        "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "updated_at" timestamptz NOT NULL DEFAULT now(),
        "created_at" timestamptz NOT NULL DEFAULT now()
    );

CREATE TABLE
    "user_addresses" (
        "id" bigserial PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
        "phone" varchar(12) NOT NULL,
        "address_1" varchar NOT NULL,
        "address_2" varchar,
        "ward" varchar(100),
        "district" varchar(100) NOT NULL,
        "city" varchar(100) NOT NULL,
        "is_primary" bool NOT NULL DEFAULT false,
        "is_deleted" bool NOT NULL DEFAULT false,
        "created_at" timestamptz NOT NULL DEFAULT now(),
        "updated_at" timestamptz NOT NULL DEFAULT now(),
        "deleted_at" timestamptz
    );

CREATE INDEX ON "user_addresses" ("user_id", "is_primary");

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");