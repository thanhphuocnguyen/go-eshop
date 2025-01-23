CREATE TABLE
    "orders" (
        "order_id" UUID NOT NULL PRIMARY KEY,
        "user_id" UUID NOT NULL REFERENCES "users" ("user_id"),
        "user_address_id" BIGINT NOT NULL REFERENCES "user_addresses" ("user_address_id"),
        "total_price" DECIMAL(10, 2) NOT NULL,
        "status" order_status NOT NULL DEFAULT 'pending',
        "confirmed_at" TIMESTAMPTZ,
        "delivered_at" TIMESTAMPTZ,
        "cancelled_at" TIMESTAMPTZ,
        "refunded_at" TIMESTAMPTZ,
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    "order_items" (
        "order_item_id" BIGSERIAL PRIMARY KEY,
        "order_id" UUID NOT NULL REFERENCES "orders" ("order_id") ON DELETE CASCADE,
        "product_id" BIGINT NOT NULL REFERENCES "products" ("product_id"),
        "variant_id" BIGINT NOT NULL REFERENCES "product_variants" ("variant_id"),
        "quantity" SMALLINT NOT NULL,
        "price" DECIMAL(10, 2) NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CHECK ("quantity" > 0),
        CHECK ("price" > 0),
        CHECK (
            "product_id" IS NOT NULL
            OR "variant_id" IS NOT NULL
        )
    );

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "orders" ("user_id");

CREATE INDEX ON "orders" ("user_id", "status");

CREATE INDEX ON "order_items" ("product_id", "order_id");