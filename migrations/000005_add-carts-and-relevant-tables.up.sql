CREATE TABLE
    "carts" (
        "id" serial PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
        "updated_at" timestamptz NOT NULL DEFAULT (now ()),
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        UNIQUE ("user_id", "id")
    );

CREATE TABLE
    "cart_items" (
        "id" serial PRIMARY KEY,
        "product_id" bigint NOT NULL REFERENCES "products" ("id"),
        "cart_id" int NOT NULL REFERENCES "carts" ("id") ON DELETE CASCADE,
        "quantity" smallint NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        UNIQUE ("product_id", "cart_id")
    );

CREATE INDEX ON "cart_items" ("product_id", "cart_id");