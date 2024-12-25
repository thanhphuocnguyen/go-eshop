CREATE TABLE
    "orders" (
        "id" bigserial PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES "users" ("id"),
        "user_address_id" bigint NOT NULL REFERENCES "user_addresses" ("id"),
        "total_price" DECIMAL(10, 2) NOT NULL,
        "status" order_status NOT NULL DEFAULT 'pending',
        "confirmed_at" timestamptz,
        "delivered_at" timestamptz,
        "cancelled_at" timestamptz,
        "refunded_at" timestamptz,
        "updated_at" timestamptz NOT NULL DEFAULT now(),
        "created_at" timestamptz NOT NULL DEFAULT now()
    );

CREATE TABLE
    "order_items" (
        "id" bigserial PRIMARY KEY,
        "product_id" bigint NOT NULL REFERENCES "products" ("id"),
        "order_id" bigint NOT NULL REFERENCES "orders" ("id") ON DELETE CASCADE,
        "quantity" int NOT NULL,
        "price" DECIMAL(10, 2) NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT now()
    );

CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        method payment_method NOT NULL, -- e.g., Credit Card, PayPal
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        gateway payment_gateway, -- e.g., Stripe, PayPal
        transaction_id VARCHAR(100), -- From the payment gateway,]
        created_at TIMESTAMP DEFAULT now(),
        updated_at TIMESTAMP DEFAULT now()
    );

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "orders" ("user_id");

CREATE INDEX ON "orders" ("user_id", "status");

CREATE INDEX ON "order_items" ("product_id", "order_id");

CREATE INDEX ON "payments" ("order_id");