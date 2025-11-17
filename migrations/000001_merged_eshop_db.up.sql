-- Merged E-Shop Database Migration File
-- This file combines all individual migration files from 000001 to 000011
-- Created: November 13, 2025

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types first
DO $$
BEGIN
    -- Create order_status enum type if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE "order_status" AS ENUM (
          'pending',
          'confirmed',
          'delivering',
          'delivered',
          'cancelled',
          'refunded',
          'completed'
        );
    END IF;

    -- Create payment_status enum type if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
        CREATE TYPE "payment_status" AS ENUM (
          'pending',
          'success',
          'failed',
          'cancelled',
          'refunded',
          'processing'
        );
    END IF;

    -- Create cart_status enum type if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cart_status') THEN
        CREATE TYPE "cart_status" AS ENUM ('active', 'checked_out');
    END IF;
END$$;

-- Create reference tables
-- Create user_roles reference table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create permissions table with Linux-style permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES user_roles(id) ON DELETE CASCADE,
    module VARCHAR(100) NOT NULL,
    r BOOLEAN NOT NULL DEFAULT false, -- read permission
    w BOOLEAN NOT NULL DEFAULT false, -- write permission
    x BOOLEAN NOT NULL DEFAULT false, -- execute permission
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(role_id, module)
);

-- Insert seed data for user_roles
INSERT INTO user_roles (code, name, description) VALUES
    ('admin', 'Administrator', 'Full system access with all permissions'),
    ('user', 'User', 'Standard user with basic permissions'),
    ('moderator', 'Moderator', 'Content moderator with limited admin permissions')
ON CONFLICT (code) DO NOTHING;

-- Insert seed data for permissions
INSERT INTO permissions (role_id, module, r, w, x) 
SELECT 
    r.id,
    perm.module,
    perm.read_perm,
    perm.write_perm,
    perm.execute_perm
FROM user_roles r
JOIN (
    -- Admin permissions - full access to all modules
    VALUES 
        ('admin', 'products', true, true, true),
        ('admin', 'orders', true, true, true),
        ('admin', 'users', true, true, true),
        ('admin', 'categories', true, true, true),
        ('admin', 'brands', true, true, true),
        ('admin', 'collections', true, true, true),
        ('admin', 'discounts', true, true, true),
        ('admin', 'payments', true, true, true),
        ('admin', 'ratings', true, true, true),
        ('admin', 'shipping', true, true, true),
        
        -- User permissions - basic access
        ('user', 'products', true, false, false),
        ('user', 'orders', true, true, false),
        ('user', 'ratings', true, true, false),
        ('user', 'cart', true, true, false),
        
        -- Moderator permissions - content management
        ('moderator', 'products', true, true, false),
        ('moderator', 'categories', true, true, false),
        ('moderator', 'brands', true, true, false),
        ('moderator', 'collections', true, true, false),
        ('moderator', 'ratings', true, true, true),
        ('moderator', 'orders', true, false, false),
        ('moderator', 'users', true, false, false)
) AS perm(role_code, module, read_perm, write_perm, execute_perm)
ON r.code = perm.role_code
ON CONFLICT (role_id, module) DO NOTHING;

-- Create payment_methods reference table
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    gateway_supported VARCHAR(100), -- e.g., 'stripe', 'paypal', null for COD
    icon_url TEXT,
    requires_account BOOLEAN NOT NULL DEFAULT false, -- Whether this method requires user account setup
    min_amount DECIMAL(10, 2), -- Minimum transaction amount for this method
    max_amount DECIMAL(10, 2), -- Maximum transaction amount for this method
    processing_fee_percentage DECIMAL(5, 4), -- Processing fee as percentage
    processing_fee_fixed DECIMAL(10, 2), -- Fixed processing fee amount
    currency_supported TEXT[], -- Array of supported currencies
    countries_supported TEXT[], -- Array of supported country codes
    metadata JSONB, -- Additional gateway-specific configuration
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert seed data for payment_methods
INSERT INTO payment_methods (
    code, name, description, gateway_supported, requires_account, 
    min_amount, max_amount, processing_fee_percentage, processing_fee_fixed,
    currency_supported, countries_supported
) VALUES
    ('credit_card', 'Credit Card', 'Credit card payment method', 'stripe', false, 1.00, 10000.00, 0.0290, 0.30, ARRAY['USD', 'EUR', 'GBP'], ARRAY['US', 'CA', 'GB', 'EU']),
    ('debit_card', 'Debit Card', 'Debit card payment method', 'stripe', false, 1.00, 5000.00, 0.0290, 0.30, ARRAY['USD', 'EUR', 'GBP'], ARRAY['US', 'CA', 'GB', 'EU']),
    ('paypal', 'PayPal', 'PayPal payment method', 'paypal', false, 1.00, 10000.00, 0.0349, 0.49, ARRAY['USD', 'EUR', 'GBP', 'CAD'], ARRAY['US', 'CA', 'GB', 'EU', 'AU']),
    ('stripe', 'Stripe', 'Stripe payment method', 'stripe', false, 0.50, 999999.99, 0.0290, 0.30, ARRAY['USD', 'EUR', 'GBP'], ARRAY['US', 'CA', 'GB', 'EU']),
    ('apple_pay', 'Apple Pay', 'Apple Pay payment method', 'stripe', false, 1.00, 10000.00, 0.0290, 0.30, ARRAY['USD', 'EUR', 'GBP'], ARRAY['US', 'CA', 'GB', 'EU']),
    ('bank_transfer', 'Bank Transfer', 'Bank transfer payment method', null, false, 10.00, 50000.00, 0.0000, 0.00, ARRAY['USD', 'EUR'], ARRAY['US', 'EU']),
    ('cod', 'Cash on Delivery', 'Cash on delivery payment method', null, false, 5.00, 500.00, 0.0000, 2.00, ARRAY['USD'], ARRAY['US'])
ON CONFLICT (code) DO NOTHING;

-- Create card_types reference table
CREATE TABLE IF NOT EXISTS card_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert seed data for card_types
INSERT INTO card_types (code, name, description) VALUES
    ('debit', 'Debit Card', 'Standard debit card'),
    ('credit', 'Credit Card', 'Standard credit card'),
    ('prepaid', 'Prepaid Card', 'Prepaid card with preloaded funds'),
    ('gift_card', 'Gift Card', 'Gift card for purchases'),
    ('corporate', 'Corporate Card', 'Corporate/company card'),
    ('virtual', 'Virtual Card', 'Virtual card for online transactions'),
    ('business', 'Business Card', 'Business credit/debit card'),
    ('rewards', 'Rewards Card', 'Rewards/loyalty card')
ON CONFLICT (code) DO NOTHING;

-- Create users and addresses tables
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        role_id UUID NOT NULL REFERENCES user_roles (id) ON DELETE RESTRICT,
        username VARCHAR UNIQUE NOT NULL,
        email VARCHAR UNIQUE NOT NULL,
        phone_number VARCHAR(20) NOT NULL CHECK (
            char_length(phone_number) >= 10
            AND char_length(phone_number) <= 20
        ),
        first_name VARCHAR NOT NULL,
        last_name VARCHAR NOT NULL,
        avatar_url VARCHAR,
        avatar_image_id VARCHAR,
        hashed_password VARCHAR NOT NULL,
        verified_email bool NOT NULL DEFAULT FALSE,
        verified_phone bool NOT NULL DEFAULT FALSE,
        locked BOOLEAN NOT NULL DEFAULT FALSE,
        password_changed_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    user_addresses (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        phone_number VARCHAR(20) NOT NULL CHECK (
            char_length(phone_number) >= 10
            AND char_length(phone_number) <= 20
        ),
        street VARCHAR NOT NULL,
        ward VARCHAR(100),
        district VARCHAR(100) NOT NULL,
        city VARCHAR(100) NOT NULL,
        is_default BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    email_verifications (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        verify_code VARCHAR(255) NOT NULL,
        is_used BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        expired_at TIMESTAMPTZ NOT NULL DEFAULT (NOW () + interval '1 day')
    );

CREATE TABLE
    user_sessions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        refresh_token VARCHAR NOT NULL,
        user_agent VARCHAR(512) NOT NULL,
        client_ip INET NOT NULL,
        blocked boolean NOT NULL DEFAULT FALSE,
        expired_at TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    user_payment_infos (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        card_number VARCHAR(19) NOT NULL,
        card_last4 VARCHAR(4) NOT NULL,
        payment_method_token VARCHAR(255) NOT NULL,
        expiration_date DATE NOT NULL,
        billing_address TEXT NOT NULL,
        is_default BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT NOW (),
        updated_at TIMESTAMPTZ DEFAULT NOW (),
        UNIQUE (user_id, payment_method_token)
    );

-- Create categories, collections, and brands tables
CREATE TABLE
    categories (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) UNIQUE NOT NULL,
        description TEXT,
        image_url TEXT,
        image_id VARCHAR,
        published BOOLEAN NOT NULL DEFAULT TRUE,
        slug VARCHAR UNIQUE NOT NULL,
        display_order INT DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    collections (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        image_id VARCHAR,
        description TEXT,
        slug VARCHAR UNIQUE NOT NULL,
        display_order INT DEFAULT 0,
        published BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    brands (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        image_id VARCHAR,
        description TEXT,
        slug VARCHAR UNIQUE NOT NULL,
        display_order INT DEFAULT 0,
        published BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- Create attributes tables
CREATE TABLE
    attributes (
        id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        name VARCHAR(100) UNIQUE NOT NULL
    );

CREATE TABLE
    attribute_values (
        id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        attribute_id INT NOT NULL REFERENCES attributes (id) ON DELETE CASCADE,
        value VARCHAR(255) NOT NULL,
        UNIQUE(attribute_id, value)
    );

-- Create products tables
CREATE TABLE
    products (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        short_description VARCHAR(1000),
        attributes int[],
        base_price DECIMAL(10, 2),
        base_sku VARCHAR(100) UNIQUE NOT NULL,
        slug VARCHAR(255) UNIQUE NOT NULL,
        is_active BOOLEAN DEFAULT TRUE,
        category_id UUID REFERENCES categories (id) ON DELETE SET NULL,
        collection_id UUID REFERENCES collections (id) ON DELETE SET NULL,
        brand_id UUID REFERENCES brands (id) ON DELETE SET NULL,
        image_url TEXT,
        image_id VARCHAR(255),
        avg_rating DECIMAL(2, 1) DEFAULT NULL,
        rating_count INT NOT NULL DEFAULT 0,
        one_star_count INT NOT NULL DEFAULT 0,
        two_star_count INT NOT NULL DEFAULT 0,
        three_star_count INT NOT NULL DEFAULT 0,
        four_star_count INT NOT NULL DEFAULT 0,
        five_star_count INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create product_images table for storing product image metadata
CREATE TABLE
    product_images (
        id BIGINT PRIMARY KEY,
        product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        image_url VARCHAR(1024) NOT NULL,
        image_id VARCHAR(255) NOT NULL,
        alt_text VARCHAR(255),
        caption TEXT,
        display_order BIGINT NOT NULL DEFAULT 0,
        uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    product_variants (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        description VARCHAR(1000),
        sku VARCHAR(100) UNIQUE NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL DEFAULT 0,
        weight DECIMAL(10, 2),
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        image_url TEXT,
        image_id VARCHAR(255)
    );

CREATE TABLE
    variant_attribute_values (
        variant_id UUID NOT NULL REFERENCES product_variants (id) ON DELETE CASCADE,
        attribute_value_id BIGINT NOT NULL REFERENCES attribute_values (id) ON DELETE CASCADE,
        PRIMARY KEY (variant_id, attribute_value_id)
    );

CREATE TABLE
    featured_sections (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) UNIQUE NOT NULL,
        slug VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        image_id VARCHAR(255),
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    featured_products (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        featured_id UUID REFERENCES featured_sections (id) ON DELETE CASCADE,
        product_id UUID REFERENCES products (id) ON DELETE CASCADE,
        sort_order SMALLINT NOT NULL CHECK (
            sort_order >= 0
            AND sort_order <= 32767
        ),
        UNIQUE (featured_id, product_id)
    );

-- Create discounts tables
CREATE TABLE IF NOT EXISTS discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value NUMERIC(10, 2) NOT NULL CHECK (discount_value > 0),
    min_purchase_amount NUMERIC(10, 2) CHECK (min_purchase_amount IS NULL OR min_purchase_amount > 0),
    max_discount_amount NUMERIC(10, 2) CHECK (max_discount_amount IS NULL OR max_discount_amount > 0),
    usage_limit INTEGER CHECK (usage_limit IS NULL OR usage_limit > 0),
    used_count INTEGER NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CHECK (expires_at > starts_at)
);

CREATE TABLE IF NOT EXISTS discount_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, product_id)
);

CREATE TABLE IF NOT EXISTS discount_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, category_id)
);

CREATE TABLE IF NOT EXISTS discount_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (discount_id, user_id)
);

-- Create orders tables
CREATE TABLE
    orders (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        customer_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
        customer_email VARCHAR(255) NOT NULL,
        customer_name VARCHAR(255) NOT NULL,
        customer_phone VARCHAR(50) NOT NULL,
        shipping_address JSONB NOT NULL,
        total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
        status order_status NOT NULL DEFAULT 'pending',
        confirmed_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        cancelled_at TIMESTAMPTZ,
        shipping_method VARCHAR(100),
        refunded_at TIMESTAMPTZ,
        order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        shipping_method_id UUID,
        shipping_rate_id UUID,
        estimated_delivery_date TIMESTAMPTZ,
        tracking_url TEXT,
        shipping_provider VARCHAR(100),
        shipping_notes TEXT
    );

CREATE TABLE
    order_items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        variant_id UUID NOT NULL REFERENCES product_variants (id) ON DELETE RESTRICT,
        quantity SMALLINT NOT NULL DEFAULT 1 CHECK (quantity > 0),
        price_per_unit_snapshot DECIMAL(10, 2) NOT NULL CHECK (price_per_unit_snapshot >= 0),
        line_total_snapshot DECIMAL(10, 2) NOT NULL CHECK (line_total_snapshot >= 0),
        product_name_snapshot VARCHAR(255) NOT NULL,
        variant_sku_snapshot VARCHAR(100) NOT NULL,
        attributes_snapshot JSONB NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        discounted_price DECIMAL(10, 2)
    );

CREATE TABLE IF NOT EXISTS order_discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    discount_amount NUMERIC(10, 2) NOT NULL CHECK (discount_amount > 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (order_id, discount_id)
);

-- Create carts tables
CREATE TABLE
    carts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        session_id VARCHAR(255),
        order_id UUID REFERENCES orders (id) ON DELETE SET NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    cart_items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        cart_id UUID NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
        variant_id UUID NOT NULL REFERENCES product_variants (id) ON DELETE CASCADE,
        quantity smallint NOT NULL DEFAULT 1 CHECK (quantity > 0),
        added_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (cart_id, variant_id)
    );

-- Create payments tables
CREATE TABLE
    payments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        payment_method_id UUID NOT NULL REFERENCES payment_methods (id) ON DELETE RESTRICT,
        amount DECIMAL(10, 2) NOT NULL,
        processing_fee DECIMAL(10, 2) DEFAULT 0.00,
        net_amount DECIMAL(10, 2) NOT NULL, -- amount - processing_fee
        status payment_status NOT NULL DEFAULT 'pending',
        gateway VARCHAR(50) DEFAULT NULL,
        gateway_reference VARCHAR(255), -- Gateway-specific transaction reference
        refund_id VARCHAR(255),
        payment_intent_id VARCHAR,
        charge_id VARCHAR(255),
        error_code VARCHAR(100),
        error_message TEXT,
        metadata JSONB, -- Additional payment-specific data
        created_at TIMESTAMPTZ DEFAULT NOW (),
        updated_at TIMESTAMPTZ DEFAULT NOW ()
    );

CREATE TABLE
    payment_transactions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        payment_id UUID NOT NULL REFERENCES payments (id) ON DELETE CASCADE,
        amount DECIMAL(12, 2) NOT NULL,
        status payment_status NOT NULL DEFAULT 'pending',
        gateway_transaction_id VARCHAR(255),
        gateway_response_code VARCHAR(100),
        gateway_response_message TEXT,
        transaction_date TIMESTAMPTZ DEFAULT NOW (),
        created_at TIMESTAMPTZ DEFAULT NOW ()
    );

-- Create ratings tables
CREATE TABLE
    product_ratings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        order_item_id UUID REFERENCES order_items(id) ON DELETE SET NULL,
        rating DECIMAL(2, 1) NOT NULL CHECK (rating >= 1.0 AND rating <= 5.0),
        review_title VARCHAR(255),
        review_content TEXT,
        verified_purchase BOOLEAN NOT NULL DEFAULT FALSE,
        is_visible BOOLEAN NOT NULL DEFAULT TRUE,
        is_approved BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE (product_id, user_id)
    );

CREATE TABLE
    rating_votes (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        rating_id UUID NOT NULL REFERENCES product_ratings(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        is_helpful BOOLEAN NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE (rating_id, user_id)
    );

CREATE TABLE
    rating_replies (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        rating_id UUID NOT NULL REFERENCES product_ratings(id) ON DELETE CASCADE,
        reply_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        is_visible BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create shipping tables
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

CREATE TABLE IF NOT EXISTS
    shipping_rate_conditions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        shipping_rate_id UUID NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
        condition_type VARCHAR(50) NOT NULL,
        min_value DECIMAL(10, 2),
        max_value DECIMAL(10, 2),
        additional_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
        category_ids UUID[],
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    shipments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
        status VARCHAR(50) NOT NULL DEFAULT 'pending',
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        tracking_number VARCHAR(100),
        tracking_url TEXT,
        shipping_provider VARCHAR(100),
        shipping_notes TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    shipment_items (
        shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
        order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
        quantity INT NOT NULL CHECK (quantity > 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        PRIMARY KEY(shipment_id, order_item_id)
    );

-- Add foreign key constraints for orders table
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_shipping_method FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id),
    ADD CONSTRAINT fk_orders_shipping_rate FOREIGN KEY (shipping_rate_id) REFERENCES shipping_rates(id);

-- Create functions and triggers
-- Function to set default user role
CREATE OR REPLACE FUNCTION set_default_user_role()
RETURNS TRIGGER AS $$
BEGIN
    -- If role_id is not provided, set it to 'user' role
    IF NEW.role_id IS NULL THEN
        SELECT id INTO NEW.role_id FROM user_roles WHERE code = 'user' LIMIT 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_default_user_role
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION set_default_user_role();

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

-- Create function to update product average ratings
CREATE OR REPLACE FUNCTION update_product_avg_rating() 
RETURNS TRIGGER AS $$
BEGIN
    -- Update the rating counts and average for the affected product
    WITH rating_stats AS (
        SELECT 
            COUNT(*) AS total_count,
            AVG(rating) AS avg,
            COUNT(*) FILTER (WHERE ROUND(rating) = 1) AS one_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 2) AS two_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 3) AS three_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 4) AS four_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 5) AS five_star
        FROM product_ratings
        WHERE product_id = NEW.product_id AND is_visible = TRUE AND is_approved = TRUE
    )
    UPDATE products
    SET 
        rating_count = rating_stats.total_count,
        one_star_count = rating_stats.one_star,
        two_star_count = rating_stats.two_star,
        three_star_count = rating_stats.three_star,
        four_star_count = rating_stats.four_star,
        five_star_count = rating_stats.five_star,
        avg_rating = CASE WHEN rating_stats.total_count > 0 THEN rating_stats.avg ELSE NULL END
    FROM rating_stats
    WHERE id = NEW.product_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update product ratings
CREATE TRIGGER after_rating_insert
AFTER INSERT ON product_ratings
FOR EACH ROW
EXECUTE FUNCTION update_product_avg_rating();

CREATE TRIGGER after_rating_update
AFTER UPDATE ON product_ratings
FOR EACH ROW
WHEN (OLD.rating != NEW.rating OR OLD.is_visible != NEW.is_visible OR OLD.is_approved != NEW.is_approved)
EXECUTE FUNCTION update_product_avg_rating();

CREATE TRIGGER after_rating_delete
AFTER DELETE ON product_ratings
FOR EACH ROW
EXECUTE FUNCTION update_product_avg_rating();

-- Create all indexes for performance optimization
-- Reference tables indexes
CREATE INDEX IF NOT EXISTS idx_payment_methods_code ON payment_methods(code);
CREATE INDEX IF NOT EXISTS idx_payment_methods_is_active ON payment_methods(is_active);
CREATE INDEX IF NOT EXISTS idx_payment_methods_gateway ON payment_methods(gateway_supported);
CREATE INDEX IF NOT EXISTS idx_user_roles_code ON user_roles(code);
CREATE INDEX IF NOT EXISTS idx_user_roles_is_active ON user_roles(is_active);
CREATE INDEX IF NOT EXISTS idx_permissions_role_id ON permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_permissions_module ON permissions(module);
CREATE INDEX IF NOT EXISTS idx_permissions_role_module ON permissions(role_id, module);
CREATE INDEX IF NOT EXISTS idx_card_types_code ON card_types(code);
CREATE INDEX IF NOT EXISTS idx_card_types_is_active ON card_types(is_active);

-- Users and related tables indexes
CREATE INDEX idx_email_verifications_expired_at ON email_verifications (expired_at);
CREATE INDEX ON user_payment_infos (user_id);
CREATE INDEX ON user_sessions (user_id);
CREATE INDEX ON user_addresses (user_id, is_default);
CREATE INDEX idx_users_role_id ON users (role_id);

-- Product images indexes
CREATE INDEX idx_product_images_product_id ON product_images (product_id);
CREATE INDEX idx_product_images_image_id ON product_images (image_id);
CREATE INDEX idx_product_images_display_order ON product_images (display_order);

-- Products indexes
CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_products_is_active ON products (is_active);
CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);
CREATE INDEX idx_product_variants_price ON product_variants (price);
CREATE INDEX idx_product_variants_is_active_stock ON product_variants (is_active, stock);
CREATE INDEX idx_variant_attribute_values_attribute_value_id ON variant_attribute_values (attribute_value_id);

-- Discounts indexes
CREATE INDEX IF NOT EXISTS idx_discounts_code ON discounts(code);
CREATE INDEX IF NOT EXISTS idx_discounts_dates ON discounts(starts_at, expires_at);
CREATE INDEX IF NOT EXISTS idx_discounts_active_dates ON discounts(starts_at, expires_at) WHERE is_active = true AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_discount_products_discount_id ON discount_products(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_products_product_id ON discount_products(product_id);
CREATE INDEX IF NOT EXISTS idx_discount_categories_discount_id ON discount_categories(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_categories_category_id ON discount_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_discount_users_discount_id ON discount_users(discount_id);
CREATE INDEX IF NOT EXISTS idx_discount_users_user_id ON discount_users(user_id);
CREATE INDEX IF NOT EXISTS idx_order_discounts_order_id ON order_discounts(order_id);
CREATE INDEX IF NOT EXISTS idx_order_discounts_discount_id ON order_discounts(discount_id);

-- Orders indexes
CREATE INDEX idx_orders_customer_id ON orders (customer_id);
CREATE INDEX idx_orders_order_date ON orders (order_date);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_orders_customer_email ON orders (customer_email);
CREATE INDEX idx_orders_created_at ON orders (created_at);
CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_variant_id ON order_items (variant_id);

-- Carts indexes
CREATE INDEX idx_carts_user_id ON carts (user_id);
CREATE INDEX idx_carts_status ON carts (order_id);
CREATE INDEX idx_carts_session_id ON carts (session_id);
CREATE INDEX idx_carts_updated_at ON carts (updated_at);
CREATE INDEX idx_cart_items_variant_id ON cart_items (variant_id);

-- Payments indexes
CREATE INDEX ON payments (order_id);
CREATE INDEX ON payments (status);
CREATE INDEX ON payments (payment_method_id);
CREATE INDEX ON payments (payment_intent_id);
CREATE INDEX ON payments (charge_id);
CREATE INDEX ON payments (gateway_reference);
CREATE INDEX ON payment_transactions (payment_id);
CREATE INDEX ON payment_transactions (status);
CREATE INDEX ON payment_transactions (gateway_transaction_id);

-- Ratings indexes
CREATE INDEX idx_product_ratings_product_id ON product_ratings(product_id);
CREATE INDEX idx_product_ratings_user_id ON product_ratings(user_id);
CREATE INDEX idx_product_ratings_rating ON product_ratings(rating);
CREATE INDEX idx_product_ratings_is_visible_is_approved ON product_ratings(is_visible, is_approved);
CREATE INDEX idx_product_ratings_created_at ON product_ratings(created_at);
CREATE INDEX idx_rating_votes_rating_id ON rating_votes(rating_id);
CREATE INDEX idx_rating_votes_user_id ON rating_votes(user_id);
CREATE INDEX idx_rating_replies_rating_id ON rating_replies(rating_id);

-- Shipping indexes
CREATE INDEX idx_shipping_methods_is_active ON shipping_methods(is_active);
CREATE INDEX idx_shipping_zones_is_active ON shipping_zones(is_active);
CREATE INDEX idx_shipping_rates_method_zone ON shipping_rates(shipping_method_id, shipping_zone_id);
CREATE INDEX idx_shipping_rates_is_active ON shipping_rates(is_active);
CREATE INDEX idx_shipping_rate_conditions_rate_id ON shipping_rate_conditions(shipping_rate_id);
CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipment_items_shipment_id ON shipment_items(shipment_id);
CREATE INDEX idx_shipment_items_order_item_id ON shipment_items(order_item_id);