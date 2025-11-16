-- name: CreateAttribute :one
INSERT INTO attributes (name) VALUES ($1) RETURNING *;

-- name: AddBulkAttributes :copyfrom
INSERT INTO attributes (name) VALUES ($1);

-- name: GetAttributeByID :one
SELECT * FROM attributes WHERE id = $1 ORDER BY id;

-- name: GetAttributes :many
SELECT 
    a.*, av.value as attr_value, av.id as attr_value_id
FROM attributes as a
LEFT JOIN attribute_values as av ON a.id = av.attribute_id
WHERE 
    CASE 
        WHEN array_length(sqlc.arg(ids)::int[], 1) > 0 THEN a.id = ANY(sqlc.arg(ids)::int[])
        ELSE true
    END
ORDER BY a.id;

-- name: GetAttributeByName :one
SELECT * FROM attributes WHERE name = $1 LIMIT 1;

-- name: GetAttributesByIDs :many
SELECT * FROM attributes WHERE id = ANY(sqlc.arg(ids)::int[]) ORDER BY attributes.id;

-- name: UpdateAttribute :one
UPDATE attributes SET name = $1 WHERE id = $2 RETURNING *;

-- name: CountAttributes :one
SELECT COUNT(*) FROM attributes;

-- name: DeleteAttribute :exec
DELETE FROM attributes WHERE id = $1;

-- Attribute values
-- name: CreateAttributeValue :one
INSERT INTO attribute_values (attribute_id, value) VALUES ($1, $2) RETURNING *;

-- name: CreateAttributeValues :copyfrom
INSERT INTO attribute_values (attribute_id, value) VALUES ($1, $2);

-- name: GetAttributeValueByID :one
SELECT * FROM attribute_values WHERE id = $1;

-- name: GetAttributeValues :many
SELECT * FROM attribute_values WHERE attribute_id = $1 ORDER BY attribute_values.id;

-- name: GetAttributeValuesByIDs :many
SELECT * FROM attribute_values WHERE id = ANY(sqlc.arg(ids)::bigint[]) ORDER BY attribute_values.id;

-- name: GetAttributeWithValuesByIDs :many
SELECT att.name as attribute_name, att.id as attribute_id, atv.value as attribute_value, atv.id as attribute_value_id 
FROM attributes as att
LEFT JOIN attribute_values as atv ON att.id = atv.attribute_id
WHERE attribute_id = ANY(sqlc.arg(ids)::int[]) 
ORDER BY atv.id;

-- name: UpdateAttributeValue :one
UPDATE attribute_values 
SET "value" = COALESCE(sqlc.narg('value'), "value")
WHERE id = $1 RETURNING *;

-- name: DeleteAttributeValue :exec
DELETE FROM attribute_values WHERE id = $1;

-- name: DeleteAttributeValueByValueID :exec
DELETE FROM attribute_values WHERE id = $2 AND attribute_id = $1;

-- Product Variant attributes
-- name: CreateProductVariantAttribute :one
INSERT INTO variant_attribute_values (variant_id, attribute_value_id) VALUES ($1, $2) RETURNING *;

-- name: CreateBulkProductVariantAttribute :copyfrom
INSERT INTO variant_attribute_values (variant_id, attribute_value_id) VALUES ($1, $2);

-- name: GetProductVariantAttributeByID :one
SELECT * FROM variant_attribute_values WHERE variant_id = $1 LIMIT 1;

-- name: GetProductVariantAttributes :many
SELECT * FROM variant_attribute_values WHERE variant_id = $1 ORDER BY variant_attribute_values.attribute_value_id;

-- name: DeleteProductVariantAttributes :exec
DELETE FROM variant_attribute_values WHERE variant_id = $1;