# Go eShop - E-Commerce Backend API

A modern e-commerce backend built with Go, featuring JWT authentication, payment processing, and clean architecture.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

## Features

- 🔐 JWT/PASETO authentication
- 🛍️ Product catalog with search and filters
- 🛒 Shopping cart management
- 💳 Payment processing (Stripe)
- 📦 Order management and tracking
- ⭐ Product reviews and ratings
- 👥 User management and profiles
- 📊 Admin dashboard
- 🔍 Elasticsearch integration
## Tech Stack

- **Backend:** Go 1.24+ with Chi router
- **Database:** PostgreSQL 14+ with SQLC ORM
- **Cache:** Redis 6+
- **Search:** Elasticsearch  
- **Jobs:** Asynq for background processing
- **Payments:** Stripe integration
- **Media:** Cloudinary for image management
- **Auth:** JWT/PASETO tokens

## Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- Make

### Setup

1. **Clone the repository**
```bash
git clone https://github.com/thanhphuocnguyen/go-eshop.git
cd go-eshop/server
```

2. **Start with Docker**
```bash
make dev-setup
```

This will start PostgreSQL, Redis, run migrations, seed data, and start the server.

3. **Manual setup** (alternative)
```bash
# Copy environment file
cp app.env.example app.env

# Start infrastructure
docker-compose up -d postgres redis

# Install dependencies
go mod download

# Run migrations
make migrate-up

# Seed sample data (optional)
make seed

# Start server
make serve-server
```

4. **Verify installation**
- API: http://localhost:4000/health
- Swagger docs: http://localhost:4000/swagger/index.html

## API Usage

### Authentication
```bash
# Register
curl -X POST http://localhost:4000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user","email":"user@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:4000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password123"}'
```

### Products
```bash
# Get products
curl http://localhost:4000/api/v1/products?page=1&limit=10

# Search products
curl http://localhost:4000/api/v1/products/search?q=laptop
```

### Cart & Orders (requires auth token)
```bash
# Add to cart
curl -X POST http://localhost:4000/api/v1/cart/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"product_id":"uuid","quantity":1}'

# Create order
curl -X POST http://localhost:4000/api/v1/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"shipping_address_id":"uuid","payment_method":"stripe"}'
```

## Development

```bash
# Run with hot reload
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## Project Structure

```
cmd/               # Application entrypoints
├── web/          # Main web server
├── migrate/      # Database migrations
└── seed/         # Data seeding

internal/          # Private application code
├── api/          # HTTP handlers and routes
├── models/       # Database models
├── dto/          # Data transfer objects
├── utils/        # Utility functions
└── worker/       # Background jobs

pkg/              # Public libraries
├── auth/         # Authentication utilities
├── cache/        # Cache implementations
└── payment/      # Payment processing

migrations/       # Database migration files
seeds/           # Seed data files
```

## Environment Variables

```env
# Server
PORT=4000
ENV=development

# Database
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop

# Redis
REDIS_URL=localhost:6380

# Auth
SYMMETRIC_KEY=your-32-character-secret-key
ACCESS_TOKEN_DURATION=24h

# Stripe (optional)
STRIPE_SECRET_KEY=sk_test_...

# Cloudinary (optional)
CLOUDINARY_URL=cloudinary://...
```

## API Documentation

- **Swagger UI:** http://localhost:4000/swagger/index.html
- **API Reference:** [docs/API.md](docs/API.md)
- **Database Schema:** [docs/DATABASE.md](docs/DATABASE.md)

## Deployment

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for detailed deployment instructions.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.