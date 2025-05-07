-- Drop indexes for `attribute_values` table
DROP INDEX IF EXISTS idx_attribute_values_display_order;

DROP INDEX IF EXISTS idx_attribute_values_attribute_id;

DROP TABLE IF EXISTS attribute_values;

DROP TABLE IF EXISTS attributes;