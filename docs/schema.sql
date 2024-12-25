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

CREATE INDEX ON "shippings" ("order_id");

CREATE INDEX ON "payment_infos" ("user_id");

ALTER TABLE "shippings" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "payment_infos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");