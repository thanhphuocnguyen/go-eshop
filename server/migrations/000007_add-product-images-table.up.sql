CREATE TABLE
    images (
        image_id SERIAL PRIMARY KEY,
        product_id bigint REFERENCES products (product_id) ON DELETE CASCADE,
        variant_id bigint REFERENCES product_variants (variant_id) ON DELETE CASCADE,
        image_url TEXT NOT NULL, -- URL/path of the image
        external_id TEXT, -- ID of the image on cloudinary or other services
        "primary" BOOLEAN DEFAULT FALSE, -- Indicates the main image
        created_at TIMESTAMP DEFAULT NOW (),
        updated_at TIMESTAMP DEFAULT NOW (),
        CHECK (
            product_id IS NOT NULL
            OR variant_id IS NOT NULL
        ) -- Ensure at least one association
    );

CREATE INDEX ON images (product_id);

CREATE INDEX ON images (variant_id);