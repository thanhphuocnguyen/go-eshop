# Database Seeding Guide

This guide explains how to seed the E-Shop database with initial data for development and testing purposes.

## Overview

The seeding system provides realistic test data including:
- **User Management**: Roles, permissions, users, and addresses
- **Product Catalog**: Categories, brands, collections, products with variants
- **Attributes**: Product attributes and their values (Color, Size, Material, etc.)
- **Commerce**: Discounts, shipping methods, zones, and payment methods
- **Sample Data**: 10+ realistic products with variants, multiple users with addresses

## Quick Start

### 1. Prerequisites
```bash
# Make sure your database is running and configured
# Check your app.env file for database connection settings
```

### 2. Run All Seeds
```bash
# From the project root directory
./scripts/seed.sh
```

### 3. Default Admin User
After seeding, you can log in with:
- **Username**: `admin`
- **Email**: `admin@eshop.com`
- **Password**: `admin123`

⚠️ **Important**: Change the admin password in production!

## Selective Seeding

You can seed specific data types individually:

```bash
# Seed specific data types
./scripts/seed.sh attributes   # Attributes and values
./scripts/seed.sh users        # Users and addresses  
./scripts/seed.sh products     # Products and variants
./scripts/seed.sh categories   # Product categories
./scripts/seed.sh brands       # Product brands
./scripts/seed.sh collections  # Product collections
./scripts/seed.sh discounts    # Discount codes
./scripts/seed.sh shipping     # Shipping methods and zones
```

## Manual Seeding (Go Command)

If you prefer to use the Go command directly:

```bash
# Build the seed command
go build ./cmd/seed

# Run all seeds
./seed

# Run specific seeds
./seed users
./seed products
./seed categories
# ... etc
```

## Seed Data Details

### Users & Authentication
- **Admin User**: Full system access
- **Regular Users**: 9 test users with realistic profiles
- **Addresses**: 15 realistic US addresses distributed among users

### Product Catalog
- **Categories**: 8 main categories (Electronics, Fashion, etc.)
- **Brands**: 10 popular brands (Apple, Nike, Samsung, etc.)
- **Collections**: 8 curated collections (Featured, Best Sellers, etc.)
- **Products**: 10+ realistic products with multiple variants
- **Attributes**: 7 attribute types with comprehensive value sets

### Commerce Features
- **Discounts**: 5 different discount types with various conditions
- **Shipping**: 5 shipping methods with different zones and rates
- **Payments**: Pre-configured payment methods (Credit Card, PayPal, etc.)

### Sample Products Include
- iPhone 15 Pro (multiple storage/color variants)
- Samsung Galaxy S24 Ultra
- Nike Air Max 270 (multiple sizes/colors)
- MacBook Pro 14\"
- Sony WH-1000XM5 Headphones
- Levi's 501 Jeans (multiple sizes)
- IKEA Furniture
- Athletic wear and more...

## Seed File Structure

```
seeds/
├── categories.json           # Product categories
├── brands.json              # Product brands  
├── collections.json         # Product collections
├── products.json            # Products with variants
├── users.json               # User accounts
├── addresses.json           # User addresses
├── attribute_values.json    # Attribute values for products
├── discounts.json           # Discount codes and rules
├── shipping_methods.json    # Shipping options
└── shipping_zones.json      # Shipping zones and rates
```

## Customization

### Adding New Seed Data

1. **Create/Edit JSON Files**: Modify files in the `seeds/` directory
2. **Follow Existing Format**: Use the same structure as existing files
3. **Re-run Seeds**: Use the script to apply changes

### Example: Adding a New Category
```json
// seeds/categories.json
{
  "name": "New Category",
  "description": "Description of the new category",
  "sort_order": 9,
  "published": true
}
```

### Example: Adding a New Product
```json
// seeds/products.json
{
  "name": "New Product",
  "description": "Product description",
  "price": 99.99,
  "stock": 100,
  "sku": "NEW-PRODUCT-SKU",
  "variants": [
    {
      "price": 99.99,
      "stock": 50,
      "sku": "NEW-PRODUCT-RED-M",
      "attributes": {
        "Color": "Red",
        "Size": "Medium"
      }
    }
  ]
}
```

## Development Workflow

### For New Features
1. Run seeds to get test data: `./scripts/seed.sh`
2. Develop your feature with realistic data
3. Test with seeded products, users, and orders

### For Testing
1. Reset database (migration down/up)
2. Run seeds: `./scripts/seed.sh`
3. Run your tests with consistent data

### For Demos
1. Use the seeded admin account for full access
2. Showcase with realistic product catalog
3. Demonstrate e-commerce flows with test users

## Troubleshooting

### Database Connection Issues
```bash
# Check if database is running
docker ps  # or pg_isready if using local PostgreSQL

# Check environment variables
cat app.env | grep DB
```

### Seed Already Exists Errors
The seeding system is idempotent - it checks if data already exists before inserting. This is normal behavior.

### Build Errors
```bash
# Make sure you're in the project root
pwd  # should show /path/to/eshop/server

# Check for missing dependencies
go mod tidy
go mod download
```

### Permission Errors
```bash
# Make sure the script is executable
chmod +x scripts/seed.sh
```

## Production Considerations

⚠️ **Never run seeds in production!**

- Seeds are for development and testing only
- Contains default passwords and test data
- May overwrite existing production data

For production initialization:
1. Use proper database migrations
2. Create admin accounts manually with secure passwords
3. Import real product catalogs through admin interface
4. Configure payment methods through secure channels

## Contributing

When adding new seed data:
1. Keep it realistic and useful for testing
2. Follow existing JSON structure and naming conventions
3. Update this README if adding new seed types
4. Test your changes with `./scripts/seed.sh`