CREATE TABLE
    "categories" (
        "id" serial PRIMARY KEY,
        "name" varchar UNIQUE NOT NULL,
        "sort_order" smallint NOT NULL,
        "image_url" varchar,
        "published" bool NOT NULL DEFAULT true,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
    );

CREATE TABLE
    "category_products" (
        "category_id" int NOT NULL REFERENCES "categories" ("id") ON DELETE CASCADE,
        "product_id" bigint NOT NULL REFERENCES "products" ("id") ON DELETE CASCADE,
        PRIMARY KEY ("category_id", "product_id")
    );