-- Create shipping_methods table
CREATE TABLE IF NOT EXISTS
    shipping_methods (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) NOT NULL,
        description TEXT,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        requires_address BOOLEAN NOT NULL DEFAULT TRUE,
        estimated_delivery_time VARCHAR(100),
        icon_url TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create shipping_zones table to handle geographical regions
CREATE TABLE IF NOT EXISTS
    shipping_zones (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) NOT NULL,
        description TEXT,
        countries TEXT[] NOT NULL,
        states TEXT[],
        zip_codes TEXT[],
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create shipping_rates table for pricing
CREATE TABLE IF NOT EXISTS
    shipping_rates (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        shipping_method_id UUID NOT NULL REFERENCES shipping_methods(id) ON DELETE CASCADE,
        shipping_zone_id UUID NOT NULL REFERENCES shipping_zones(id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        base_rate DECIMAL(10, 2) NOT NULL,
        min_order_amount DECIMAL(10, 2),
        max_order_amount DECIMAL(10, 2),
        free_shipping_threshold DECIMAL(10, 2),
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE(shipping_method_id, shipping_zone_id)
    );

-- Create shipping_rate_conditions table for special pricing rules
CREATE TABLE IF NOT EXISTS
    shipping_rate_conditions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        shipping_rate_id UUID NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
        condition_type VARCHAR(50) NOT NULL, -- e.g., 'weight_range', 'item_count', 'product_category'
        min_value DECIMAL(10, 2),
        max_value DECIMAL(10, 2),
        additional_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
        category_ids UUID[],
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Add shipping-related fields to orders table
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS shipping_method_id UUID REFERENCES shipping_methods(id);
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS shipping_rate_id UUID REFERENCES shipping_rates(id);
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS estimated_delivery_date TIMESTAMPTZ;
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS tracking_url TEXT;
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS shipping_provider VARCHAR(100);
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS shipping_notes TEXT;

-- Create a table to track shipment status
CREATE TABLE IF NOT EXISTS
    shipments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
        status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'shipped', 'delivered', 'failed'
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        tracking_number VARCHAR(100),
        tracking_url TEXT,
        shipping_provider VARCHAR(100),
        shipping_notes TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create table for shipment items
CREATE TABLE IF NOT EXISTS
    shipment_items (
        shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
        order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
        quantity INT NOT NULL CHECK (quantity > 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        PRIMARY KEY(shipment_id, order_item_id)
    );


CREATE INDEX idx_shipping_methods_is_active ON shipping_methods(is_active);
CREATE INDEX idx_shipping_zones_is_active ON shipping_zones(is_active);
CREATE INDEX idx_shipping_rates_method_zone ON shipping_rates(shipping_method_id, shipping_zone_id);
CREATE INDEX idx_shipping_rates_is_active ON shipping_rates(is_active);
CREATE INDEX idx_shipping_rate_conditions_rate_id ON shipping_rate_conditions(shipping_rate_id);
CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipment_items_shipment_id ON shipment_items(shipment_id);
CREATE INDEX idx_shipment_items_order_item_id ON shipment_items(order_item_id);