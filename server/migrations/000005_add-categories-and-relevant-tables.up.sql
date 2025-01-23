CREATE TABLE
    "categories" (
        "category_id" SERIAL PRIMARY KEY,
        "name" VARCHAR UNIQUE NOT NULL,
        "description" TEXT,
        "sort_order" SMALLINT NOT NULL CHECK (
            "sort_order" >= 0
            AND "sort_order" <= 32767
        ),
        "published" BOOL NOT NULL DEFAULT TRUE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (category_id, sort_order)
    );

CREATE TABLE
    "category_products" (
        "category_id" INT NOT NULL REFERENCES "categories" ("category_id") ON DELETE CASCADE,
        "product_id" BIGINT NOT NULL REFERENCES "products" ("product_id") ON DELETE CASCADE,
        "sort_order" SMALLINT NOT NULL CHECK (
            "sort_order" >= 0
            AND "sort_order" <= 32767
        ),
        PRIMARY KEY ("category_id", "product_id")
    );