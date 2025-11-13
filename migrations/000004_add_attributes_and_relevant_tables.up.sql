-- Create attributes table
CREATE TABLE
    attributes (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        "name" VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Create attribute_values table
CREATE TABLE
    attribute_values (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        attribute_id UUID NOT NULL REFERENCES attributes (id) ON DELETE CASCADE,
        "name" VARCHAR NOT NULL,
        code VARCHAR NOT NULL,
        is_active BOOLEAN DEFAULT TRUE,
        display_order SMALLINT NOT NULL DEFAULT 0 CHECK (
            display_order >= 0
            AND display_order <= 32767
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (attribute_id, code),
        UNIQUE(attribute_id, "name")
    );

-- Indexes for `attribute_values` table
-- Composite index to efficiently query attribute_values by attribute_id and display_order
CREATE INDEX idx_attribute_values_attribute_id_display_order ON attribute_values (attribute_id, display_order);
-- Index to optimize queries filtering or ordering by display_order
