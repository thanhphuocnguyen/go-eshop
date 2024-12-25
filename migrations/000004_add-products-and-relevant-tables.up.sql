CREATE TABLE
    "products" (
        "id" bigserial PRIMARY KEY,
        "name" varchar NOT NULL,
        "description" text NOT NULL,
        "sku" varchar NOT NULL,
        "image_url" varchar,
        "stock" int NOT NULL,
        "archived" bool NOT NULL DEFAULT false,
        "price" DECIMAL(10, 2) NOT NULL,
        "updated_at" timestamptz NOT NULL DEFAULT now(),
        "created_at" timestamptz NOT NULL DEFAULT now()
    );

CREATE TABLE
    product_variants (
        variant_id bigserial PRIMARY KEY,
        product_id bigint NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        sku VARCHAR(100) UNIQUE NOT NULL, -- Unique stock-keeping unit
        price DECIMAL(10, 2) NOT NULL, -- Override price for the variant
        stock bigint DEFAULT 0, -- Inventory for this variant
        created_at TIMESTAMP DEFAULT now(),
        updated_at TIMESTAMP DEFAULT now()
    );

CREATE TABLE
    attributes (
        attribute_id bigserial PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL -- e.g., "Size", "Color"
    );

CREATE TABLE
    attribute_values (
        value_id bigserial PRIMARY KEY,
        attribute_id bigint NOT NULL REFERENCES attributes (attribute_id) ON DELETE CASCADE,
        value VARCHAR(100) NOT NULL -- e.g., "Small", "Red"
    );

CREATE TABLE
    variant_attributes (
        id bigserial PRIMARY KEY,
        variant_id bigint NOT NULL REFERENCES product_variants (variant_id) ON DELETE CASCADE,
        value_id bigint NOT NULL REFERENCES attribute_values (value_id) ON DELETE CASCADE,
        UNIQUE (variant_id, value_id) -- Ensure no duplicate value for a variant
    );

CREATE INDEX ON "products" ("price");

CREATE INDEX ON "products" ("archived");