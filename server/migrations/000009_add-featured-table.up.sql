CREATE TABLE
    featured_sections (
        featured_id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        description TEXT,
        sort_order SMALLINT NOT NULL CHECK (
            sort_order >= 0
            AND sort_order <= 32767
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    featured_products (
        id SERIAL PRIMARY KEY,
        featured_id INT REFERENCES featured_sections (featured_id) ON DELETE CASCADE,
        product_id UUID REFERENCES products (product_id) ON DELETE CASCADE,
        sort_order SMALLINT NOT NULL CHECK (
            sort_order >= 0
            AND sort_order <= 32767
        ),
        UNIQUE (featured_id, product_id)
    );

CREATE INDEX ON featured_products (featured_id);

CREATE INDEX ON featured_products (product_id);