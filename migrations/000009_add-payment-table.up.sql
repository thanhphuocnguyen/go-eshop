CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL REFERENCES orders (order_id) ON DELETE CASCADE,
        amount DECIMAL(10, 2) NOT NULL,
        method payment_method NOT NULL, -- e.g., Credit Card, PayPal
        status payment_status NOT NULL DEFAULT 'pending', -- Pending, Success, Failed
        gateway payment_gateway, -- e.g., Stripe, PayPal
        transaction_id VARCHAR(100), -- From the payment gateway,]
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now ()
    );

CREATE TABLE
    user_payment_infos (
        payment_method_id SERIAL PRIMARY KEY,
        user_id BIGINT REFERENCES users (user_id),
        card_number VARCHAR(16) NOT NULL,
        cardholder_name VARCHAR(100) NOT NULL,
        expiration_date DATE NOT NULL,
        billing_address TEXT NOT NULL,
        is_default BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "user_payment_infos" ("user_id");