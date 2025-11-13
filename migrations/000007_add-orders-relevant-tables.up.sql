CREATE TABLE
    "orders" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "customer_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE RESTRICT,
        "customer_email" VARCHAR(255) NOT NULL,
        "customer_name" VARCHAR(255) NOT NULL,
        "customer_phone" VARCHAR(50) NOT NULL,
        "shipping_address" JSONB NOT NULL,
        "total_price" DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
        "status" order_status NOT NULL DEFAULT 'pending',
        "confirmed_at" TIMESTAMPTZ,
        "delivered_at" TIMESTAMPTZ,
        "cancelled_at" TIMESTAMPTZ,
        "shipping_method" VARCHAR(100),
        "refunded_at" TIMESTAMPTZ,
        "order_date" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    "order_items" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "order_id" UUID NOT NULL REFERENCES "orders" ("id") ON DELETE CASCADE,
        "variant_id" UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE RESTRICT,
        "quantity" SMALLINT NOT NULL DEFAULT 1 CHECK (quantity > 0),
        "price_per_unit_snapshot" DECIMAL(10, 2) NOT NULL CHECK (price_per_unit_snapshot >= 0),
        "line_total_snapshot" DECIMAL(10, 2) NOT NULL CHECK (line_total_snapshot >= 0),
        "product_name_snapshot" VARCHAR(255) NOT NULL,
        "variant_sku_snapshot" VARCHAR(100) NOT NULL,
        "attributes_snapshot" JSONB NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Indexes for `orders` table
CREATE INDEX idx_orders_customer_id ON orders (customer_id);

-- Crucial FK index (if used)
CREATE INDEX idx_orders_order_date ON orders (order_date);

CREATE INDEX idx_orders_status ON orders (status);

CREATE INDEX idx_orders_customer_email ON orders (customer_email);

-- For guest lookups
CREATE INDEX idx_orders_created_at ON orders (created_at);

-- Indexes for `order_items` table
CREATE INDEX idx_order_items_order_id ON order_items (order_id);

-- Crucial FK index
CREATE INDEX idx_order_items_variant_id ON order_items (variant_id);

-- Crucial FK index
-- Optional: Index snapshot fields if frequently searched directly
-- CREATE INDEX idx_order_items_variant_sku_snapshot ON order_items (variant_sku_snapshot);
-- CREATE INDEX idx_order_items_attributes_snapshot ON order_items USING GIN (attributes_snapshot); -- Example for PostgreSQL JSONB indexing