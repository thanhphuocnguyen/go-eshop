-- Merged E-Shop Database Rollback Migration File
-- This file combines all individual down migration files from 000011 to 000001 in reverse order
-- Created: November 13, 2025

-- Drop shipping tables and constraints (Migration 000011)
DROP INDEX IF EXISTS idx_shipment_items_order_item_id;
DROP INDEX IF EXISTS idx_shipment_items_shipment_id;
DROP INDEX IF EXISTS idx_shipments_status;
DROP INDEX IF EXISTS idx_shipments_order_id;
DROP INDEX IF EXISTS idx_shipping_rate_conditions_rate_id;
DROP INDEX IF EXISTS idx_shipping_rates_is_active;
DROP INDEX IF EXISTS idx_shipping_rates_method_zone;
DROP INDEX IF EXISTS idx_shipping_zones_is_active;
DROP INDEX IF EXISTS idx_shipping_methods_is_active;

-- Drop shipment_items table (must be dropped before shipments due to FK constraint)
DROP TABLE IF EXISTS shipment_items CASCADE;

-- Drop shipments table
DROP TABLE IF EXISTS shipments CASCADE;

-- Remove shipping-related fields from orders table
ALTER TABLE orders
DROP COLUMN IF EXISTS shipping_notes,
DROP COLUMN IF EXISTS shipping_provider,
DROP COLUMN IF EXISTS tracking_url,
DROP COLUMN IF EXISTS estimated_delivery_date,
DROP COLUMN IF EXISTS shipping_rate_id,
DROP COLUMN IF EXISTS shipping_method_id;

-- Drop shipping_rate_conditions table
DROP TABLE IF EXISTS shipping_rate_conditions;

-- Drop shipping_rates table
DROP TABLE IF EXISTS shipping_rates;

-- Drop shipping_zones table
DROP TABLE IF EXISTS shipping_zones;

-- Drop shipping_methods table
DROP TABLE IF EXISTS shipping_methods;

-- Drop ratings tables and related triggers/functions (Migration 000010)
-- Drop triggers first
DROP TRIGGER IF EXISTS after_rating_insert ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_update ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_delete ON product_ratings;
DROP TRIGGER IF EXISTS trigger_set_default_user_role ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_product_avg_rating();
DROP FUNCTION IF EXISTS set_default_user_role();

-- Drop indexes
DROP INDEX IF EXISTS idx_product_ratings_product_id;
DROP INDEX IF EXISTS idx_product_ratings_user_id;
DROP INDEX IF EXISTS idx_product_ratings_rating;
DROP INDEX IF EXISTS idx_product_ratings_is_visible_is_approved;
DROP INDEX IF EXISTS idx_product_ratings_created_at;
DROP INDEX IF EXISTS idx_rating_votes_rating_id;
DROP INDEX IF EXISTS idx_rating_votes_user_id;
DROP INDEX IF EXISTS idx_rating_replies_rating_id;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS rating_replies;
DROP TABLE IF EXISTS rating_votes;
DROP TABLE IF EXISTS product_ratings;

-- Remove rating columns from products table if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'products'
    ) THEN
        ALTER TABLE products
            DROP COLUMN IF EXISTS avg_rating,
            DROP COLUMN IF EXISTS rating_count,
            DROP COLUMN IF EXISTS one_star_count,
            DROP COLUMN IF EXISTS two_star_count,
            DROP COLUMN IF EXISTS three_star_count,
            DROP COLUMN IF EXISTS four_star_count,
            DROP COLUMN IF EXISTS five_star_count;
    END IF;
END
$$;

-- Drop payment tables (Migration 000009)
DROP INDEX IF EXISTS "payments_order_id_idx";
DROP INDEX IF EXISTS "payments_status_idx";
DROP INDEX IF EXISTS "payments_payment_method_id_idx";
DROP INDEX IF EXISTS "payments_payment_intent_id_idx";
DROP INDEX IF EXISTS "payments_charge_id_idx";
DROP INDEX IF EXISTS "payments_gateway_reference_idx";
DROP INDEX IF EXISTS "payment_transactions_payment_id_idx";
DROP INDEX IF EXISTS "payment_transactions_status_idx";
DROP INDEX IF EXISTS "payment_transactions_gateway_transaction_id_idx";

-- Drop order_items table here (references orders and product_variants)
DROP INDEX IF EXISTS idx_order_items_variant_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP TABLE IF EXISTS order_items CASCADE;

-- Drop order_discounts table (references orders and discounts)
DROP INDEX IF EXISTS idx_order_discounts_discount_id;
DROP INDEX IF EXISTS idx_order_discounts_order_id;
DROP TABLE IF EXISTS order_discounts CASCADE;

-- Drop payment tables
DROP TABLE IF EXISTS payment_transactions CASCADE;
DROP TABLE IF EXISTS payments CASCADE;

-- Drop cart tables (Migration 000008)
-- Drop indexes for `cart_items` table
DROP INDEX IF EXISTS idx_cart_items_variant_id;
DROP INDEX IF EXISTS idx_cart_items_cart_id;

-- Drop indexes for `carts` table  
DROP INDEX IF EXISTS idx_carts_updated_at;
DROP INDEX IF EXISTS idx_carts_session_id;
DROP INDEX IF EXISTS idx_carts_status; -- This is actually on order_id column
DROP INDEX IF EXISTS idx_carts_user_id;

-- Drop tables in proper dependency order
-- cart_items references carts, so drop it first
DROP TABLE IF EXISTS cart_items CASCADE;
DROP TABLE IF EXISTS carts CASCADE;

-- Drop order tables (Migration 000007)
-- NOTE: order_items and other FK-dependent tables are dropped in other sections
-- We only need to drop the orders table here after all referencing tables are gone

-- Drop indexes for `orders` table
DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_customer_email;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_order_date;
DROP INDEX IF EXISTS idx_orders_customer_id;

-- Drop foreign key constraints from orders table first
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_shipping_method;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_shipping_rate;

-- Drop orders table (after all dependent tables have been dropped)
DROP TABLE IF EXISTS orders CASCADE;

-- Drop discount tables (Migration 000005) - MOVED HERE before products
-- Note: order_items and order_discounts tables have already been dropped in the payment section above

-- Drop indexes (order_discounts indexes already dropped in payment section)
DROP INDEX IF EXISTS idx_discount_users_user_id;
DROP INDEX IF EXISTS idx_discount_users_discount_id;
DROP INDEX IF EXISTS idx_discount_categories_category_id;
DROP INDEX IF EXISTS idx_discount_categories_discount_id;
DROP INDEX IF EXISTS idx_discount_products_product_id;
DROP INDEX IF EXISTS idx_discount_products_discount_id;
DROP INDEX IF EXISTS idx_discounts_active_dates;
DROP INDEX IF EXISTS idx_discounts_dates;
DROP INDEX IF EXISTS idx_discounts_code;

-- Drop trigger and function
DROP TRIGGER IF EXISTS update_discount_timestamp ON discounts;
DROP FUNCTION IF EXISTS update_discount_updated_at();

-- Drop tables (order_discounts already dropped in payment section)
-- discount_products must be dropped before products table
DROP TABLE IF EXISTS discount_users CASCADE;
DROP TABLE IF EXISTS discount_categories CASCADE;
DROP TABLE IF EXISTS discount_products CASCADE;
DROP TABLE IF EXISTS discounts CASCADE;

-- Drop products and related tables (Migration 000006)
-- Drop images table first as it's not referenced by others
DROP INDEX IF EXISTS idx_product_images_product_id;
DROP INDEX IF EXISTS idx_product_images_image_id;
DROP INDEX IF EXISTS idx_product_images_is_primary;
DROP INDEX IF EXISTS idx_product_images_display_order;
DROP TABLE IF EXISTS product_images;

-- Drop indexes for `variant_attribute_values` table
DROP INDEX IF EXISTS idx_variant_attribute_values_attribute_value_id;

-- Drop indexes for `product_variants` table
DROP INDEX IF EXISTS idx_product_variants_is_active_stock;
DROP INDEX IF EXISTS idx_product_variants_price;
DROP INDEX IF EXISTS idx_product_variants_product_id;

-- Drop indexes for `products` table
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_name;

-- Drop tables in proper dependency order
-- Tables that reference product_variants must be dropped first
-- variant_attribute_values references product_variants
DROP TABLE IF EXISTS variant_attribute_values;

-- featured_products references products  
DROP TABLE IF EXISTS featured_products;

-- featured_sections has no dependencies
DROP TABLE IF EXISTS featured_sections;

-- product_variants references products, so drop it before products
DROP TABLE IF EXISTS product_variants;

-- products is referenced by other tables, so drop it last
DROP TABLE IF EXISTS products CASCADE;

-- Drop attributes tables (Migration 000004)
-- Drop indexes for `attribute_values` table
DROP INDEX IF EXISTS idx_attribute_values_attribute_id_display_order;
DROP INDEX IF EXISTS idx_attribute_values_display_order;
DROP INDEX IF EXISTS idx_attribute_values_attribute_id;

-- Drop attribute_values table (references attributes)
DROP TABLE IF EXISTS attribute_values;

-- Drop attributes table (parent table)
DROP TABLE IF EXISTS attributes;

-- Drop categories, collections, brands tables (Migration 000003)
-- Drop main tables
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "brands";
DROP TABLE IF EXISTS "collections";

-- Drop users and related tables (Migration 000002)
-- Drop indexes first
DROP INDEX IF EXISTS idx_email_verifications_expired_at;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS "user_addresses_user_id_default_idx";
DROP INDEX IF EXISTS "sessions_user_id_idx";
DROP INDEX IF EXISTS "idx_user_addresses_user_id";
DROP INDEX IF EXISTS "user_payment_infos_user_id_idx";

-- Drop tables in proper dependency order
-- Tables that reference users must be dropped first
DROP TABLE IF EXISTS "user_addresses";
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS "user_sessions";
DROP TABLE IF EXISTS user_payment_infos;

-- Drop users table last (parent table)
DROP TABLE IF EXISTS "users";

-- Drop reference tables and enum types (Migration 000001)
-- Drop indexes for reference tables
DROP INDEX IF EXISTS idx_card_types_is_active;
DROP INDEX IF EXISTS idx_card_types_code;
DROP INDEX IF EXISTS idx_payment_methods_display_order;
DROP INDEX IF EXISTS idx_payment_methods_gateway;
DROP INDEX IF EXISTS idx_payment_methods_code;
DROP INDEX IF EXISTS idx_payment_methods_is_active;
DROP INDEX IF EXISTS idx_user_roles_code;
DROP INDEX IF EXISTS idx_user_roles_is_active;

-- Drop reference tables
-- Drop permissions table first (it references user_roles)
DROP INDEX IF EXISTS idx_permissions_role_module;
DROP INDEX IF EXISTS idx_permissions_module;
DROP INDEX IF EXISTS idx_permissions_role_id;
DROP TABLE IF EXISTS permissions;

-- Drop other reference tables
DROP TABLE IF EXISTS card_types;
DROP TABLE IF EXISTS payment_methods;

-- Drop user_roles table last (after permissions)
DROP TABLE IF EXISTS user_roles;

-- Drop enum types
DROP TYPE IF EXISTS "cart_status";
DROP TYPE IF EXISTS "payment_status";
DROP TYPE IF EXISTS "order_status";

-- Drop UUID extension (optional - only if you want to remove it completely)
-- DROP EXTENSION IF EXISTS "uuid-ossp";