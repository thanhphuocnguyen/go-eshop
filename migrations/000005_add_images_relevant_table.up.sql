-- This migration script creates the images and image_assignments tables
-- for storing image metadata and their associations with various entities.
-- It also creates indexes to optimize queries on these tables.
-- The images table stores information about each image, including its URL,
CREATE TABLE
    images (
        id SERIAL PRIMARY KEY,
        external_id VARCHAR(255) UNIQUE NOT NULL,
        url VARCHAR(1024) NOT NULL,
        alt_text VARCHAR(255),
        caption TEXT,
        mime_type VARCHAR(50),
        file_size BIGINT,
        width INT,
        height INT,
        uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Create image_assignments table
CREATE TABLE
    image_assignments (
        id SERIAL PRIMARY KEY,
        image_id INT NOT NULL REFERENCES images (id) ON DELETE CASCADE,
        entity_id UUID NOT NULL,
        entity_type VARCHAR(50) NOT NULL,
        display_order SMALLINT NOT NULL DEFAULT 0 CHECK (
            display_order >= 0
            AND display_order <= 32767
        ),
        role VARCHAR(50) NOT NULL, -- e.g., 'product', 'category', etc.
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE (entity_id, entity_type, image_id)
    );

-- Create index on images table for faster lookups by URL
CREATE INDEX idx_images_url ON images (url);

-- Create index on image_assignments table for faster lookups by entity_id and entity_type
CREATE INDEX idx_image_assignments_entity ON image_assignments (entity_id, entity_type);

CREATE INDEX idx_image_assignments_image_id ON image_assignments (image_id);

CREATE INDEX idx_image_assignments_entity_order ON image_assignments (entity_id, entity_type, display_order);