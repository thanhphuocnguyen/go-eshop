-- name: CreateProduct :one
INSERT INTO products 
    (id, name, description, base_price, base_sku, slug, brand_id, collection_id, category_id) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: CreateProductVariant :one
INSERT INTO product_variants
    (id, product_id, sku, price, stock, weight)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateBulkProductVariants :copyfrom
INSERT INTO product_variants 
    (id, product_id, sku, price, stock, weight) 
VALUES 
    ($1, $2, $3, $4, $5, $6);

-- name: GetProductByID :one
SELECT
    products.*
FROM
    products
WHERE
    products.id = $1 AND
    is_active = COALESCE(sqlc.narg('is_active'), TRUE)
GROUP BY
    products.id;

-- name: GetProductVariantByID :one
SELECT
    product_variants.*
FROM
    product_variants
WHERE
    id = $1;

-- name: GetProductDetail :many
SELECT
    p.id as product_id, p.name, p.description, p.base_price, p.base_sku, p.slug, p.updated_at, p.created_at, p.is_active,
    c.id AS category_id, c.name AS category_name,
    cl.id AS collection_id, cl.name AS collection_name,
    b.id AS brand_id, b.name AS brand_name
    -- img.id AS img_id, img.url AS img_url, img.alt_text AS img_alt, 
    -- img.caption AS img_cap, img.mime_type AS img_mime_type, img.file_size AS image_size, 
    -- img.width AS img_w, img.height AS img_h, img.external_id AS img_external_id,
    -- ia.display_order AS img_assignment_display_order, ia.role AS img_assignment_role
FROM
    products p
JOIN categories as c ON p.category_id = c.id
JOIN brands AS b ON p.brand_id = b.id
LEFT JOIN collections as cl ON p.collection_id = cl.id
WHERE
    p.id = $1 AND
    p.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
ORDER BY
    p.id;

-- name: GetProductVariants :many
SELECT
    v.*,
    a.id as attr_id, a.name as attr_name,
    av.id as attr_val_id, av.value as attr_value, av.display_order as attr_display_order, 
    av.is_active as attr_val_is_active, av.display_value as attr_display_value
    -- img.id AS img_id, img.url AS img_url, img.alt_text AS img_alt, 
    -- img.caption AS img_cap, img.mime_type AS img_mime_type, img.file_size AS image_size, 
    -- img.width AS img_w, img.height AS img_h, img.external_id AS img_external_id,
    -- ia.display_order AS img_assignment_display_order, ia.role AS img_assignment_role
FROM
    product_variants AS v
JOIN variant_attribute_values as vav ON v.id = vav.variant_id
JOIN attribute_values as av ON vav.attribute_value_id = av.id
JOIN attributes as a ON av.attribute_id = a.id
WHERE
    v.product_id = $1 AND
    v.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
ORDER BY
    a.id, av.display_order, v.created_at DESC;

-- name: GetProducts :many
SELECT
    p.*,
    first_img.id AS img_id, first_img.url AS img_url,
    COUNT(v.id) AS variant_count
FROM products as p
LEFT JOIN product_variants as v ON p.id = v.product_id
LEFT JOIN LATERAL (
    SELECT img.id, img.url
    FROM image_assignments as ia
    JOIN images as img ON img.id = ia.image_id
    WHERE ia.entity_id = p.id AND ia.entity_type = 'product'
    ORDER BY ia.display_order ASC, ia.id ASC
    LIMIT 1
) AS first_img ON true
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), p.name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), p.base_sku)
GROUP BY
    p.id, first_img.id, first_img.url
ORDER BY
    p.id
LIMIT $1 OFFSET $2;

-- name: GetProductWithImage :one
SELECT
    products.*,
    img.id AS img_id, img.url AS img_url
FROM
    products
LEFT JOIN image_assignments AS ia ON products.id = img.external_id AND img.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE
    products.id = $1 AND
    is_active = COALESCE(sqlc.narg('is_active'), false);

-- name: GetProductsByCategoryID :many
SELECT
    p.*,
    first_img.id AS img_id, first_img.url AS img_url
FROM
    products AS p
LEFT JOIN LATERAL (
    SELECT img.id, img.url
    FROM image_assignments as ia
    JOIN images as img ON img.id = ia.image_id
    WHERE ia.entity_id = p.id AND ia.entity_type = 'product'
    ORDER BY ia.display_order ASC, ia.id ASC
    LIMIT 1
) AS first_img ON true
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), base_sku) AND
    p.category_id = $1
ORDER BY
    p.id
LIMIT
    $2
OFFSET
    $3;

-- name: GetProductsByCollectionID :many
SELECT
    p.*,
    first_img.id AS img_id, first_img.url AS img_url
FROM products AS p
LEFT JOIN LATERAL (
    SELECT img.id, img.url
    FROM image_assignments as ia
    JOIN images as img ON img.id = ia.image_id
    WHERE ia.entity_id = p.id AND ia.entity_type = 'product'
    ORDER BY ia.display_order ASC, ia.id ASC
    LIMIT 1
) AS first_img ON true
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), base_sku) AND
    p.collection_id = $1
ORDER BY p.id
LIMIT $2
OFFSET $3;

-- name: GetProductsByBrandID :many
SELECT
    p.*,
    img.id AS img_id,
    img.url AS img_url,
    MIN(p.base_price)::DECIMAL AS min_price,
    MAX(p.base_price)::DECIMAL AS max_price,
    MAX(p.base_price)::SMALLINT AS discount
FROM products AS p
LEFT JOIN LATERAL (
    SELECT img.id, img.url
    FROM image_assignments as ia
    JOIN images as img ON img.id = ia.image_id
    WHERE ia.entity_id = p.id AND ia.entity_type = 'product'
    ORDER BY ia.display_order ASC, ia.id ASC
    LIMIT 1
) AS img ON true
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), base_sku) AND
    p.brand_id = $1
GROUP BY
    p.id, img.id
ORDER BY
    p.id
LIMIT
    $2
OFFSET
    $3;

-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    name ILIKE COALESCE(sqlc.narg('name'), name);

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    brand_id = coalesce(sqlc.narg('brand_id'), brand_id),
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
DELETE FROM
    products
WHERE
    id = $1;

-- name: ArchiveProduct :exec
UPDATE
    products
SET
    is_active = $1,
    updated_at = NOW()
WHERE
    id = $2;
    
-- name: UpdateProductVariant :one
UPDATE
    product_variants
SET
    sku = coalesce(sqlc.narg('sku'), sku),
    price = coalesce(sqlc.narg('price'), price),
    stock = coalesce(sqlc.narg('stock'), stock),
    weight = coalesce(sqlc.narg('weight'), weight),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProductVariant :exec
DELETE FROM
    product_variants
WHERE
    id = $1;

-- name: ArchiveProductVariant :exec
UPDATE
    product_variants
SET
    is_active = $1,
    updated_at = NOW()
WHERE
    id = $2;

-- name: AddBulkProducts :copyfrom
INSERT INTO products (id, category_id, collection_id, brand_id, name, description) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateProductStock :one
UPDATE
    product_variants
SET
    stock = stock - $1
WHERE
    id = $2
RETURNING *;