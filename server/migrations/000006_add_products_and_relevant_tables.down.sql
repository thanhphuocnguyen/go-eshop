-- Drop indexes for `variant_attribute_values` table
DROP INDEX IF EXISTS idx_variant_attribute_values_attribute_value_id;

-- Drop indexes for `product_variants` table
DROP INDEX IF EXISTS idx_product_variants_is_active_stock;

DROP INDEX IF EXISTS idx_product_variants_price;

DROP INDEX IF EXISTS idx_product_variants_product_id;

-- Drop indexes for `products` table
DROP INDEX IF EXISTS idx_products_is_active;

DROP INDEX IF EXISTS idx_products_name;

-- Drop tables
DROP TABLE IF EXISTS featured_products;

DROP TABLE IF EXISTS featured_sections;

DROP TABLE IF EXISTS variant_attribute_values;

DROP TABLE IF EXISTS product_variants;

DROP TABLE IF EXISTS products;