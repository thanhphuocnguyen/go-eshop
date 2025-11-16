-- name: CreateProduct :one
INSERT INTO products (name, description, short_description, base_price, base_sku, slug, attributes, brand_id, collection_id, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: CreateProductVariant :one
INSERT INTO product_variants (product_id, description, sku, price, stock, weight) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: CreateBulkProductVariants :copyfrom
INSERT INTO product_variants (product_id, sku, price, stock, weight) VALUES ($1, $2, $3, $4, $5);

-- name: GetProductByID :one
SELECT products.* FROM products WHERE products.id = $1 AND is_active = COALESCE(sqlc.narg('is_active'), TRUE) GROUP BY products.id;

-- name: GetProductBySlug :one
SELECT products.* FROM products WHERE products.slug = $1 AND is_active = COALESCE(sqlc.narg('is_active'), TRUE);

-- name: GetProductVariantByID :one
SELECT product_variants.* FROM product_variants WHERE id = $1;

-- name: GetProductDetail :one
SELECT
    p.id as product_id, p.name, p.description, p.base_price,
    p.base_sku, p.slug, p.updated_at, p.created_at, p.is_active,
    p.short_description, p.attributes,
    p.rating_count, p.one_star_count, p.two_star_count,
    p.three_star_count, p.four_star_count, p.five_star_count,
    c.id AS category_id, c.name AS category_name,
    cl.id AS collection_id, cl.name AS collection_name,
    b.id AS brand_id, b.name AS brand_name,
    d.id as discount_id, MAX(d.discount_value) AS max_discount_value, d.discount_type
FROM products p
LEFT JOIN discount_products AS pd ON p.id = pd.product_id
LEFT JOIN discount_categories AS pc ON p.category_id = pc.category_id
LEFT JOIN categories as c ON p.category_id = c.id
LEFT JOIN brands AS b ON p.brand_id = b.id
LEFT JOIN collections as cl ON p.collection_id = cl.id
LEFT JOIN discounts AS d ON pd.discount_id = d.id OR pc.discount_id = d.id
WHERE (p.id = $1 OR p.slug = $2) AND p.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
GROUP BY p.id, c.id, cl.id, b.id, d.id
LIMIT 1;

-- name: GetProductVariants :many
SELECT
    v.id, v.product_id, v.description, v.sku, v.price, v.stock, v.weight, v.is_active, v.created_at, v.updated_at,
    a.id as attr_id, a.name as attr_name,
    av.id as attr_val_id,
    av.value as attr_value
FROM product_variants AS v
JOIN variant_attribute_values as vav ON v.id = vav.variant_id
JOIN attribute_values as av ON vav.attribute_value_id = av.id
JOIN attributes as a ON av.attribute_id = a.id
WHERE v.product_id = $1 AND v.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
ORDER BY a.id, v.created_at DESC;

-- name: GetFilterListForCollectionID :many
SELECT c.name as category_name, c.id as category_id, br.id as brand_id, br.name AS brand_name, p.attributes
FROM products p
LEFT JOIN categories c ON c.id = p.category_id
LEFT JOIN collections cl ON p.collection_id = cl.id
LEFT JOIN brands br ON p.brand_id = br.id
WHERE cl.id = $1
GROUP BY c.id, br.id, p.attributes
ORDER BY c.id;

-- name: GetProducts :many
SELECT
    p.*,
    first_img.id AS img_id, first_img.image_url AS img_url,
    COUNT(v.id) AS variant_count, MIN(v.price)::DECIMAL AS min_price, MAX(v.price)::DECIMAL AS max_price
FROM products as p
LEFT JOIN product_variants as v ON p.id = v.product_id
LEFT JOIN LATERAL (
    SELECT pi.id, pi.image_url
    FROM product_images as pi
    WHERE pi.product_id = p.id
    ORDER BY pi.id ASC
    LIMIT 1
) AS first_img ON true
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) 
    AND (p.name ILIKE COALESCE(sqlc.narg('search'), p.name) OR 
        p.base_sku ILIKE COALESCE(sqlc.narg('search'), p.base_sku) OR 
        p.description ILIKE COALESCE(sqlc.narg('search'), p.description))
    -- 
    AND (ARRAY_LENGTH(sqlc.arg(category_ids)::uuid[], 1) IS NULL OR p.category_id = ANY(sqlc.arg(category_ids)::uuid[]))
    AND p.collection_id = COALESCE(sqlc.narg('collection_id'), p.collection_id)
    AND p.brand_id = COALESCE(sqlc.narg('brand_id'), p.brand_id)
    AND p.slug ILIKE COALESCE(sqlc.narg('slug'), p.slug)
GROUP BY
    p.id, first_img.id, first_img.image_url
ORDER BY
    @orderBy::text
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    name ILIKE COALESCE(sqlc.narg('name'), name)
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
    attributes = coalesce(sqlc.narg('attributes'), attributes),
    collection_id = coalesce(sqlc.narg('collection_id'), collection_id),
    category_id = coalesce(sqlc.narg('category_id'), category_id),
    slug = coalesce(sqlc.narg('slug'), slug),
    base_price = coalesce(sqlc.narg('base_price'), base_price),
    base_sku = coalesce(sqlc.narg('base_sku'), base_sku),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ArchiveProduct :exec
UPDATE products SET is_active = $1, updated_at = NOW() WHERE id = $2;
    
-- name: UpdateProductVariant :one
UPDATE product_variants
SET
    sku = coalesce(sqlc.narg('sku'), sku),
    price = coalesce(sqlc.narg('price'), price),
    stock = coalesce(sqlc.narg('stock'), stock),
    weight = coalesce(sqlc.narg('weight'), weight),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteProductVariant :exec
DELETE FROM product_variants WHERE id = $1;

-- name: ArchiveProductVariant :exec
UPDATE product_variants SET is_active = $1, updated_at = NOW() WHERE id = $2;

-- name: AddBulkProducts :copyfrom
INSERT INTO products (category_id, collection_id, brand_id, name, description) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateProductStock :one
UPDATE product_variants SET stock = stock - $1 WHERE id = $2 RETURNING *;