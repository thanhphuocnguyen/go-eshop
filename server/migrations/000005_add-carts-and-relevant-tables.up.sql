CREATE TABLE
    "carts" (
        "cart_id" serial PRIMARY KEY,
        "user_id" bigint NOT NULL UNIQUE REFERENCES "users" ("user_id") ON DELETE CASCADE,
        "updated_at" timestamptz NOT NULL DEFAULT (now ()),
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "cart_items" (
        "cart_item_id" serial PRIMARY KEY,
        "product_id" bigint NOT NULL REFERENCES "products" ("product_id"),
        "cart_id" int NOT NULL REFERENCES "carts" ("cart_id") ON DELETE CASCADE,
        "quantity" smallint NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        UNIQUE ("product_id", "cart_id")
    );

CREATE INDEX ON "cart_items" ("product_id", "cart_id");