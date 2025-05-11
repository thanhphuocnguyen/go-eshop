-- Create discounts table
CREATE TABLE IF NOT EXISTS discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value NUMERIC(10, 2) NOT NULL,
    min_purchase_amount NUMERIC(10, 2),
    max_discount_amount NUMERIC(10, 2),
    usage_limit INTEGER,
    used_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create discount_products table for product-specific discounts
CREATE TABLE IF NOT EXISTS discount_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, product_id)
);

-- Create discount_categories table for category-specific discounts
CREATE TABLE IF NOT EXISTS discount_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, category_id)
);

-- Create discount_users table for user-specific discounts
CREATE TABLE IF NOT EXISTS discount_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, user_id)
);

-- Create order_discounts table to track discount usage in orders
CREATE TABLE IF NOT EXISTS order_discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    discount_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (order_id, discount_id)
);

-- Add triggers for updating timestamps
CREATE OR REPLACE FUNCTION update_discount_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_discount_timestamp
BEFORE UPDATE ON discounts
FOR EACH ROW
EXECUTE FUNCTION update_discount_updated_at();

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_discounts_code ON discounts(code);
CREATE INDEX IF NOT EXISTS idx_discounts_is_active ON discounts(is_active);
CREATE INDEX IF NOT EXISTS idx_discounts_dates ON discounts(starts_at, expires_at);
CREATE INDEX IF NOT EXISTS idx_discount_products_discount_id ON discount_products(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_products_product_id ON discount_products(product_id);
CREATE INDEX IF NOT EXISTS idx_discount_categories_discount_id ON discount_categories(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_categories_category_id ON discount_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_discount_users_discount_id ON discount_users(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_users_user_id ON discount_users(user_id);
CREATE INDEX IF NOT EXISTS idx_order_discounts_order_id ON order_discounts(order_id);
CREATE INDEX IF NOT EXISTS idx_order_discounts_discount_id ON order_discounts(discount_id);

-- Add discounted_price column to order_items table to track original vs. discounted price
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS discounted_price NUMERIC(10, 2);
