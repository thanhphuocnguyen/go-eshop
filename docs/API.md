# API Documentation

## Overview

This document provides detailed information about the e-commerce API endpoints, authentication, and data models.

## Base URL

```
http://localhost:4000/api/v1
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Token Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Login user |
| POST | `/auth/refresh` | Refresh access token |
| POST | `/auth/logout` | Logout user |
| POST | `/auth/verify-email` | Verify user email |

## User Management

### User Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users/profile` | Get current user profile | Yes |
| PUT | `/users/profile` | Update user profile | Yes |
| POST | `/users/avatar` | Upload user avatar | Yes |

### Address Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/address` | Get user addresses | Yes |
| POST | `/address` | Create new address | Yes |
| PUT | `/address/:id` | Update address | Yes |
| DELETE | `/address/:id` | Delete address | Yes |
| PUT | `/address/:id/default` | Set default address | Yes |

## Product Catalog

### Public Product Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/products` | Get products (paginated) | No |
| GET | `/products/:id` | Get product details | No |
| GET | `/products/:id/variants` | Get product variants | No |
| GET | `/products/:id/ratings` | Get product ratings | No |

### Categories

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/categories` | Get all categories | No |
| GET | `/categories/:slug` | Get category by slug | No |

### Brands

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/brands` | Get all brands | No |
| GET | `/brands/:slug` | Get brand by slug | No |

### Collections

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/collections` | Get all collections | No |
| GET | `/collections/:slug` | Get collection by slug | No |

## Shopping Cart

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/cart` | Get user's cart | Yes |
| POST | `/cart/items` | Add item to cart | Yes |
| PUT | `/cart/items/:id` | Update cart item | Yes |
| DELETE | `/cart/items/:id` | Remove item from cart | Yes |
| DELETE | `/cart/clear` | Clear entire cart | Yes |

## Orders

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/orders` | Get user orders | Yes |
| GET | `/orders/:id` | Get order details | Yes |
| POST | `/orders/checkout` | Checkout cart | Yes |
| PUT | `/orders/:id/cancel` | Cancel order | Yes |

## Payments

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/payments/intent` | Create payment intent | Yes |
| GET | `/payments/methods` | Get payment methods | No |
| POST | `/webhook/v1/stripe` | Stripe webhook endpoint | No |

## Admin Endpoints

### Admin User Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/users` | Get all users | Admin |
| GET | `/admin/users/:id` | Get user details | Admin |

### Admin Product Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/products` | Get all products | Admin |
| POST | `/admin/products` | Create product | Admin |
| PUT | `/admin/products/:id` | Update product | Admin |
| DELETE | `/admin/products/:id` | Delete product | Admin |
| POST | `/admin/products/:id/images` | Upload product images | Admin |
| POST | `/admin/products/:id/variants` | Create product variant | Admin |
| PUT | `/admin/products/:id/variants/:variantId` | Update variant | Admin |
| DELETE | `/admin/products/:id/variants/:variantId` | Delete variant | Admin |

### Admin Category Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/categories` | Create category | Admin |
| PUT | `/admin/categories/:id` | Update category | Admin |
| DELETE | `/admin/categories/:id` | Delete category | Admin |
| POST | `/admin/categories/:id/image` | Upload category image | Admin |

### Admin Brand Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/brands` | Create brand | Admin |
| PUT | `/admin/brands/:id` | Update brand | Admin |
| DELETE | `/admin/brands/:id` | Delete brand | Admin |
| POST | `/admin/brands/:id/image` | Upload brand image | Admin |

### Admin Collection Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/collections` | Create collection | Admin |
| PUT | `/admin/collections/:id` | Update collection | Admin |
| DELETE | `/admin/collections/:id` | Delete collection | Admin |
| POST | `/admin/collections/:id/image` | Upload collection image | Admin |

### Admin Order Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/orders` | Get all orders | Admin |
| PUT | `/admin/orders/:id/status` | Update order status | Admin |

### Admin Discount Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/discounts` | Get all discounts | Admin |
| POST | `/admin/discounts` | Create discount | Admin |
| PUT | `/admin/discounts/:id` | Update discount | Admin |
| DELETE | `/admin/discounts/:id` | Delete discount | Admin |

## Data Models

### User

```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "phone_number": "string",
  "first_name": "string",
  "last_name": "string",
  "avatar_url": "string",
  "role": "admin|user|moderator",
  "verified_email": "boolean",
  "verified_phone": "boolean",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Product

```json
{
  "id": "uuid",
  "name": "string",
  "slug": "string",
  "description": "string",
  "brand_id": "uuid",
  "category_id": "uuid",
  "collection_id": "uuid",
  "base_price": "decimal",
  "sku": "string",
  "status": "active|inactive|out_of_stock",
  "published": "boolean",
  "featured": "boolean",
  "weight": "decimal",
  "dimensions": "object",
  "images": ["string"],
  "variants": ["object"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Order

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "status": "pending|confirmed|delivering|delivered|cancelled|refunded|completed",
  "total_amount": "decimal",
  "subtotal": "decimal",
  "tax_amount": "decimal",
  "shipping_cost": "decimal",
  "discount_amount": "decimal",
  "shipping_address": "object",
  "billing_address": "object",
  "items": ["object"],
  "payment_status": "pending|success|failed|cancelled|refunded|processing",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Cart Item

```json
{
  "id": "uuid",
  "cart_id": "uuid",
  "product_id": "uuid",
  "variant_id": "uuid",
  "quantity": "integer",
  "unit_price": "decimal",
  "total_price": "decimal",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details (optional)"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `INTERNAL_ERROR` | 500 | Internal server error |

## Pagination

List endpoints support pagination with the following query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field
- `order`: Sort order (`asc` or `desc`)

Example:
```
GET /api/v1/products?page=2&limit=20&sort=created_at&order=desc
```

Response includes pagination metadata:

```json
{
  "data": [...],
  "pagination": {
    "page": 2,
    "limit": 20,
    "total": 150,
    "pages": 8
  }
}
```

## Rate Limiting

The API implements rate limiting:

- **Public endpoints**: 100 requests per minute per IP
- **Authenticated endpoints**: 1000 requests per minute per user
- **Admin endpoints**: 500 requests per minute per admin user

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## File Upload

File upload endpoints accept multipart/form-data with the following constraints:

- **Maximum file size**: 10MB
- **Supported formats**: JPEG, PNG, WebP, GIF
- **Image processing**: Automatic resizing and optimization via Cloudinary

Example:
```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -F "image=@product-image.jpg" \
  http://localhost:4000/api/v1/admin/products/123/images
```