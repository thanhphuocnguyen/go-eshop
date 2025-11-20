-- name: CreateCategory :one
INSERT INTO categories (name, slug, description, image_url, image_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories WHERE slug = $1 LIMIT 1;

-- name: GetCategories :many
SELECT * FROM categories WHERE published = COALESCE(sqlc.narg('published'), true) ORDER BY display_order LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = COALESCE(sqlc.narg('name'), name), 
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    image_id = COALESCE(sqlc.narg('image_id'), image_id),
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE id = $1 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CountCategories :one
SELECT count(*) FROM categories;

-- name: SeedCategories :copyfrom
INSERT INTO categories (name, slug, description, image_url, image_id) VALUES ($1, $2, $3, $4, $5);

-- name: AddProductsToCategory :copyfrom
INSERT INTO category_products (category_id, product_id) VALUES ($1, $2);

-- name: RemoveProductsFromCategory :exec
DELETE FROM category_products WHERE product_id = $1;