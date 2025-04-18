-- name: CreateProduct :one
INSERT INTO products 
    (id, name, description, base_price, base_sku, slug, brand_id, collection_id, category_id) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: CreateProductVariant :one
INSERT INTO product_variants
    (id, product_id, sku, price, stock, weight, image_url, image_id)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CreateBulkProductVariants :copyfrom
INSERT INTO product_variants 
    (id, product_id, sku, price, stock, weight, image_url) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7);

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

-- name: GetProductVariants :many
SELECT
    p.id as product_id, p.name, p.description, p.base_price, p.base_sku, p.slug, p.updated_at, p.created_at, p.is_active,
    v.id as variant_id, v.sku as variant_sku, v.price as variant_price, v.stock as variant_stock, 
    v.weight as variant_weight, v.is_active as variant_is_active, v.image_url as variant_image_url, v.image_id as variant_image_id,
    a.id as attribute_id, a.name as attribute_name,
    av.id as attribute_value_id, av.value as attribute_value, av.display_order as attribute_display_order, 
    av.is_active as attribute_value_is_active, av.display_value as attribute_display_value,
    c.id AS category_id, c.name AS category_name,
    cl.id AS collection_id, cl.name AS collection_name,
    b.id AS brand_id, b.name AS brand_name,
    img.id AS image_id, img.url AS image_url, img.alt_text AS image_alt_text, 
    img.caption AS image_caption, img.mime_type AS image_mime_type, img.file_size AS image_file_size, 
    img.width AS image_width, img.height AS image_height, img.external_id AS image_external_id,
    ia.display_order AS image_display_order, ia.role AS image_role
FROM
    products p
LEFT JOIN product_variants as v ON p.id = v.product_id
LEFT JOIN variant_attribute_values as vav ON v.id = vav.variant_id
LEFT JOIN attribute_values as av ON vav.attribute_value_id = av.id
LEFT JOIN attributes as a ON av.attribute_id = a.id
LEFT JOIN categories as c ON p.category_id = c.id
LEFT JOIN collections as cl ON p.collection_id = cl.id
LEFT JOIN brands AS b ON p.brand_id = b.id
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE
    p.id = $1 AND
    p.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
ORDER BY
    p.id, v.id, a.id, av.display_order;

-- name: GetProducts :many
SELECT
    p.*,
    img.id AS image_id, img.url AS image_url,
    COUNT(v.id) AS variant_count
FROM products as p
LEFT JOIN product_variants as v ON p.id = v.product_id
LEFT JOIN (
    SELECT
        image_assignments.entity_id, image_assignments.image_id, image_assignments.display_order, image_assignments.role,
        images.id, images.url
    FROM image_assignments
    LEFT JOIN images ON images.id = image_assignments.image_id
    WHERE image_assignments.entity_type = 'product'
    LIMIT 1
) AS img ON p.id = img.entity_id
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), p.name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), p.base_sku)
GROUP BY
    p.id, img.id, img.url
ORDER BY
    p.id
LIMIT $1 OFFSET $2;

-- name: GetProductWithImage :one
SELECT
    products.*,
    img.id AS image_id, img.url AS image_url
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
    img.id AS image_id, img.url AS image_url
FROM
    products AS p
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), base_sku) AND
    p.category_id = $1
GROUP BY
    p.id, img.id
ORDER BY
    p.id
LIMIT
    $2
OFFSET
    $3;

-- name: GetProductsByCollectionID :many
SELECT
    p.*,
    img.id AS image_id, img.url AS image_url
FROM products AS p
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), is_active) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    p.base_sku ILIKE COALESCE(sqlc.narg('base_sku'), base_sku) AND
    p.collection_id = $1
GROUP BY p.id, img.id
ORDER BY p.id
LIMIT $2
OFFSET $3;

-- name: GetProductsByBrandID :many
SELECT
    p.*,
    img.id AS image_id,
    img.url AS image_url,
    MIN(p.base_price)::DECIMAL AS min_price,
    MAX(p.base_price)::DECIMAL AS max_price,
    MAX(p.base_price)::SMALLINT AS discount
FROM products AS p
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
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
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    image_id = coalesce(sqlc.narg('image_id'), image_id),
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