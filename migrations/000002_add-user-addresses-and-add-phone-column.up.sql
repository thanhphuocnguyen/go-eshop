CREATE TABLE
    "user_addresses" (
        "id" bigserial PRIMARY KEY,
        "user_id" bigint NOT NULL,
        "phone" varchar(12) NOT NULL,
        "address_1" varchar NOT NULL,
        "address_2" varchar,
        "ward" varchar(100),
        "district" varchar(100) NOT NULL,
        "city" varchar(100) NOT NULL,
        "is_primary" bool NOT NULL DEFAULT false,
        "is_deleted" bool NOT NULL DEFAULT false,
        "created_at" timestamptz NOT NULL DEFAULT now (),
        "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "deleted_at" timestamptz
    );

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_address_id") REFERENCES "user_addresses" ("id");

ALTER TABLE "user_addresses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");