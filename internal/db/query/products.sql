-- name: CreateProduct :one
INSERT INTO products (name, description, short_description, base_price, base_sku, slug, brand_id, collection_id, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: CreateBulkProductVariants :copyfrom
INSERT INTO product_variants (product_id, sku, price, stock, weight) VALUES ($1, $2, $3, $4, $5);

-- name: GetProductByID :one
SELECT products.* FROM products WHERE products.id = $1 AND is_active = COALESCE(sqlc.narg('is_active'), TRUE) GROUP BY products.id;

-- name: GetProductBySlug :one
SELECT products.* FROM products WHERE products.slug = $1 AND is_active = COALESCE(sqlc.narg('is_active'), TRUE);

-- name: GetProductVariantByID :one
SELECT * FROM product_variants WHERE id = $1 AND product_id = $2 AND is_active = COALESCE(sqlc.narg('is_active'), TRUE) LIMIT 1;

-- name: GetVariantDetailByID :many
SELECT product_variants.*, attribute_values.id as attribute_value_id, attribute_values.value as attribute_value FROM product_variants
LEFT JOIN variant_attribute_values ON product_variants.id = variant_attribute_values.variant_id
LEFT JOIN attribute_values ON variant_attribute_values.attribute_value_id = attribute_values.id
WHERE product_variants.id = $1 AND product_variants.product_id = $2 AND product_variants.is_active = COALESCE(sqlc.narg('is_active'), product_variants.is_active)
GROUP BY product_variants.id, attribute_values.id;

-- name: GetProductDetail :one
SELECT
    p.*,
    c.id AS category_id, c.name AS category_name,
    cl.id AS collection_id, cl.name AS collection_name,
    b.id AS brand_id, b.name AS brand_name
FROM products p
LEFT JOIN categories as c ON p.category_id = c.id
LEFT JOIN brands AS b ON p.brand_id = b.id
LEFT JOIN collections as cl ON p.collection_id = cl.id
WHERE (p.id = $1 OR p.slug = $2) AND p.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
GROUP BY p.id, c.id, cl.id, b.id
LIMIT 1;

-- name: GetProductVariantList :many
SELECT * FROM product_variants WHERE product_id = $1 AND is_active = COALESCE(sqlc.narg('is_active'), is_active) ORDER BY id, created_at DESC;

-- name: GetFilterListForCollectionID :many
SELECT c.name as category_name, c.id as category_id, br.id as brand_id, br.name AS brand_name
FROM products p
LEFT JOIN categories c ON c.id = p.category_id
LEFT JOIN collections cl ON p.collection_id = cl.id
LEFT JOIN brands br ON p.brand_id = br.id
WHERE cl.id = $1
GROUP BY c.id, br.id
ORDER BY c.id;

-- name: GetAdminProducts :many
SELECT
    p.*
FROM products as p
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) 
    AND (p.name ILIKE COALESCE(sqlc.narg('search'), p.name) OR p.base_sku ILIKE COALESCE(sqlc.narg('search'), p.base_sku) OR p.description ILIKE COALESCE(sqlc.narg('search'), p.description))
    AND (ARRAY_LENGTH(sqlc.arg('category_ids')::uuid[], 1) IS NULL OR p.category_id = ANY(sqlc.arg('category_ids')::uuid[]))
    AND (p.collection_id IS NULL OR p.collection_id = COALESCE(sqlc.narg('collection_id'), p.collection_id))
    AND p.brand_id = COALESCE(sqlc.narg('brand_id'), p.brand_id)
    AND p.slug ILIKE COALESCE(sqlc.narg('slug'), p.slug)
GROUP BY
    p.id
ORDER BY
    @orderBy::text
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
    AND name ILIKE COALESCE(sqlc.narg('name'), name)
    AND category_id = COALESCE(sqlc.narg('category_id'), category_id)
    AND collection_id = COALESCE(sqlc.narg('collection_id'), collection_id)
    AND brand_id = COALESCE(sqlc.narg('brand_id'), brand_id);

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    short_description = coalesce(sqlc.narg('short_description'), short_description),
    brand_id = coalesce(sqlc.narg('brand_id'), brand_id),
    collection_id = coalesce(sqlc.narg('collection_id'), collection_id),
    category_id = coalesce(sqlc.narg('category_id'), category_id),
    slug = coalesce(sqlc.narg('slug'), slug),
    base_price = coalesce(sqlc.narg('base_price'), base_price),
    base_sku = coalesce(sqlc.narg('base_sku'), base_sku),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    image_id = coalesce(sqlc.narg('image_id'), image_id),
    updated_at = NOW()
WHERE
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ArchiveProduct :exec
UPDATE products SET is_active = $1, updated_at = NOW() WHERE id = $2;
    
-- Product Variants --
-- name: CreateProductVariant :one
INSERT INTO product_variants (product_id, description, sku, price, stock, weight) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateProductVariant :one
UPDATE product_variants
SET
    sku = coalesce(sqlc.narg('sku'), sku),
    price = coalesce(sqlc.narg('price'), price),
    stock = coalesce(sqlc.narg('stock'), stock),
    description = coalesce(sqlc.narg('description'), description),
    weight = coalesce(sqlc.narg('weight'), weight),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    image_id = coalesce(sqlc.narg('image_id'), image_id),
    updated_at = NOW()
WHERE id = sqlc.arg('id') AND product_id = sqlc.arg('product_id') RETURNING *;

-- name: DeleteProductVariant :exec
DELETE FROM product_variants WHERE id = $1 AND product_id = $2;

-- name: ArchiveProductVariant :exec
UPDATE product_variants SET is_active = $1, updated_at = NOW() WHERE id = $2;

-- name: AddBulkProducts :copyfrom
INSERT INTO products (category_id, collection_id, brand_id, name, description) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateProductStock :one
UPDATE product_variants SET stock = stock - $1 WHERE id = $2 RETURNING *;