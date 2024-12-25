CREATE TABLE
    "categories" (
        "id" serial PRIMARY KEY,
        "name" varchar UNIQUE NOT NULL,
        "sort_order" smallint NOT NULL UNIQUE CHECK ("sort_order" >= 0 AND "sort_order" <= 32767), 
        "image_url" varchar,
        "published" bool NOT NULL DEFAULT true,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        "updated_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "category_products" (
        "category_id" int NOT NULL REFERENCES "categories" ("id") ON DELETE CASCADE,
        "product_id" bigint NOT NULL REFERENCES "products" ("id") ON DELETE CASCADE,
        PRIMARY KEY ("category_id", "product_id")
    );