-- Drop any remaining tables that might reference categories, brands, or collections
-- (in case previous migrations didn't run properly)
DROP TABLE IF EXISTS discount_categories;
DROP TABLE IF EXISTS products;

-- Drop main tables
DROP TABLE IF EXISTS "categories";

DROP TABLE IF EXISTS "brands";

DROP TABLE IF EXISTS "collections";