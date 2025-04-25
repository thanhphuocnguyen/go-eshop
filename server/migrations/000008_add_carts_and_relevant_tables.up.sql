CREATE TABLE
    "carts" (
        "id" UUID NOT NULL PRIMARY KEY,
        "user_id" UUID UNIQUE REFERENCES "users" ("id") ON DELETE CASCADE,
        "session_id" VARCHAR(255) UNIQUE,
        "order_id" UUID REFERENCES "orders" ("id") ON DELETE SET NULL,
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    cart_items (
        id UUID PRIMARY KEY,
        cart_id UUID NOT NULL REFERENCES "carts" ("id") ON DELETE CASCADE,
        variant_id UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE CASCADE,
        quantity smallint NOT NULL DEFAULT 1 CHECK ("quantity" > 0),
        added_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (cart_id, variant_id) -- This is a composite unique constraint
    );

-- Indexes for `carts` table
CREATE INDEX idx_carts_user_id ON carts (user_id);

CREATE INDEX idx_carts_status ON carts (order_id);

-- For finding logged-in user's cart
CREATE INDEX idx_carts_session_id ON carts (session_id);

-- For finding guest's cart
CREATE INDEX idx_carts_updated_at ON carts (updated_at);

-- For cleaning up old carts
-- Indexes for `cart_items` table
-- The UNIQUE constraint on (cart_id, variant_id) creates the most important index.
-- Add separate indexes if lookups by only one column are common and not covered well by the unique index:
CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);

-- Often needed to fetch all items for a cart
CREATE INDEX idx_cart_items_variant_id ON cart_items (variant_id);

-- Less common, but potentially useful