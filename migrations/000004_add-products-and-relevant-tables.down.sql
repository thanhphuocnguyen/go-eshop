DROP INDEX IF EXISTS "products_price_idx";
DROP INDEX IF EXISTS "products_archived_idx";

DROP INDEX IF EXISTS "product_variants_product_id_idx";
DROP INDEX IF EXISTS "product_variants_sku_idx";

DROP INDEX IF EXISTS "attributes_name_idx";

DROP INDEX IF EXISTS "attribute_values_attribute_id_idx";
DROP INDEX IF EXISTS "attribute_values_value_idx";

DROP INDEX IF EXISTS "variant_attributes_variant_id_idx";
DROP INDEX IF EXISTS "variant_attributes_value_id_idx";

DROP TABLE IF EXISTS variant_attributes;
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attributes;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;