CREATE TABLE
    payments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        method payment_method NOT NULL, -- e.g., card, bank_transfer
        gateway VARCHAR DEFAULT NULL, -- e.g., Stripe, PayPal
        refund_id VARCHAR, -- From the payment gateway, e.g., Stripe
        payment_intent_id VARCHAR, -- From the payment gateway, e.g., Stripe
        charge_id VARCHAR, -- From the payment gateway, e.g., Stripe
        error_code VARCHAR, -- From the payment gateway, e.g., Stripe
        error_message VARCHAR, -- From the payment gateway, e.g., Stripe
        created_at TIMESTAMPTZ DEFAULT NOW (),
        updated_at TIMESTAMPTZ DEFAULT NOW ()
    );

CREATE TABLE
    payment_transactions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        payment_id UUID NOT NULL REFERENCES payments (id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        gateway_transaction_id VARCHAR, -- From the payment gateway, e.g., Stripe
        gateway_response_code VARCHAR, -- From the payment gateway, e.g., Stripe
        gateway_response_message VARCHAR, -- From the payment gateway, e.g., Stripe
        transaction_date TIMESTAMPTZ DEFAULT NOW (),
        created_at TIMESTAMPTZ DEFAULT NOW ()
    );

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("status");

CREATE INDEX ON "payments" ("method");

CREATE INDEX ON "payments" ("payment_intent_id");

CREATE INDEX ON "payments" ("gateway_charge_id");

CREATE INDEX ON "payment_transactions" ("payment_id");

CREATE INDEX ON "payment_transactions" ("status");

CREATE INDEX ON "payment_transactions" ("gateway_transaction_id");