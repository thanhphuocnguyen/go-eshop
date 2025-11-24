-- name: CreateShippingMethod :one
INSERT INTO shipping_methods (name, description, is_active, requires_address, estimated_delivery_time, icon_url) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetShippingMethods :many
SELECT * FROM shipping_methods WHERE is_active = COALESCE(sqlc.narg('is_active'), TRUE) ORDER BY name;

-- name: GetShippingMethodByID :one
SELECT * FROM shipping_methods WHERE id = $1;

-- name: UpdateShippingMethod :one
UPDATE shipping_methods SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    requires_address = COALESCE(sqlc.narg('requires_address'), requires_address),
    estimated_delivery_time = COALESCE(sqlc.narg('estimated_delivery_time'), estimated_delivery_time),
    icon_url = COALESCE(sqlc.narg('icon_url'), icon_url),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteShippingMethod :exec
DELETE FROM shipping_methods WHERE id = $1;

-- name: CountShippingMethods :one
SELECT COUNT(*) FROM shipping_methods;

-- name: SeedShippingMethods :copyfrom
INSERT INTO shipping_methods (name, description, is_active, requires_address, estimated_delivery_time) 
VALUES ($1, $2, $3, $4, $5);

-- SHIPPING ZONES
-- name: CreateShippingZone :one
INSERT INTO shipping_zones (name, description, countries, states, zip_codes, is_active) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetShippingZones :many
SELECT * FROM shipping_zones WHERE is_active = COALESCE(sqlc.narg('is_active'), TRUE) ORDER BY name;

-- name: GetShippingZoneByID :one
SELECT * FROM shipping_zones WHERE id = $1;

-- name: UpdateShippingZone :one
UPDATE shipping_zones SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    countries = COALESCE(sqlc.narg('countries'), countries),
    states = COALESCE(sqlc.narg('states'), states),
    zip_codes = COALESCE(sqlc.narg('zip_codes'), zip_codes),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteShippingZone :exec
DELETE FROM shipping_zones WHERE id = $1;

-- name: CountShippingZones :one
SELECT COUNT(*) FROM shipping_zones;

-- name: SeedShippingZones :copyfrom
INSERT INTO shipping_zones (name, description, countries, states, zip_codes, is_active) 
VALUES ($1, $2, $3, $4, $5, $6);

-- SHIPPING RATES
-- name: CreateShippingRate :one
INSERT INTO shipping_rates (shipping_method_id, shipping_zone_id, name, base_rate, min_order_amount, max_order_amount, free_shipping_threshold, is_active) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetShippingRates :many
SELECT sr.*, sm.name as method_name, sz.name as zone_name
FROM shipping_rates sr
JOIN shipping_methods sm ON sr.shipping_method_id = sm.id
JOIN shipping_zones sz ON sr.shipping_zone_id = sz.id
WHERE sr.is_active = COALESCE(sqlc.narg('is_active'), TRUE)
ORDER BY sr.base_rate;

-- name: GetShippingRateByID :one
SELECT sr.*, sm.name as method_name, sz.name as zone_name
FROM shipping_rates sr
JOIN shipping_methods sm ON sr.shipping_method_id = sm.id
JOIN shipping_zones sz ON sr.shipping_zone_id = sz.id
WHERE sr.id = $1;

-- name: GetShippingRatesByZone :many
SELECT sr.*, sm.name as method_name
FROM shipping_rates sr
JOIN shipping_methods sm ON sr.shipping_method_id = sm.id
WHERE sr.shipping_zone_id = $1 AND sr.is_active = TRUE
ORDER BY sr.base_rate;

-- name: UpdateShippingRate :one
UPDATE shipping_rates SET
    name = COALESCE(sqlc.narg('name'), name),
    base_rate = COALESCE(sqlc.narg('base_rate'), base_rate),
    min_order_amount = COALESCE(sqlc.narg('min_order_amount'), min_order_amount),
    max_order_amount = COALESCE(sqlc.narg('max_order_amount'), max_order_amount),
    free_shipping_threshold = COALESCE(sqlc.narg('free_shipping_threshold'), free_shipping_threshold),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteShippingRate :exec
DELETE FROM shipping_rates WHERE id = $1;

-- name: CountShippingRates :one
SELECT COUNT(*) FROM shipping_rates;