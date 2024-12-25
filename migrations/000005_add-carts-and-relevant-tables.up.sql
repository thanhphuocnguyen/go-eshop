CREATE TABLE
    "carts" (
        "id" bigserial PRIMARY KEY,
        "checkout_at" timestamptz,
        "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
        "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        UNIQUE ("user_id", "checkout_at")
    );

CREATE TABLE
    "cart_items" (
        "id" bigserial PRIMARY KEY,
        "product_id" bigint NOT NULL REFERENCES "products" ("id"),
        "cart_id" bigint NOT NULL REFERENCES "carts" ("id") ON DELETE CASCADE,
        "quantity" smallint NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        UNIQUE ("product_id", "cart_id")
    );

CREATE INDEX ON "cart_items" ("product_id", "cart_id");

CREATE UNIQUE INDEX unique_user_cart ON carts (user_id)
WHERE
    checkout_at IS NULL;