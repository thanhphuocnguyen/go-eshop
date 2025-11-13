-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create user_role enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE "user_role" AS ENUM ('admin', 'user', 'moderator');
    END IF;
END$$;

-- Create order_status enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE "order_status" AS ENUM (
          'pending',
          'confirmed',
          'delivering',
          'delivered',
          'cancelled',
          'refunded',
          'completed'
        );
    END IF;
END$$;

-- Create payment_status enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
        CREATE TYPE "payment_status" AS ENUM (
          'pending',
          'success',
          'failed',
          'cancelled',
          'refunded',
          'processing'
        );
    END IF;
END$$;

-- Create payment_method enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
        CREATE TYPE "payment_method" AS ENUM (
          'credit_card',
          'debit_card',
          'paypal',
          'stripe',
          'apple_pay',
          'bank_transfer',
          'cod'
        );
    END IF;
END$$;

-- Create cart_status enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cart_status') THEN
        CREATE TYPE "cart_status" AS ENUM ('active', 'checked_out');
    END IF;
END$$;

-- Create card_type enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'card_type') THEN
        CREATE TYPE "card_type" AS ENUM ('debit', 'credit');
    END IF;
END$$;

-- Create entity_type enum type if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entity_type') THEN
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
    END IF;
END$$;