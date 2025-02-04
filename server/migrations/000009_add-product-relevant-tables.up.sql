CREATE TABLE
    "brand" (
        "id" SERIAL PRIMARY KEY,
        "name" VARCHAR(100) NOT NULL
    );

CREATE TABLE
    "brand_product" (
        "brand_id" SERIAL REFERENCES brand (id) ON DELETE CASCADE,
        "product_id" UUID REFERENCES products (product_id) ON DELETE CASCADE,
        PRIMARY KEY ("brand_id", "product_id")
    );

CREATE TABLE
    "section" (
        "id" SERIAL PRIMARY KEY,
        "category_id" UUID REFERENCES category (category_id) ON DELETE CASCADE,
        "name" VARCHAR(100) NOT NULL
    );

CREATE TABLE
    "rating" (
        "id" SERIAL PRIMARY KEY,
        "product_id" UUID REFERENCES products (product_id) ON DELETE CASCADE,
        "rating" SMALLINT NOT NULL,
        "comment" TEXT,
        "created_at" TIMESTAMP NOT NULL DEFAULT NOW ()
    );

CREATE INDEX ON "brand_product" ("brand_id");

CREATE INDEX ON "brand_product" ("product_id");

CREATE INDEX ON "section" ("category_id");

CREATE INDEX ON "rating" ("product_id");