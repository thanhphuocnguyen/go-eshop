# Database Schema Documentation

## Overview

This document describes the database schema for the e-commerce platform. The database uses PostgreSQL with UUID primary keys and follows normalized design principles.

## Database Design Principles

- **UUID Primary Keys**: All tables use UUID v4 as primary keys for better scalability and security
- **Soft Deletes**: Important entities support soft deletion to maintain data integrity
- **Audit Trail**: All tables include `created_at` and `updated_at` timestamps
- **Referential Integrity**: Foreign key constraints ensure data consistency
- **Enumerated Types**: PostgreSQL ENUM types for status fields to ensure data validity

## Core Tables

### Users and Authentication

#### `users`
Primary user information table.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique user identifier |
| `role_id` | UUID | NOT NULL, FK to user_roles | User role reference |
| `username` | VARCHAR | UNIQUE, NOT NULL | Unique username |
| `email` | VARCHAR | UNIQUE, NOT NULL | User email address |
| `phone_number` | VARCHAR(20) | NOT NULL, LENGTH(10-20) | Phone number |
| `first_name` | VARCHAR | NOT NULL | First name |
| `last_name` | VARCHAR | NOT NULL | Last name |
| `avatar_url` | VARCHAR | | Profile picture URL |
| `avatar_image_id` | VARCHAR | | Cloudinary image ID |
| `hashed_password` | VARCHAR | NOT NULL | Bcrypt hashed password |
| `verified_email` | BOOLEAN | DEFAULT FALSE | Email verification status |
| `verified_phone` | BOOLEAN | DEFAULT FALSE | Phone verification status |
| `locked` | BOOLEAN | DEFAULT FALSE | Account lock status |
| `password_changed_at` | TIMESTAMPTZ | DEFAULT '0001-01-01' | Last password change |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Account creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `user_roles`
User role definitions with hierarchical permissions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Role identifier |
| `code` | VARCHAR(50) | UNIQUE, NOT NULL | Role code (admin, user, moderator) |
| `name` | VARCHAR(100) | NOT NULL | Display name |
| `description` | TEXT | | Role description |
| `is_active` | BOOLEAN | DEFAULT TRUE | Active status |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Default Roles:**
- `admin`: Full system access
- `user`: Standard customer access
- `moderator`: Content management access

#### `permissions`
Linux-style permission system (read, write, execute) per module.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Permission identifier |
| `role_id` | UUID | FK to user_roles | Role reference |
| `module` | VARCHAR(100) | NOT NULL | Module name |
| `r` | BOOLEAN | DEFAULT FALSE | Read permission |
| `w` | BOOLEAN | DEFAULT FALSE | Write permission |
| `x` | BOOLEAN | DEFAULT FALSE | Execute permission |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Unique constraint:** `(role_id, module)`

#### `user_sessions`
Active user sessions for JWT token management.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Session identifier |
| `user_id` | UUID | FK to users | User reference |
| `refresh_token` | VARCHAR | NOT NULL | Refresh token |
| `user_agent` | VARCHAR(512) | NOT NULL | Client user agent |
| `client_ip` | INET | NOT NULL | Client IP address |
| `blocked` | BOOLEAN | DEFAULT FALSE | Session block status |
| `expired_at` | TIMESTAMPTZ | NOT NULL | Session expiration |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Session start time |

#### `user_addresses`
User shipping and billing addresses.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Address identifier |
| `user_id` | UUID | FK to users | User reference |
| `phone_number` | VARCHAR(20) | LENGTH(10-20) | Contact phone |
| `street` | VARCHAR | NOT NULL | Street address |
| `ward` | VARCHAR(100) | | Ward/district |
| `district` | VARCHAR(100) | NOT NULL | District |
| `city` | VARCHAR(100) | NOT NULL | City |
| `is_default` | BOOLEAN | DEFAULT FALSE | Default address flag |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `email_verifications`
Email verification tokens.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Verification identifier |
| `user_id` | UUID | FK to users | User reference |
| `email` | VARCHAR(255) | NOT NULL | Email to verify |
| `verify_code` | VARCHAR(255) | NOT NULL | Verification code |
| `is_used` | BOOLEAN | DEFAULT FALSE | Usage status |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `expired_at` | TIMESTAMPTZ | DEFAULT NOW() + 1 day | Expiration time |

### Product Catalog

#### `categories`
Product categories with hierarchical structure support.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Category identifier |
| `name` | VARCHAR(255) | UNIQUE, NOT NULL | Category name |
| `description` | TEXT | | Category description |
| `image_url` | TEXT | | Category image URL |
| `image_id` | VARCHAR | | Cloudinary image ID |
| `published` | BOOLEAN | DEFAULT TRUE | Visibility status |
| `slug` | VARCHAR | UNIQUE, NOT NULL | URL-friendly identifier |
| `display_order` | INTEGER | DEFAULT 0 | Sort order |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `brands`
Product brands/manufacturers.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Brand identifier |
| `name` | VARCHAR(255) | UNIQUE, NOT NULL | Brand name |
| `description` | TEXT | | Brand description |
| `image_url` | TEXT | | Brand logo URL |
| `image_id` | VARCHAR | | Cloudinary image ID |
| `slug` | VARCHAR | UNIQUE, NOT NULL | URL-friendly identifier |
| `display_order` | INTEGER | DEFAULT 0 | Sort order |
| `published` | BOOLEAN | DEFAULT TRUE | Visibility status |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `collections`
Product collections for marketing and grouping.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Collection identifier |
| `name` | VARCHAR(255) | UNIQUE, NOT NULL | Collection name |
| `description` | TEXT | | Collection description |
| `image_url` | TEXT | | Collection image URL |
| `image_id` | VARCHAR | | Cloudinary image ID |
| `slug` | VARCHAR | UNIQUE, NOT NULL | URL-friendly identifier |
| `display_order` | INTEGER | DEFAULT 0 | Sort order |
| `published` | BOOLEAN | DEFAULT TRUE | Visibility status |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `products`
Core product information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Product identifier |
| `name` | VARCHAR(255) | NOT NULL | Product name |
| `slug` | VARCHAR | UNIQUE, NOT NULL | URL-friendly identifier |
| `description` | TEXT | | Product description |
| `short_description` | TEXT | | Brief description |
| `brand_id` | UUID | FK to brands | Brand reference |
| `category_id` | UUID | FK to categories | Category reference |
| `collection_id` | UUID | FK to collections | Collection reference |
| `base_price` | DECIMAL(10,2) | NOT NULL | Base price |
| `sku` | VARCHAR | UNIQUE | Stock keeping unit |
| `status` | product_status | DEFAULT 'active' | Product status |
| `published` | BOOLEAN | DEFAULT FALSE | Visibility status |
| `featured` | BOOLEAN | DEFAULT FALSE | Featured product flag |
| `weight` | DECIMAL(8,2) | | Product weight (kg) |
| `dimensions` | JSONB | | Dimensions object |
| `metadata` | JSONB | | Additional product data |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Product Status Enum:** `'active', 'inactive', 'out_of_stock', 'discontinued'`

#### `product_variants`
Product variations (size, color, etc.).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Variant identifier |
| `product_id` | UUID | FK to products | Product reference |
| `name` | VARCHAR(255) | NOT NULL | Variant name |
| `sku` | VARCHAR | UNIQUE | Variant SKU |
| `price` | DECIMAL(10,2) | NOT NULL | Variant price |
| `compare_at_price` | DECIMAL(10,2) | | Original/compare price |
| `cost_per_item` | DECIMAL(10,2) | | Cost basis |
| `inventory_quantity` | INTEGER | DEFAULT 0 | Stock quantity |
| `weight` | DECIMAL(8,2) | | Variant weight |
| `requires_shipping` | BOOLEAN | DEFAULT TRUE | Shipping requirement |
| `taxable` | BOOLEAN | DEFAULT TRUE | Tax applicability |
| `barcode` | VARCHAR | | Barcode/UPC |
| `position` | INTEGER | DEFAULT 0 | Display order |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `images`
Product and variant images.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Image identifier |
| `product_id` | UUID | FK to products | Product reference |
| `variant_id` | UUID | FK to product_variants | Variant reference |
| `url` | VARCHAR | NOT NULL | Image URL |
| `cloudinary_id` | VARCHAR | NOT NULL | Cloudinary public ID |
| `alt_text` | VARCHAR(255) | | Alt text for accessibility |
| `position` | INTEGER | DEFAULT 0 | Display order |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Upload time |

#### `attributes`
Product attribute definitions (color, size, material).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Attribute identifier |
| `name` | VARCHAR(100) | UNIQUE, NOT NULL | Attribute name |

#### `attribute_values`
Possible values for each attribute.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Value identifier |
| `attribute_id` | INTEGER | FK to attributes | Attribute reference |
| `value` | VARCHAR(100) | NOT NULL | Attribute value |

**Unique constraint:** `(attribute_id, value)`

#### `variant_attributes`
Links variants to their attribute values.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `variant_id` | UUID | FK to product_variants | Variant reference |
| `attribute_value_id` | INTEGER | FK to attribute_values | Attribute value reference |

**Primary key:** `(variant_id, attribute_value_id)`

### Shopping and Orders

#### `shopping_carts`
User shopping carts.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Cart identifier |
| `user_id` | UUID | FK to users | User reference |
| `status` | cart_status | DEFAULT 'active' | Cart status |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Cart Status Enum:** `'active', 'checked_out'`

#### `cart_items`
Items in shopping carts.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Item identifier |
| `cart_id` | UUID | FK to shopping_carts | Cart reference |
| `product_id` | UUID | FK to products | Product reference |
| `variant_id` | UUID | FK to product_variants | Variant reference |
| `quantity` | INTEGER | NOT NULL, > 0 | Item quantity |
| `unit_price` | DECIMAL(10,2) | NOT NULL | Price per unit |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Unique constraint:** `(cart_id, variant_id)`

#### `orders`
Customer orders.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Order identifier |
| `user_id` | UUID | FK to users | Customer reference |
| `order_number` | VARCHAR | UNIQUE, NOT NULL | Human-readable order number |
| `status` | order_status | DEFAULT 'pending' | Order status |
| `currency` | VARCHAR(3) | DEFAULT 'USD' | Currency code |
| `subtotal` | DECIMAL(10,2) | NOT NULL | Subtotal amount |
| `tax_amount` | DECIMAL(10,2) | DEFAULT 0 | Tax amount |
| `shipping_cost` | DECIMAL(10,2) | DEFAULT 0 | Shipping cost |
| `discount_amount` | DECIMAL(10,2) | DEFAULT 0 | Discount amount |
| `total_amount` | DECIMAL(10,2) | NOT NULL | Final total |
| `shipping_address` | JSONB | NOT NULL | Shipping address object |
| `billing_address` | JSONB | | Billing address object |
| `notes` | TEXT | | Order notes |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Order time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Order Status Enum:** `'pending', 'confirmed', 'delivering', 'delivered', 'cancelled', 'refunded', 'completed'`

#### `order_items`
Items within orders.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Item identifier |
| `order_id` | UUID | FK to orders | Order reference |
| `product_id` | UUID | FK to products | Product reference |
| `variant_id` | UUID | FK to product_variants | Variant reference |
| `quantity` | INTEGER | NOT NULL | Ordered quantity |
| `unit_price` | DECIMAL(10,2) | NOT NULL | Price per unit |
| `total_price` | DECIMAL(10,2) | NOT NULL | Total line price |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |

### Payments

#### `payment_methods`
Available payment method configurations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Method identifier |
| `code` | VARCHAR(50) | UNIQUE, NOT NULL | Method code |
| `name` | VARCHAR(100) | NOT NULL | Display name |
| `description` | TEXT | | Method description |
| `is_active` | BOOLEAN | DEFAULT TRUE | Active status |
| `gateway_supported` | VARCHAR(100) | | Gateway provider |
| `requires_account` | BOOLEAN | DEFAULT FALSE | Account requirement |
| `min_amount` | DECIMAL(10,2) | | Minimum transaction |
| `max_amount` | DECIMAL(10,2) | | Maximum transaction |
| `processing_fee_percentage` | DECIMAL(5,4) | | Percentage fee |
| `processing_fee_fixed` | DECIMAL(10,2) | | Fixed fee amount |
| `currency_supported` | TEXT[] | | Supported currencies |
| `countries_supported` | TEXT[] | | Supported countries |
| `metadata` | JSONB | | Additional configuration |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

#### `payments`
Payment transactions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Payment identifier |
| `order_id` | UUID | FK to orders | Order reference |
| `user_id` | UUID | FK to users | User reference |
| `payment_method_id` | UUID | FK to payment_methods | Method reference |
| `amount` | DECIMAL(10,2) | NOT NULL | Payment amount |
| `currency` | VARCHAR(3) | DEFAULT 'USD' | Currency code |
| `status` | payment_status | DEFAULT 'pending' | Payment status |
| `gateway_payment_id` | VARCHAR(255) | | External payment ID |
| `gateway_response` | JSONB | | Gateway response data |
| `failure_reason` | TEXT | | Failure description |
| `processed_at` | TIMESTAMPTZ | | Processing timestamp |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Payment Status Enum:** `'pending', 'success', 'failed', 'cancelled', 'refunded', 'processing'`

#### `user_payment_infos`
Saved user payment information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Info identifier |
| `user_id` | UUID | FK to users | User reference |
| `card_number` | VARCHAR(19) | NOT NULL | Masked card number |
| `card_last4` | VARCHAR(4) | NOT NULL | Last 4 digits |
| `payment_method_token` | VARCHAR(255) | NOT NULL | Gateway token |
| `expiration_date` | DATE | NOT NULL | Card expiration |
| `billing_address` | TEXT | NOT NULL | Billing address |
| `is_default` | BOOLEAN | DEFAULT FALSE | Default payment flag |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Unique constraint:** `(user_id, payment_method_token)`

### Reviews and Ratings

#### `ratings`
Product ratings and reviews.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Rating identifier |
| `product_id` | UUID | FK to products | Product reference |
| `user_id` | UUID | FK to users | User reference |
| `rating` | INTEGER | CHECK (1 <= rating <= 5) | Star rating (1-5) |
| `title` | VARCHAR(255) | | Review title |
| `content` | TEXT | | Review content |
| `verified_purchase` | BOOLEAN | DEFAULT FALSE | Purchase verification |
| `helpful_count` | INTEGER | DEFAULT 0 | Helpful votes |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Review time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Unique constraint:** `(product_id, user_id)`

### Discounts and Promotions

#### `discounts`
Discount and promotion rules.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Discount identifier |
| `name` | VARCHAR(255) | NOT NULL | Discount name |
| `code` | VARCHAR(50) | UNIQUE | Discount code |
| `description` | TEXT | | Description |
| `type` | discount_type | NOT NULL | Discount type |
| `value` | DECIMAL(10,2) | NOT NULL | Discount value |
| `minimum_amount` | DECIMAL(10,2) | | Minimum order amount |
| `maximum_amount` | DECIMAL(10,2) | | Maximum discount amount |
| `usage_limit` | INTEGER | | Total usage limit |
| `used_count` | INTEGER | DEFAULT 0 | Current usage count |
| `user_limit` | INTEGER | DEFAULT 1 | Per-user usage limit |
| `is_active` | BOOLEAN | DEFAULT TRUE | Active status |
| `starts_at` | TIMESTAMPTZ | | Start date |
| `ends_at` | TIMESTAMPTZ | | End date |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation time |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update time |

**Discount Type Enum:** `'percentage', 'fixed_amount', 'buy_x_get_y', 'free_shipping', 'tiered'`

#### `discount_usages`
Track discount usage by users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Usage identifier |
| `discount_id` | UUID | FK to discounts | Discount reference |
| `user_id` | UUID | FK to users | User reference |
| `order_id` | UUID | FK to orders | Order reference |
| `used_at` | TIMESTAMPTZ | DEFAULT NOW() | Usage time |

**Unique constraint:** `(discount_id, user_id, order_id)`

## Indexes

### Performance Indexes

```sql
-- User lookups
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role_id ON users(role_id);

-- Product searches
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_collection_id ON products(collection_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_published ON products(published);

-- Order queries
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- Cart operations
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_shopping_carts_user_id ON shopping_carts(user_id);

-- Session management
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
```

### Full-text Search Indexes

```sql
-- Product search
CREATE INDEX idx_products_search ON products 
USING GIN(to_tsvector('english', name || ' ' || COALESCE(description, '')));

-- Category search
CREATE INDEX idx_categories_search ON categories 
USING GIN(to_tsvector('english', name || ' ' || COALESCE(description, '')));
```

## Database Functions and Triggers

### Automatic Timestamp Updates

```sql
-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### Order Number Generation

```sql
-- Function to generate sequential order numbers
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
    NEW.order_number = 'ORD-' || LPAD(nextval('order_number_seq')::TEXT, 8, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE SEQUENCE order_number_seq;

CREATE TRIGGER set_order_number BEFORE INSERT ON orders
    FOR EACH ROW EXECUTE FUNCTION generate_order_number();
```

## Backup and Maintenance

### Regular Maintenance Tasks

1. **Vacuum and Analyze**
   ```sql
   VACUUM ANALYZE;
   ```

2. **Reindex for Performance**
   ```sql
   REINDEX DATABASE eshop;
   ```

3. **Clean Old Sessions**
   ```sql
   DELETE FROM user_sessions WHERE expired_at < NOW() - INTERVAL '7 days';
   ```

4. **Archive Old Orders**
   ```sql
   -- Move orders older than 2 years to archive table
   INSERT INTO orders_archive SELECT * FROM orders 
   WHERE created_at < NOW() - INTERVAL '2 years';
   ```

### Backup Strategy

1. **Daily Full Backup**
   ```bash
   pg_dump -h localhost -U postgres -d eshop -f backup_$(date +%Y%m%d).sql
   ```

2. **Point-in-Time Recovery**
   ```bash
   # Enable WAL archiving in postgresql.conf
   archive_mode = on
   archive_command = 'cp %p /path/to/archive/%f'
   ```

## Security Considerations

1. **Password Storage**: All passwords are hashed using bcrypt
2. **Sensitive Data**: Payment tokens are encrypted and never stored in plain text
3. **Audit Trail**: All data modifications are logged with timestamps
4. **Access Control**: Row-level security can be implemented for multi-tenant scenarios
5. **Data Encryption**: Consider encrypting PII fields at rest

## Performance Optimization

1. **Connection Pooling**: Use pgbouncer or similar for connection management
2. **Query Optimization**: Regular analysis of slow queries
3. **Partitioning**: Consider partitioning large tables (orders, payments) by date
4. **Caching**: Implement Redis caching for frequently accessed data
5. **Read Replicas**: Use read replicas for reporting and analytics queries