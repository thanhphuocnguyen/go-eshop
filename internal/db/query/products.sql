-- name: CreateProduct :one
INSERT INTO products (name, description, short_description, base_price, base_sku, slug, brand_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

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
SELECT p.*,
    JSON_BUILD_OBJECT('id', b.id, 'name', b.name) AS brand,
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT('id', c.id, 'name', c.name)) FILTER (WHERE c.id IS NOT NULL) AS categories,
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT('id', cl.id, 'name', cl.name)) FILTER (WHERE cl.id IS NOT NULL) AS collections,
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT('attributeId', a.id,'attributeName', a.name)) FILTER (WHERE a.id IS NOT NULL) AS attributes,
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', pv.id,
        'sku', pv.sku,
        'price', pv.price,
        'stock', pv.stock,
        'isActive', pv.is_active,
        'imageUrl', pv.image_url,
        'imageId', pv.image_id
    )) FILTER (WHERE pv.id IS NOT NULL) AS variants
FROM products p
LEFT JOIN brands AS b ON p.brand_id = b.id
LEFT JOIN category_products AS cp ON p.id = cp.product_id
LEFT JOIN collection_products AS colp ON p.id = colp.product_id
LEFT JOIN categories as c ON cp.category_id = c.id
LEFT JOIN collections as cl ON colp.collection_id = cl.id
LEFT JOIN product_attributes pa ON p.id = pa.product_id
LEFT JOIN attributes a ON pa.attribute_id = a.id
LEFT JOIN product_variants pv ON pv.product_id = p.id
WHERE (p.id = $1 OR p.slug = $2) AND p.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
GROUP BY p.id, b.id LIMIT 1;

-- name: GetProductVariantList :many
SELECT v.*, 
    JSONB_AGG(
        DISTINCT JSONB_BUILD_OBJECT('id', av.id, 'value', av.value, 'attribute_id', av.attribute_id)
    ) FILTER (WHERE av.id IS NOT NULL) AS attribute_values
FROM product_variants v
LEFT JOIN variant_attribute_values pva ON v.id = pva.variant_id
LEFT JOIN attribute_values av ON pva.attribute_value_id = av.id
WHERE v.product_id = $1 AND v.is_active = COALESCE(sqlc.narg('is_active'), v.is_active)
GROUP BY v.id
ORDER BY v.id, v.created_at DESC;

-- name: GetAdminProductList :many
SELECT p.* FROM products as p
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) 
    AND p.name ILIKE COALESCE(sqlc.narg('search'), p.name)
    AND p.slug ILIKE COALESCE(sqlc.narg('slug'), p.slug)
GROUP BY p.id ORDER BY @orderBy::text LIMIT $1 OFFSET $2;

-- name: GetProductList :many
SELECT p.*, MIN(pv.price) as min_price, COUNT(pv.id) as variant_count FROM products as p
LEFT JOIN collection_products cp ON p.id = cp.product_id
LEFT JOIN collections c ON cp.collection_id = c.id
LEFT JOIN category_products catp ON p.id = catp.product_id
LEFT JOIN categories cat ON catp.category_id = cat.id
LEFT JOIN brands b ON p.brand_id = b.id
LEFT JOIN product_variants pv ON pv.product_id = p.id
WHERE
    p.is_active = COALESCE(sqlc.narg('is_active'), p.is_active) 
    AND p.name ILIKE COALESCE(sqlc.narg('search'), '%')
    AND (sqlc.narg('brand_ids')::uuid[] is null or p.brand_id = ANY(sqlc.narg('brand_ids')::uuid[]))
    AND (sqlc.narg('collection_ids')::uuid[] is null or c.id = ANY(sqlc.narg('collection_ids')::uuid[]))
    AND (sqlc.narg('category_ids')::uuid[] is null or cat.id = ANY(sqlc.narg('category_ids')::uuid[]))
    AND pv.stock > 0
GROUP BY p.id
ORDER BY @orderBy::text LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE
    is_active = COALESCE(sqlc.narg('is_active'), is_active) AND name ILIKE COALESCE(sqlc.narg('name'), '%')
    AND brand_id = COALESCE(sqlc.narg('brand_id'), brand_id);

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    short_description = coalesce(sqlc.narg('short_description'), short_description),
    brand_id = coalesce(sqlc.narg('brand_id'), brand_id),
    slug = coalesce(sqlc.narg('slug'), slug),
    base_price = coalesce(sqlc.narg('base_price'), base_price),
    base_sku = coalesce(sqlc.narg('base_sku'), base_sku),
    is_active = coalesce(sqlc.narg('is_active'), is_active),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    image_id = coalesce(sqlc.narg('image_id'), image_id),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

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
INSERT INTO products (brand_id, name, description) VALUES ($1, $2, $3);

-- name: UpdateProductStock :one
UPDATE product_variants SET stock = stock - $1 WHERE id = $2 RETURNING *;




--- SHOP PRODUCT DISPLAY ----
-- name: GetShopProducts :many