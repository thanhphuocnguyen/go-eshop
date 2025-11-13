-- Create product_ratings table
CREATE TABLE
    product_ratings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        order_item_id UUID REFERENCES order_items(id) ON DELETE SET NULL,
        rating DECIMAL(2, 1) NOT NULL CHECK (rating >= 1.0 AND rating <= 5.0),
        review_title VARCHAR(255),
        review_content TEXT,
        verified_purchase BOOLEAN NOT NULL DEFAULT FALSE,
        is_visible BOOLEAN NOT NULL DEFAULT TRUE,
        is_approved BOOLEAN NOT NULL DEFAULT FALSE,

        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

        UNIQUE (product_id, user_id)
    );

-- Create rating_votes table to track user votes on reviews
CREATE TABLE
    rating_votes (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        rating_id UUID NOT NULL REFERENCES product_ratings(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        is_helpful BOOLEAN NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE (rating_id, user_id)
    );

-- Create rating_replies table for store/admin responses to reviews
CREATE TABLE
    rating_replies (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        rating_id UUID NOT NULL REFERENCES product_ratings(id) ON DELETE CASCADE,
        reply_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        is_visible BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Add aggregate rating fields to products table
ALTER TABLE products 
ADD COLUMN avg_rating DECIMAL(2, 1) DEFAULT NULL,
ADD COLUMN rating_count INT NOT NULL DEFAULT 0,
ADD COLUMN one_star_count INT NOT NULL DEFAULT 0,
ADD COLUMN two_star_count INT NOT NULL DEFAULT 0,
ADD COLUMN three_star_count INT NOT NULL DEFAULT 0,
ADD COLUMN four_star_count INT NOT NULL DEFAULT 0,
ADD COLUMN five_star_count INT NOT NULL DEFAULT 0;

-- Create indexes for product_ratings table
CREATE INDEX idx_product_ratings_product_id ON product_ratings(product_id);
CREATE INDEX idx_product_ratings_user_id ON product_ratings(user_id);
CREATE INDEX idx_product_ratings_rating ON product_ratings(rating);
CREATE INDEX idx_product_ratings_is_visible_is_approved ON product_ratings(is_visible, is_approved);
CREATE INDEX idx_product_ratings_created_at ON product_ratings(created_at);

-- Create indexes for rating_votes table
CREATE INDEX idx_rating_votes_rating_id ON rating_votes(rating_id);
CREATE INDEX idx_rating_votes_user_id ON rating_votes(user_id);

-- Create indexes for rating_replies table
CREATE INDEX idx_rating_replies_rating_id ON rating_replies(rating_id);

-- Create function to update product average ratings
CREATE OR REPLACE FUNCTION update_product_avg_rating() 
RETURNS TRIGGER AS $$
BEGIN
    -- Update the rating counts and average for the affected product
    WITH rating_stats AS (
        SELECT 
            COUNT(*) AS total_count,

            COUNT(*) FILTER (WHERE ROUND(rating) = 1) AS one_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 2) AS two_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 3) AS three_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 4) AS four_star,
            COUNT(*) FILTER (WHERE ROUND(rating) = 5) AS five_star
        FROM product_ratings
        WHERE product_id = NEW.product_id AND is_visible = TRUE AND is_approved = TRUE
    )
    UPDATE products
    SET 
        rating_count = rating_stats.total_count,
        one_star_count = rating_stats.one_star,
        two_star_count = rating_stats.two_star,
        three_star_count = rating_stats.three_star,
        four_star_count = rating_stats.four_star,
        five_star_count = rating_stats.five_star,
        avg_rating = CASE WHEN rating_stats.total_count > 0 THEN rating_stats.avg ELSE NULL END
    FROM rating_stats
    WHERE id = NEW.product_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update product ratings
CREATE TRIGGER after_rating_insert
AFTER INSERT ON product_ratings
FOR EACH ROW
EXECUTE FUNCTION update_product_avg_rating();

CREATE TRIGGER after_rating_update
AFTER UPDATE ON product_ratings
FOR EACH ROW
WHEN (OLD.rating != NEW.rating OR OLD.is_visible != NEW.is_visible OR OLD.is_approved != NEW.is_approved)
EXECUTE FUNCTION update_product_avg_rating();

CREATE TRIGGER after_rating_delete
AFTER DELETE ON product_ratings
FOR EACH ROW
EXECUTE FUNCTION update_product_avg_rating();