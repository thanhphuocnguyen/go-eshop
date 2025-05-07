-- Create products table
CREATE TABLE
    products (
        id UUID PRIMARY KEY,
        name VARCHAR NOT NULL,
        description TEXT NOT NULL,
        short_description VARCHAR(1000),
        attributes INTEGER[],
        base_price DECIMAL(10, 2),
        base_sku VARCHAR(100) UNIQUE NOT NULL,
        slug VARCHAR UNIQUE NOT NULL,
        is_active BOOLEAN DEFAULT TRUE,
        category_id UUID REFERENCES categories (id) ON DELETE SET NULL,
        collection_id UUID REFERENCES collections (id) ON DELETE SET NULL,
        brand_id UUID REFERENCES brands (id) ON DELETE SET NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create Products variants table
CREATE TABLE
    product_variants (
        id UUID PRIMARY KEY,
        product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        description VARCHAR(1000),
        sku VARCHAR(100) UNIQUE NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL DEFAULT 0,
        weight DECIMAL(10, 2),
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create variant_attributes table
CREATE TABLE
    variant_attribute_values (
        variant_id UUID NOT NULL REFERENCES product_variants (id) ON DELETE CASCADE,
        attribute_value_id INT NOT NULL REFERENCES attribute_values (id) ON DELETE CASCADE,
        PRIMARY KEY (variant_id, attribute_value_id)
    );

-- Create product_attributes table
CREATE TABLE
    featured_sections (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        slug VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        image_id VARCHAR,
        description TEXT,
        remarkable BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    featured_products (
        id SERIAL PRIMARY KEY,
        featured_id INT REFERENCES featured_sections (id) ON DELETE CASCADE,
        product_id UUID REFERENCES products (id) ON DELETE CASCADE,
        sort_order SMALLINT NOT NULL CHECK (
            sort_order >= 0
            AND sort_order <= 32767
        ),
        UNIQUE (featured_id, product_id)
    );

-- Indexes for `products` table
CREATE INDEX idx_products_name ON products (name);

CREATE INDEX idx_products_is_active ON products (is_active);

-- base_sku is likely indexed by its UNIQUE constraint
-- The UNIQUE constraint on (attribute_id, value) likely provides a useful index too
-- Indexes for `product_variants` table
CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);

-- Crucial FK index
CREATE INDEX idx_product_variants_price ON product_variants (price);

CREATE INDEX idx_product_variants_is_active_stock ON product_variants (is_active, stock);

-- Composite for finding sellable items
-- sku is likely indexed by its UNIQUE constraint
-- Indexes for `variant_attribute_values` table
-- The composite PK (variant_id, attribute_value_id) is automatically indexed.
-- We need an index for lookups based *only* on attribute_value_id:
CREATE INDEX idx_variant_attribute_values_attribute_value_id ON variant_attribute_values (attribute_value_id);

-- An index on just variant_id might also be useful if querying attributes for a variant often,
-- though the composite PK might cover this depending on the DB. Add if needed:
-- CREATE INDEX idx_variant_attribute_values_variant_id ON variant_attribute_values (variant_id);