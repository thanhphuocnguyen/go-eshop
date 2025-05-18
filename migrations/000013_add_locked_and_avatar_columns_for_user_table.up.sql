-- Add locked, avatar_url, and avatar_image_id columns to users table

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS locked BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS avatar_url TEXT DEFAULT NULL,
ADD COLUMN IF NOT EXISTS avatar_image_id UUID DEFAULT NULL;

-- Add foreign key constraint for avatar_image_id referencing the images table
ALTER TABLE users
ADD CONSTRAINT fk_users_avatar_image_id 
FOREIGN KEY (avatar_image_id) REFERENCES images(id) ON DELETE SET NULL;
