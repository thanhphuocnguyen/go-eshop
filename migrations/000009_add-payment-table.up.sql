CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        method payment_method NOT NULL, -- e.g., Credit Card, PayPal
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        gateway payment_gateway, -- e.g., Stripe, PayPal
        transaction_id VARCHAR(100), -- From the payment gateway,]
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now ()
    );

CREATE INDEX ON "payments" ("order_id");