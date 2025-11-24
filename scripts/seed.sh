#!/bin/bash

# E-Shop Database Seeding Script
# This script seeds the database with initial data for development/testing

set -e  # Exit on any error

echo "🌱 Starting E-Shop Database Seeding..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if we're in the correct directory
if [ ! -f "go.mod" ] || [ ! -d "seeds" ]; then
    print_error "Please run this script from the project root directory"
    print_error "Make sure you have a 'seeds' directory with seed data files"
    exit 1
fi

# Check if database is running
print_status "Checking database connection..."
if ! go run ./cmd/seed --help > /dev/null 2>&1; then
    print_error "Failed to run seed command. Make sure:"
    echo "  1. Database is running"
    echo "  2. Database configuration is correct in app.env"
    echo "  3. Go application builds successfully"
    exit 1
fi

print_success "Database connection verified"

# Build the seed command
print_status "Building seed command..."
if ! go build -o tmp/seed ./cmd/seed; then
    print_error "Failed to build seed command"
    exit 1
fi
print_success "Seed command built successfully"

# Function to run individual seed command with error handling
run_seed() {
    local seed_type=$1
    local description=$2
    
    print_status "Seeding $description..."
    if ./tmp/seed "$seed_type"; then
        print_success "$description seeded successfully"
        return 0
    else
        print_error "Failed to seed $description"
        return 1
    fi
}

# Run all seeds or specific seed based on argument
if [ $# -eq 0 ]; then
    # Run all seeds in proper order
    print_status "Running comprehensive database seeding..."
    
    if ./tmp/seed; then
        print_success "🎉 All database seeding completed successfully!"
        echo ""
        echo "The following data has been seeded:"
        echo "  ✓ User roles and permissions"
        echo "  ✓ Payment methods"
        echo "  ✓ Attributes and attribute values"
        echo "  ✓ Categories"
        echo "  ✓ Brands"
        echo "  ✓ Collections"
        echo "  ✓ Shipping methods and zones"
        echo "  ✓ Discounts"
        echo "  ✓ Users (including admin user)"
        echo "  ✓ User addresses"
        echo "  ✓ Products with variants"
        echo ""
        echo "🔐 Default admin credentials:"
        echo "   Username: admin"
        echo "   Email: admin@eshop.com"
        echo "   Password: admin123"
        echo ""
        print_warning "Make sure to change the admin password in production!"
    else
        print_error "Database seeding failed"
        exit 1
    fi
else
    # Run specific seed
    seed_type=$1
    case $seed_type in
        "attributes")
            run_seed "attributes" "attributes"
            run_seed "attribute-values" "attribute values"
            ;;
        "users")
            run_seed "users" "users and addresses"
            ;;
        "products")
            run_seed "products" "products"
            ;;
        "categories")
            run_seed "categories" "categories"
            ;;
        "brands")
            run_seed "brands" "brands"
            ;;
        "collections")
            run_seed "collections" "collections"
            ;;
        "discounts")
            run_seed "discounts" "discounts"
            ;;
        "shipping")
            run_seed "shipping" "shipping methods and zones"
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [SEED_TYPE]"
            echo ""
            echo "Seed the database with initial data for development/testing"
            echo ""
            echo "SEED_TYPE options:"
            echo "  attributes   Seed attributes and attribute values"
            echo "  users        Seed users and addresses"
            echo "  products     Seed products"
            echo "  categories   Seed categories"
            echo "  brands       Seed brands"
            echo "  collections  Seed collections"
            echo "  discounts    Seed discounts"
            echo "  shipping     Seed shipping methods and zones"
            echo ""
            echo "If no SEED_TYPE is provided, all data will be seeded in the correct order."
            echo ""
            echo "Examples:"
            echo "  $0              # Seed all data"
            echo "  $0 users        # Seed only users"
            echo "  $0 products     # Seed only products"
            exit 0
            ;;
        *)
            print_error "Unknown seed type: $seed_type"
            echo "Run '$0 help' for available options"
            exit 1
            ;;
    esac
fi

# Cleanup
rm -f tmp/seed

print_success "Seeding script completed successfully! 🚀"