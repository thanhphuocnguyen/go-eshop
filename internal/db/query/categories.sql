-- name: CreateCategory :one
INSERT INTO categories (name, slug, description, remarkable, image_url, image_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetCategoryByID :one
SELECT c.* FROM categories c WHERE c.id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT c.* FROM categories c WHERE c.slug = $1 LIMIT 1;

-- name: GetCategoryProductsByID :many
SELECT sqlc.embed(c), p.id, p.name as product_name, p.description as product_description FROM categories c LEFT JOIN products p ON c.id = p.id LEFT JOIN product_images pi ON pi.product_id = p.id WHERE c.id = $1;

-- name: GetCategories :many
SELECT * FROM categories WHERE published = COALESCE(sqlc.narg('published'), true) AND remarkable = COALESCE(sqlc.narg('remarkable'), remarkable) ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = COALESCE(sqlc.narg('name'), name), 
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    image_id = COALESCE(sqlc.narg('image_id'), image_id),
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    remarkable = COALESCE(sqlc.narg('remarkable'), remarkable),
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE id = $1 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CountCategories :one
SELECT count(*) FROM categories;

-- name: SeedCategories :copyfrom
INSERT INTO categories (name, slug, description, image_url, image_id) VALUES ($1, $2, $3, $4, $5);
