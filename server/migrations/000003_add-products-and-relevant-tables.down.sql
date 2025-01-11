DROP INDEX IF EXISTS products_price_idx;
DROP INDEX IF EXISTS products_archived_idx;
DROP INDEX IF EXISTS product_variants_product_id_idx;
DROP INDEX IF EXISTS product_variants_variant_sku_idx;
DROP INDEX IF EXISTS attributes_attribute_name_idx;
DROP INDEX IF EXISTS attribute_values_attribute_id_idx;
DROP INDEX IF EXISTS attribute_values_attribute_value_idx;
DROP INDEX IF EXISTS variant_attributes_variant_id_idx;
DROP INDEX IF EXISTS variant_attributes_attribute_value_id_idx;

DROP TABLE IF EXISTS variant_attributes;
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attributes;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;