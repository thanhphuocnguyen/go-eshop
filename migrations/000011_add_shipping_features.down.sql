-- Drop indexes
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
DROP TABLE IF EXISTS shipment_items;

-- Drop shipments table
DROP TABLE IF EXISTS shipments;
-- Remove shipping-related fields from orders table
ALTER TABLE orders
DROP COLUMN IF EXISTS shipping_notes,
DROP COLUMN IF EXISTS shipping_provider,
DROP COLUMN IF EXISTS tracking_url,
DROP COLUMN IF EXISTS tracking_number,
DROP COLUMN IF EXISTS estimated_delivery_date,
DROP COLUMN IF EXISTS shipping_cost,
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