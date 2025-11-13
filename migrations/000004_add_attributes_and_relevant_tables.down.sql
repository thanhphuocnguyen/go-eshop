-- Drop indexes for `attribute_values` table
DROP INDEX IF EXISTS idx_attribute_values_display_order;

DROP INDEX IF EXISTS idx_attribute_values_attribute_id;

-- Drop any remaining tables that might reference attribute_values
-- (in case previous migrations didn't run properly)
DROP TABLE IF EXISTS variant_attribute_values;

-- Drop attribute_values table (references attributes)
DROP TABLE IF EXISTS attribute_values;

-- Drop attributes table (parent table)
DROP TABLE IF EXISTS attributes;