-- Drop indexes for `image_assignments` table
DROP INDEX IF EXISTS idx_image_assignments_entity_order;

DROP INDEX IF EXISTS idx_image_assignments_image_id;

DROP INDEX IF EXISTS idx_image_assignments_entity;

-- Drop index for `images` table
DROP INDEX IF EXISTS idx_images_url;

-- Drop tables
DROP TABLE IF EXISTS image_assignments;

DROP TABLE IF EXISTS images;