-- name: InsertProductRating :one
INSERT INTO product_ratings (product_id, user_id, order_item_id, rating, review_title, review_content, verified_purchase) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateProductRating :one
UPDATE product_ratings SET 
    rating = COALESCE(sqlc.narg(rating), rating), 
    review_title = COALESCE(sqlc.narg(review_title), review_title), 
    review_content = COALESCE(sqlc.narg(review_content), review_content), 
    is_visible = COALESCE(sqlc.narg(is_visible), is_visible), 
    is_approved = COALESCE(sqlc.narg(is_approved), is_approved),
    verified_purchase = COALESCE(sqlc.narg(verified_purchase), verified_purchase),
    helpful_votes = COALESCE(sqlc.narg(helpful_votes), helpful_votes),
    unhelpful_votes = COALESCE(sqlc.narg(unhelpful_votes), unhelpful_votes),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteProductRating :exec
DELETE FROM product_ratings WHERE id = $1;

-- name: GetProductRating :one
SELECT * FROM product_ratings WHERE id = $1;

-- name: GetProductRatings :many
SELECT 
    sqlc.embed(pr), 
    u.id AS user_id, u.fullname, u.email
FROM product_ratings AS pr
JOIN users AS u ON u.id = pr.user_id
WHERE pr.product_id = $1 AND pr.is_visible = TRUE
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProductRatingsCount :one
SELECT COUNT(*) FROM product_ratings WHERE product_id = $1 AND is_visible = TRUE;

-- name: GetProductRatingsByUserID :many
SELECT 
    sqlc.embed(pr), 
    u.id AS user_id, u.fullname, u.email
FROM product_ratings AS pr
JOIN users AS u ON u.id = pr.user_id
WHERE pr.user_id = $1 AND pr.is_visible = TRUE
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProductRatingsByOrderItemIDs :many
SELECT 
    sqlc.embed(pr), 
    u.id AS user_id, u.fullname, u.email
FROM product_ratings AS pr
JOIN users AS u ON u.id = pr.user_id
WHERE pr.order_item_id = ANY(sqlc.arg(ids)::uuid[]) AND pr.is_visible = TRUE
ORDER BY pr.created_at DESC;


-- name: CountProductRatings :one
SELECT COUNT(*) FROM product_ratings WHERE product_id = $1 AND is_visible = TRUE;


-- name: InsertRatingVotes :one
INSERT INTO rating_votes (rating_id, user_id, is_helpful) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateRatingVotes :one
UPDATE rating_votes SET 
    is_helpful = COALESCE(sqlc.narg(is_helpful), is_helpful)
WHERE id = $1 RETURNING *;
-- name: DeleteRatingVotes :exec
DELETE FROM rating_votes WHERE id = $1;
-- name: GetRatingVotes :one
SELECT * FROM rating_votes WHERE id = $1;
-- name: GetRatingVotesByRatingID :many
SELECT * FROM rating_votes WHERE rating_id = $1;
-- name: GetRatingVotesByUserID :many
SELECT * FROM rating_votes WHERE user_id = $1;

-- name: GetRatingVotesCount :one
SELECT COUNT(*) FROM rating_votes WHERE rating_id = $1 AND is_helpful = TRUE;

-- name: GetRatingVotesCountByUserID :one
SELECT COUNT(*) FROM rating_votes WHERE user_id = $1 AND is_helpful = TRUE;

-- name: InsertRatingReply :one
INSERT INTO rating_replies (rating_id, reply_by, content) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateRatingReplies :one
UPDATE rating_replies SET 
    content = COALESCE(sqlc.narg(reply_content), content)
WHERE id = $1 RETURNING *;

-- name: DeleteRatingReplies :exec
DELETE FROM rating_replies WHERE id = $1;

-- name: GetRatingReplies :one
SELECT * FROM rating_replies WHERE id = $1;

-- name: GetRatingRepliesByRatingID :many
SELECT 
    sqlc.embed(rr), 
    u.id AS reply_by, u.fullname, u.email
FROM rating_replies AS rr
JOIN users AS u ON u.id = rr.reply_by
WHERE rr.rating_id = $1;

-- name: GetRatingRepliesByUserID :many
SELECT 
    sqlc.embed(rr), 
    u.id AS reply_by, u.fullname, u.email
FROM rating_replies AS rr
JOIN users AS u ON u.id = rr.reply_by
WHERE rr.reply_by = $1;