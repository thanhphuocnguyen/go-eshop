CREATE TABLE
    "categories" (
        "category_id" serial PRIMARY KEY,
        "name" varchar UNIQUE NOT NULL,
        "description" text,
        "sort_order" smallint NOT NULL UNIQUE CHECK (
            "sort_order" >= 0
            AND "sort_order" <= 32767
        ),
        "published" bool NOT NULL DEFAULT true,
        "created_at" timestamptz NOT NULL DEFAULT now (),
        "updated_at" timestamptz NOT NULL DEFAULT now ()
    );

CREATE TABLE
    "category_products" (
        "category_id" int NOT NULL REFERENCES "categories" ("category_id") ON DELETE CASCADE,
        "product_id" bigint NOT NULL REFERENCES "products" ("product_id") ON DELETE CASCADE,
        "sort_order" smallint NOT NULL UNIQUE CHECK (
            "sort_order" >= 0
            AND "sort_order" <= 32767
        ),
        PRIMARY KEY ("category_id", "product_id")
    );