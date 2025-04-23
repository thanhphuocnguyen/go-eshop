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

CREATE TYPE "payment_status" AS ENUM (
  'pending',
  'success',
  'failed',
  'cancelled',
  'refunded',
  'processing'
);

CREATE TYPE "payment_method" AS ENUM (
  'card',
  'cod',
  'wallet',
  'postpaid',
  'bank_transfer'
);

CREATE TYPE "payment_gateway" AS ENUM (
  'stripe',
  'paypal',
  'visa',
  'mastercard',
  'apple_pay',
  'google_pay',
  'postpaid',
  'momo',
  'zalo_pay',
  'vn_pay'
);

CREATE TYPE "cart_status" AS ENUM ('active', 'checked_out');

CREATE TYPE "card_type" AS ENUM ('debit', 'credit');

CREATE TYPE "image_role" AS ENUM (
  'gallery',
  'thumbnail',
  'banner',
  'avatar',
  'cover',
  'logo',
  'icon',
  'background',
  'product',
  'category',
  'brand',
  'user',
  'order',
  'cart',
  'payment'
);

CREATE TYPE "entity_type" AS ENUM (
  'product',
  'product_variant',
  'category',
  'brand',
  'user',
  'order',
  'cart',
  'payment'
);