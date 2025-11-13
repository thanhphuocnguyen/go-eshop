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
DROP TABLE IF EXISTS products; 