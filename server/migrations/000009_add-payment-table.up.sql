CREATE TABLE
    payments (
        id VARCHAR PRIMARY KEY,
        order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        payment_method payment_method NOT NULL, -- e.g., card, bank_transfer
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        payment_gateway payment_gateway DEFAULT NULL, -- e.g., Stripe, PayPal
        refund_id VARCHAR, -- From the payment gateway, e.g., Stripe
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

CREATE TABLE
    user_payment_infos (
        "id" SERIAL PRIMARY KEY,
        "user_id" UUID REFERENCES users (id) ON DELETE CASCADE,
        "card_number" VARCHAR(16) NOT NULL CHECK (card_number ~ '^[0-9]{16}$'),
        "cardholder_name" VARCHAR(100) NOT NULL,
        "expiration_date" DATE NOT NULL,
        "billing_address" TEXT NOT NULL,
        "default" BOOLEAN DEFAULT FALSE,
        "created_at" TIMESTAMPTZ DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ DEFAULT NOW(),
        UNIQUE ("user_id", "card_number")
    );

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "user_payment_infos" ("user_id");