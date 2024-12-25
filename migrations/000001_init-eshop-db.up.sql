CREATE TYPE "user_role" AS ENUM ('admin', 'user', 'moderator');

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

CREATE TYPE "payment_method" AS ENUM (
  'credit_card',
  'paypal',
  'cod',
  'debit_card',
  'apple_pay',
  'wallet',
  'postpaid'
);

CREATE TYPE "payment_gateway" AS ENUM (
  'stripe',
  'paypal',
  'razorpay',
  'visa',
  'mastercard',
  'amex',
  'apple_pay',
  'google_pay',
  'amazon_pay',
  'phone_pe',
  'paytm',
  'upi',
  'wallet',
  'cod',
  'postpaid'
);

CREATE TYPE "cart_status" AS ENUM ('active', 'checked_out');

CREATE TYPE "card_type" AS ENUM ('debit', 'credit');