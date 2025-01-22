CREATE TABLE
    "carts" (
        "cart_id" UUID NOT NULL PRIMARY KEY,
        "user_id" UUID NOT NULL UNIQUE REFERENCES "users" ("user_id") ON DELETE CASCADE,
        "updated_at" timestamptz NOT NULL DEFAULT (now ()),
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "cart_items" (
        "cart_item_id" serial PRIMARY KEY,
        "cart_id" UUID NOT NULL REFERENCES "carts" ("cart_id") ON DELETE CASCADE,
        "product_id" bigint NOT NULL REFERENCES "products" ("product_id") ON DELETE CASCADE,
        "variant_id" bigint NOT NULL REFERENCES "product_variants" ("variant_id") ON DELETE CASCADE,
        "quantity" smallint NOT NULL DEFAULT 1 CHECK ("quantity" > 0),
        "updated_at" timestamptz NOT NULL DEFAULT now(),
        "created_at" timestamptz NOT NULL DEFAULT now(),
        UNIQUE ("cart_id", "variant_id", "product_id") -- This is a composite unique constraint
    );

CREATE INDEX ON "cart_items" ("product_id", "cart_id");