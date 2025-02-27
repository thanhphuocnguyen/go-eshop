-- Create products table
CREATE TABLE
    products (
        product_id UUID PRIMARY KEY,
        category_id INT REFERENCES categories (category_id) ON DELETE SET NULL,
        collection_id INT REFERENCES collections (collection_id) ON DELETE SET NULL,
        brand_id INT REFERENCES brands (brand_id) ON DELETE SET NULL,
        name VARCHAR NOT NULL,
        description TEXT NOT NULL,
        archived BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Create product_variants table
CREATE TABLE
    product_variants (
        variant_id UUID PRIMARY KEY,
        product_id UUID NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
        price DECIMAL(10, 2) NOT NULL,
        discount SMALLINT NOT NULL DEFAULT 0 CHECK (
            discount >= 0
            AND discount <= 100
        ),
        stock_quantity INT NOT NULL DEFAULT 0,
        sku VARCHAR(100) UNIQUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Create attributes table
CREATE TABLE
    attributes (
        attribute_id SERIAL PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Create variant_attributes table
CREATE TABLE
    variant_attributes (
        variant_attribute_id SERIAL PRIMARY KEY,
        variant_id UUID NOT NULL REFERENCES product_variants (variant_id) ON DELETE CASCADE,
        attribute_id INT NOT NULL REFERENCES attributes (attribute_id) ON DELETE CASCADE,
        value VARCHAR NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (variant_id, attribute_id, value) -- Ensure no duplicate combination of variant_id, attribute_id, and value
    );

-- Create indexes
CREATE INDEX idx_products_archived ON products (archived);

CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);

CREATE INDEX idx_product_variants_sku ON product_variants (sku);

CREATE INDEX idx_attributes_name ON attributes (name);

CREATE INDEX idx_variant_attributes_variant_id ON variant_attributes (variant_id);

CREATE INDEX idx_variant_attributes_value ON variant_attributes (value);