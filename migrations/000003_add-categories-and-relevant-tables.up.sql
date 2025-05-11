CREATE TABLE
    "categories" (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "name" VARCHAR UNIQUE NOT NULL,
        "description" TEXT,
        "image_url" TEXT,
        "image_id" VARCHAR,
        "published" BOOLEAN NOT NULL DEFAULT TRUE,
        "remarkable" BOOLEAN DEFAULT FALSE,
        "slug" VARCHAR UNIQUE NOT NULL,
        "display_order" INT DEFAULT 0,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    collections (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "name" VARCHAR(255) UNIQUE NOT NULL,
        "image_url" TEXT,
        "image_id" VARCHAR,
        "description" TEXT,
        "slug" VARCHAR UNIQUE NOT NULL,
        "remarkable" BOOLEAN DEFAULT FALSE,
        "display_order" INT DEFAULT 0,
        "published" BOOLEAN NOT NULL DEFAULT TRUE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE TABLE
    brands (
        "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "name" VARCHAR(255) UNIQUE NOT NULL,
        "image_url" TEXT,
        "image_id" VARCHAR,
        "description" TEXT,
        "slug" VARCHAR UNIQUE NOT NULL,
        "remarkable" BOOLEAN DEFAULT FALSE,
        "display_order" INT DEFAULT 0,
        "published" BOOLEAN NOT NULL DEFAULT TRUE,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now ()
    );