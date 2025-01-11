CREATE TABLE
    products (
        product_id bigserial PRIMARY KEY,
        name varchar NOT NULL,
        description text NOT NULL,
        sku varchar UNIQUE, -- Ensure SKU is unique
        stock int NOT NULL,
        discount int NOT NULL DEFAULT 0 CHECK (
            discount >= 0
            AND discount <= 100
        ),
        archived bool NOT NULL DEFAULT FALSE,
        price DECIMAL(10, 2) NOT NULL,
        updated_at timestamptz NOT NULL DEFAULT now (),
        created_at timestamptz NOT NULL DEFAULT now ()
    );

CREATE TABLE
    product_variants (
        variant_id bigserial PRIMARY KEY,
        product_id bigint NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
        variant_name VARCHAR(100) NOT NULL, -- e.g., Small, Medium, Large
        variant_sku VARCHAR(100) UNIQUE, -- Unique stock-keeping unit
        variant_price DECIMAL(10, 2) NOT NULL, -- Override price for the variant
        variant_stock int NOT NULL DEFAULT 0, -- Inventory for this variant
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    attributes (
        attribute_id serial PRIMARY KEY,
        attribute_name VARCHAR(100) UNIQUE NOT NULL -- e.g., Size, Color
    );

CREATE TABLE
    attribute_values (
        attribute_value_id serial PRIMARY KEY,
        attribute_id int NOT NULL REFERENCES attributes (attribute_id) ON DELETE CASCADE,
        attribute_value VARCHAR(100) NOT NULL -- Ensure value is not null
    );

CREATE TABLE
    variant_attributes (
        variant_attribute_id serial PRIMARY KEY,
        variant_id bigint REFERENCES product_variants (variant_id) ON DELETE CASCADE,
        attribute_value_id int NOT NULL REFERENCES attribute_values (attribute_value_id) ON DELETE CASCADE,
        UNIQUE (variant_id, attribute_value_id) -- Ensure no duplicate value for a variant
    );

CREATE INDEX ON products (price);

CREATE INDEX ON products (archived);

CREATE INDEX ON product_variants (product_id);

CREATE INDEX ON product_variants (variant_sku);

CREATE INDEX ON attributes (attribute_name);

CREATE INDEX ON attribute_values (attribute_id);

CREATE INDEX ON attribute_values (attribute_value);

CREATE INDEX ON variant_attributes (variant_id);

CREATE INDEX ON variant_attributes (attribute_value_id);