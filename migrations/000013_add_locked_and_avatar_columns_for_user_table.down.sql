-- Remove columns added in the up migration

-- First remove the foreign key constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_avatar_image_id;

-- Then drop the columns
ALTER TABLE users
DROP COLUMN IF EXISTS locked,
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS avatar_image_id;
