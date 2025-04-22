-- Create attributes table
CREATE TABLE
    attributes (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Create attribute_values table
CREATE TABLE
    attribute_values (
        id SERIAL PRIMARY KEY,
        attribute_id INT NOT NULL REFERENCES attributes (id) ON DELETE CASCADE,
        value VARCHAR NOT NULL UNIQUE,
        is_active BOOLEAN DEFAULT TRUE,
        display_value VARCHAR, -- TODO: change name to code
        display_order SMALLINT NOT NULL DEFAULT 0 CHECK (
            display_order >= 0
            AND display_order <= 32767
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Indexes for `attribute_values` table
CREATE INDEX idx_attribute_values_attribute_id ON attribute_values (attribute_id);

-- Crucial FK index
CREATE INDEX idx_attribute_values_display_order ON attribute_values (display_order);