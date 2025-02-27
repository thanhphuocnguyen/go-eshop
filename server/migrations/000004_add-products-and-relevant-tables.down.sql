-- Drop indexes
DROP INDEX IF EXISTS idx_products_archived;

DROP INDEX IF EXISTS idx_product_variants_product_id;

DROP INDEX IF EXISTS idx_product_variants_sku;

DROP INDEX IF EXISTS idx_attributes_name;

DROP INDEX IF EXISTS idx_variant_attributes_variant_id;

DROP INDEX IF EXISTS idx_variant_attributes_value;

-- Drop tables
DROP TABLE IF EXISTS variant_attributes;

DROP TABLE IF EXISTS attributes;

DROP TABLE IF EXISTS product_variants;

DROP TABLE IF EXISTS products;