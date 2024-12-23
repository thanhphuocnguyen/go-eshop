CREATE TYPE "user_role" AS ENUM ('admin', 'user');

CREATE TYPE "order_status" AS ENUM (
  'wait_for_confirming',
  'confirmed',
  'delivering',
  'delivered',
  'cancelled',
  'refunded',
  'completed'
);

CREATE TYPE "payment_status" AS ENUM ('not_paid', 'paid');

CREATE TYPE "payment_type" AS ENUM ('cash', 'transfer');

CREATE TYPE "card_type" AS ENUM ('debit', 'credit');

CREATE TABLE
  "users" (
    "id" bigserial PRIMARY KEY,
    "role" user_role NOT NULL DEFAULT 'user',
    "username" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "phone" varchar(12) NOT NULL,
    "full_name" varchar NOT NULL,
    "hashed_password" varchar NOT NULL,
    "verified_email" bool NOT NULL DEFAULT false,
    "verified_phone" bool NOT NULL DEFAULT false,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "products" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "description" text NOT NULL,
    "sku" varchar NOT NULL,
    "image_url" varchar,
    "stock" int NOT NULL,
    "archived" bool NOT NULL DEFAULT false,
    "price" numeric(8, 2) NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "carts" (
    "id" bigserial PRIMARY KEY,
    "checkout_at" timestamptz,
    "user_id" bigint NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now ()),
    UNIQUE ("id", "checkout_at")
  );

CREATE TABLE
  "cart_items" (
    "id" bigserial PRIMARY KEY,
    "product_id" bigint NOT NULL,
    "cart_id" bigint NOT NULL,
    "quantity" smallint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "orders" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "user_address_id" bigint NOT NULL,
    "status" order_status NOT NULL DEFAULT 'wait_for_confirming',
    "shipping_id" bigint,
    "payment_type" payment_type NOT NULL,
    "payment_status" payment_status NOT NULL DEFAULT 'not_paid',
    "is_cod" bool NOT NULL DEFAULT false,
    "cart_id" bigint NOT NULL,
    "confirmed_at" timestamptz,
    "cancelled_at" timestamptz,
    "delivered_at" timestamptz,
    "refunded_at" timestamptz,
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "order_items" (
    "id" bigserial PRIMARY KEY,
    "product_id" bigint NOT NULL,
    "order_id" bigint NOT NULL,
    "quantity" int NOT NULL,
    "price" decimal NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "shippings" (
    "id" bigserial PRIMARY KEY,
    "vendor" varchar NOT NULL,
    "order_id" bigint NOT NULL,
    "fee" decimal NOT NULL,
    "phone" varchar NOT NULL,
    "estimated_days" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "attribute_values" (
    "id" bigserial PRIMARY KEY,
    "attribute_id" bigint NOT NULL,
    "value" varchar NOT NULL,
    "color" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "attributes" (
    "id" bigserial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now ()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
  );

CREATE TABLE
  "categories" (
    "id" bigserial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
    "image_url" varchar,
    "published" bool NOT NULL DEFAULT true,
    "created_at" timestamptz NOT NULL DEFAULT (now ()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
  );

CREATE TABLE
  "category_products" (
    "category_id" bigint NOT NULL,
    "product_id" bigint NOT NULL
  );

CREATE TABLE
  "payment_infos" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "card_number" varchar NOT NULL,
    "expired_date" timestamptz NOT NULL,
    "vcc_code" varchar NOT NULL,
    "card_type" card_type NOT NULL,
    "is_verified" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "sessions" (
    "id" uuid PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expired_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

CREATE TABLE
  "user_addresses" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "phone" varchar(12) NOT NULL,
    "address_1" varchar NOT NULL,
    "address_2" varchar,
    "ward" varchar(100),
    "district" varchar(100) NOT NULL,
    "city" varchar(100) NOT NULL,
    "is_primary" bool NOT NULL DEFAULT false,
    "is_deleted" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT now (),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "deleted_at" timestamptz
  );

CREATE INDEX "idx_user_addresses_user_id" ON "user_addresses" ("user_id");

CREATE INDEX ON "products" ("price");

CREATE INDEX ON "products" ("archived");

CREATE INDEX ON "cart_items" ("product_id", "cart_id");

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "orders" ("shipping_id");

CREATE INDEX ON "orders" ("user_id");

CREATE INDEX ON "orders" ("user_id", "status");

CREATE INDEX ON "order_items" ("product_id", "order_id");

CREATE INDEX ON "shippings" ("order_id");

CREATE INDEX ON "category_products" ("category_id", "product_id");

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "cart_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "cart_items" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("shipping_id") REFERENCES "shippings" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "shippings" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "attribute_values" ADD FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("id");

ALTER TABLE "category_products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "category_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "payment_infos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_address_id") REFERENCES "user_addresses" ("id");

ALTER TABLE "user_addresses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");