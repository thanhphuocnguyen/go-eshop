CREATE TABLE
    "categories" (
        "category_id" SERIAL PRIMARY KEY,
        "name" VARCHAR UNIQUE NOT NULL,
        "image_url" TEXT,
        "sort_order" SMALLINT NOT NULL CHECK (
            "sort_order" >= 0
            AND "sort_order" <= 32767
        ),
        "published" BOOL NOT NULL DEFAULT TRUE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (category_id, sort_order)
    );

CREATE TABLE
    collections (
        collection_id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    brands (
        brand_id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        image_url TEXT,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );