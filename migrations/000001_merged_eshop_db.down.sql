-- Down migration for E-Shop Database
-- This file reverses all changes from 000001_merged_eshop_db.up.sql
-- Created: November 16, 2025

-- Drop all indexes first
DROP INDEX IF EXISTS idx_shipment_items_order_item_id;
DROP INDEX IF EXISTS idx_shipment_items_shipment_id;
DROP INDEX IF EXISTS idx_shipments_order_id;
DROP INDEX IF EXISTS idx_shipments_status;
DROP INDEX IF EXISTS idx_shipping_rate_conditions_shipping_rate_id;
DROP INDEX IF EXISTS idx_shipping_rates_shipping_method_id;
DROP INDEX IF EXISTS idx_shipping_rates_shipping_zone_id;
DROP INDEX IF EXISTS idx_shipping_zones_countries;
DROP INDEX IF EXISTS idx_shipping_methods_is_active;
DROP INDEX IF EXISTS idx_rating_replies_rating_id;
DROP INDEX IF EXISTS idx_rating_votes_rating_id;
DROP INDEX IF EXISTS idx_rating_votes_user_id;
DROP INDEX IF EXISTS idx_product_ratings_product_id;
DROP INDEX IF EXISTS idx_product_ratings_user_id;
DROP INDEX IF EXISTS idx_product_ratings_verified_purchase;
DROP INDEX IF EXISTS idx_product_ratings_visible_approved;
DROP INDEX IF EXISTS idx_payment_transactions_payment_id;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_gateway_reference;
DROP INDEX IF EXISTS idx_cart_items_cart_id;
DROP INDEX IF EXISTS idx_cart_items_variant_id;
DROP INDEX IF EXISTS idx_carts_user_id;
DROP INDEX IF EXISTS idx_carts_session_id;
DROP INDEX IF EXISTS idx_order_discounts_order_id;
DROP INDEX IF EXISTS idx_order_discounts_discount_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_order_items_variant_id;
DROP INDEX IF EXISTS idx_orders_customer_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_order_date;
DROP INDEX IF EXISTS idx_discount_users_discount_id;
DROP INDEX IF EXISTS idx_discount_users_user_id;
DROP INDEX IF EXISTS idx_discount_categories_discount_id;
DROP INDEX IF EXISTS idx_discount_categories_category_id;
DROP INDEX IF EXISTS idx_discount_products_discount_id;
DROP INDEX IF EXISTS idx_discount_products_product_id;
DROP INDEX IF EXISTS idx_discounts_code;
DROP INDEX IF EXISTS idx_discounts_dates;
DROP INDEX IF EXISTS idx_discounts_active_dates;
DROP INDEX IF EXISTS idx_featured_products_featured_id;
DROP INDEX IF EXISTS idx_featured_products_product_id;
DROP INDEX IF EXISTS idx_variant_attribute_values_attribute_value_id;
DROP INDEX IF EXISTS idx_product_variants_product_id;
DROP INDEX IF EXISTS idx_product_variants_price;
DROP INDEX IF EXISTS idx_product_variants_is_active_stock;
DROP INDEX IF EXISTS idx_product_images_product_id;
DROP INDEX IF EXISTS idx_product_images_image_id;
DROP INDEX IF EXISTS idx_product_images_is_primary;
DROP INDEX IF EXISTS idx_product_images_display_order;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_attribute_values_attribute_id;
DROP INDEX IF EXISTS idx_brands_published;
DROP INDEX IF EXISTS idx_brands_slug;
DROP INDEX IF EXISTS idx_collections_published;
DROP INDEX IF EXISTS idx_collections_slug;
DROP INDEX IF EXISTS idx_categories_published;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_user_addresses_user_id;
DROP INDEX IF EXISTS idx_email_verifications_expired_at;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_card_types_code;
DROP INDEX IF EXISTS idx_card_types_is_active;
DROP INDEX IF EXISTS idx_permissions_role_id;
DROP INDEX IF EXISTS idx_permissions_module;
DROP INDEX IF EXISTS idx_permissions_role_module;
DROP INDEX IF EXISTS idx_payment_methods_code;
DROP INDEX IF EXISTS idx_payment_methods_is_active;
DROP INDEX IF EXISTS idx_payment_methods_gateway;
DROP INDEX IF EXISTS idx_user_roles_code;
DROP INDEX IF EXISTS idx_user_roles_is_active;

-- Drop all triggers
DROP TRIGGER IF EXISTS after_rating_delete ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_update ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_insert ON product_ratings;
DROP TRIGGER IF EXISTS update_discount_timestamp ON discounts;
DROP TRIGGER IF EXISTS trigger_set_default_user_role ON users;

-- Drop all functions
DROP FUNCTION IF EXISTS update_product_avg_rating();
DROP FUNCTION IF EXISTS update_discount_updated_at();
DROP FUNCTION IF EXISTS set_default_user_role();

-- Drop all foreign key constraints that were added later
ALTER TABLE IF EXISTS orders DROP CONSTRAINT IF EXISTS fk_orders_shipping_method;
ALTER TABLE IF EXISTS orders DROP CONSTRAINT IF EXISTS fk_orders_shipping_rate;

-- Drop all tables in reverse order of creation (respecting foreign key dependencies)
DROP TABLE IF EXISTS shipment_items;
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS shipping_rate_conditions;
DROP TABLE IF EXISTS shipping_rates;
DROP TABLE IF EXISTS shipping_zones;
DROP TABLE IF EXISTS shipping_methods;
DROP TABLE IF EXISTS rating_replies;
DROP TABLE IF EXISTS rating_votes;
DROP TABLE IF EXISTS product_ratings;
DROP TABLE IF EXISTS payment_transactions;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS order_discounts;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS discount_users;
DROP TABLE IF EXISTS discount_categories;
DROP TABLE IF EXISTS discount_products;
DROP TABLE IF EXISTS discounts;
DROP TABLE IF EXISTS featured_products;
DROP TABLE IF EXISTS featured_sections;
DROP TABLE IF EXISTS variant_attribute_values;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS category_products;
DROP TABLE IF EXISTS collection_products;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attributes;
DROP TABLE IF EXISTS product_attributes;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS user_payment_infos;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS card_types;
DROP TABLE IF EXISTS payment_methods;

-- Drop all enum types
DROP TYPE IF EXISTS cart_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;

-- Drop the UUID extension (only if it's safe to do so)
-- Note: Commented out as other applications might be using it
-- DROP EXTENSION IF EXISTS "uuid-ossp";