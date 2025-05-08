-- Drop triggers first
DROP TRIGGER IF EXISTS after_rating_insert ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_update ON product_ratings;
DROP TRIGGER IF EXISTS after_rating_delete ON product_ratings;

-- Drop function
DROP FUNCTION IF EXISTS update_product_avg_rating();

-- Drop indexes
DROP INDEX IF EXISTS idx_product_ratings_product_id;
DROP INDEX IF EXISTS idx_product_ratings_user_id;
DROP INDEX IF EXISTS idx_product_ratings_rating;
DROP INDEX IF EXISTS idx_product_ratings_is_visible_is_approved;
DROP INDEX IF EXISTS idx_product_ratings_created_at;
DROP INDEX IF EXISTS idx_rating_votes_rating_id;
DROP INDEX IF EXISTS idx_rating_votes_user_id;
DROP INDEX IF EXISTS idx_rating_replies_rating_id;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS rating_replies;
DROP TABLE IF EXISTS rating_votes;
DROP TABLE IF EXISTS product_ratings;

-- Remove columns from products table
ALTER TABLE products
DROP COLUMN IF EXISTS avg_rating,
DROP COLUMN IF EXISTS rating_count,
DROP COLUMN IF EXISTS one_star_count,
DROP COLUMN IF EXISTS two_star_count,
DROP COLUMN IF EXISTS three_star_count,
DROP COLUMN IF EXISTS four_star_count,
DROP COLUMN IF EXISTS five_star_count;