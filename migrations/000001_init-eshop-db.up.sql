CREATE TYPE "user_role" AS ENUM ('admin', 'user');

CREATE TYPE "order_status" AS ENUM (
  'pending',
  'confirmed',
  'delivering',
  'delivered',
  'cancelled',
  'refunded',
  'completed'
);

CREATE TYPE "payment_status" AS ENUM ('pending', 'success', 'failed');

CREATE TYPE "payment_method" AS ENUM ('credit_card', 'paypal', 'cod');

CREATE TYPE "card_type" AS ENUM ('debit', 'credit');