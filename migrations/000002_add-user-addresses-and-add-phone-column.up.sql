ALTER TABLE "users" ADD "phone" varchar NOT NULL;

CREATE TABLE
    "user_addresses" (
        "id" bigserial PRIMARY KEY,
        "user_id" bigint NOT NULL,
        "phone" varchar NOT NULL,
        "address_1" varchar NOT NULL,
        "address_2" varchar,
        "ward" varchar,
        "district" varchar,
        "city" varchar NOT NULL,
        "is_primary" bool NOT NULL DEFAULT false
    );

ALTER TABLE "user_addresses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");