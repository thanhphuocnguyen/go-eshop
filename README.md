# 🛍️ eShop - Modern E-Commerce Platform

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-20+-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

A comprehensive, production-ready e-commerce platform built with Go, featuring a robust REST API, modern architecture, and enterprise-level features. Designed for scalability, performance, and maintainability.

## 📑 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

## 🌟 Features

### 🛒 Customer Experience
- **🔐 Secure Authentication** - JWT/PASETO with refresh tokens and email verification
- **🔍 Advanced Product Search** - Full-text search with filters and sorting
- **🛍️ Smart Shopping Cart** - Persistent cart with real-time updates
- **💳 Multiple Payment Options** - Stripe integration with multiple payment methods
- **📦 Order Management** - Real-time order tracking and history
- **👤 User Profiles** - Comprehensive profile management and address book
- **⭐ Product Reviews** - Rating and review system with verification
- **🏷️ Discount System** - Coupons, promotions, and dynamic pricing

### 🎛️ Admin & Management
- **📊 Analytics Dashboard** - Sales metrics, user insights, and performance analytics
- **📝 Product Management** - Complete CRUD operations with variant support
- **📋 Order Processing** - Order status updates and fulfillment tracking
- **👥 User Management** - Customer account management and role-based access
- **💰 Payment Processing** - Transaction monitoring and refund handling
- **🎨 Content Management** - Categories, brands, collections, and media management
- **📈 Inventory Control** - Stock management with low-inventory alerts
- **🔧 System Settings** - Configurable business rules and system parameters

### 🔧 Technical Features
- **🚀 High Performance** - Optimized queries, caching, and connection pooling
- **🔒 Enterprise Security** - Input validation, rate limiting, and security headers
- **📱 API-First Design** - RESTful API with comprehensive Swagger documentation
- **🔄 Background Jobs** - Asynchronous task processing with Redis/Asynq
- **📧 Email Service** - Transactional emails with HTML templates
- **☁️ Cloud Integration** - Cloudinary for image management and CDN
- **🔍 Observability** - Structured logging, metrics, and health checks
- **🐳 Containerization** - Docker support with multi-stage builds

## 🏗️ Architecture

### 🎯 Clean Architecture Principles
```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│   Presentation      │    │   Business Logic    │    │   Data Access       │
│   (HTTP Handlers)   │───▶│   (Services)        │───▶│   (Repository)      │
│   • REST API        │    │   • Domain Logic    │    │   • PostgreSQL      │
│   • Middleware      │    │   • Validation      │    │   • SQLC Generated  │
│   • Request/Response│    │   • Business Rules  │    │   • Transactions    │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
           │                          │                          │
           ▼                          ▼                          ▼
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│   External Services │    │   Infrastructure    │    │   Database          │
│   • Stripe API      │    │   • Redis Cache     │    │   • PostgreSQL 14+  │
│   • Cloudinary      │    │   • Background Jobs │    │   • Migrations      │
│   • Email Service   │    │   • Logging         │    │   • Indexes         │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

### 🔧 Technology Stack

#### 🖥️ Backend (Go/Golang)
- **🌐 Web Framework**: Gin (HTTP router and middleware)
- **🗄️ Database**: PostgreSQL 14+ (primary database with ACID compliance)
- **⚡ Cache**: Redis 6+ (session storage, caching, and rate limiting)
- **🔐 Authentication**: PASETO/JWT tokens with refresh token rotation
- **💳 Payments**: Stripe integration with webhook support
- **☁️ File Storage**: Cloudinary (image hosting, optimization, and CDN)
- **📊 API Docs**: Swagger/OpenAPI 3.0 with automatic generation
- **🔄 Background Jobs**: Asynq (Redis-based job queue)
- **📝 Logging**: Zerolog (structured JSON logging)
- **🧪 Testing**: Testify, Gomock (unit and integration testing)
- **🔗 Database ORM**: SQLC (type-safe SQL code generation)
- **🚀 Migrations**: golang-migrate (version-controlled schema changes)

#### 📱 Frontend Options (Recommended)
- **⚛️ React 19** with TypeScript for type safety
- **🎨 UI Framework**: Tailwind CSS for utility-first styling
- **📡 Data Fetching**: SWR or TanStack Query for server state management
- **📋 Forms**: React Hook Form with validation
- **🎭 Components**: Headless UI, Radix UI, or Shadcn/ui
- **⚡ Build Tool**: Vite or Next.js for optimized builds

#### 🛠️ Development & DevOps
- **🐳 Containerization**: Docker with multi-stage builds
- **🔄 CI/CD**: GitHub Actions (recommended)
- **📊 Monitoring**: Prometheus + Grafana (optional)
- **📋 Linting**: golangci-lint with custom rules
- **🔧 Task Runner**: Make with comprehensive task definitions

## 🚀 Quick Start

### 📋 Prerequisites
- **Go 1.24+** - [Download here](https://golang.org/dl/)
- **PostgreSQL 14+** - [Installation guide](https://www.postgresql.org/download/)
- **Redis 6+** - [Installation guide](https://redis.io/docs/getting-started/)
- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)
- **Make** - Usually pre-installed on Unix systems

### ⚡ Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/thanhphuocnguyen/go-eshop.git
   cd go-eshop/server
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Environment Setup**
   ```bash
   # Copy environment template
   cp app.env.example app.env
   
   # Edit configuration (update database URLs, API keys, etc.)
   nano app.env
   ```

4. **Start Infrastructure Services**
   ```bash
   # Start PostgreSQL and Redis using Docker
   docker-compose up -d postgres redis
   
   # Or install locally and start services
   # brew services start postgresql redis  # macOS
   # sudo systemctl start postgresql redis  # Linux
   ```

5. **Database Setup**
   ```bash
   # Run database migrations
   make migrate-up
   
   # Seed with sample data (optional)
   make seed
   ```

6. **Start the Development Server**
   ```bash
   # Method 1: Using Make (recommended)
   make serve-server
   
   # Method 2: Direct Go command
   go run ./cmd/web api
   
   # Method 3: Using Air for live reload
   air
   ```

The API server will be available at: **http://localhost:4000**

### 🔧 Essential Configuration

Update your `app.env` file with the following key settings:

```env
# Database Configuration
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
MAX_POOL_SIZE=10

# Redis Configuration
REDIS_URL=localhost:6380

# Server Configuration
DOMAIN=localhost
PORT=4000
ENV=development

# Authentication (generate secure keys)
SYMMETRIC_KEY=your-32-character-secret-key-here
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h

# External Services (optional for development)
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
STRIPE_SECRET_KEY=sk_test_your_stripe_test_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_test_key
```

### 🎯 Verification

Test your installation with these quick checks:

```bash
# Health check
curl http://localhost:4000/health

# API version
curl http://localhost:4000/api/v1/health

# Swagger documentation
open http://localhost:4000/swagger/index.html
```

## 📚 API Documentation

### 🌐 Interactive Documentation

Once the server is running, comprehensive API documentation is available at:

**Swagger UI**: [http://localhost:4000/swagger/index.html](http://localhost:4000/swagger/index.html)

### 🔑 Authentication

The API uses Bearer token authentication. Include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:4000/api/v1/protected-endpoint
```

### 🚀 Quick API Examples

#### User Registration
```bash
curl -X POST http://localhost:4000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "phone_number": "+1234567890",
    "first_name": "John",
    "last_name": "Doe",
    "password": "SecurePassword123!"
  }'
```

#### User Login
```bash
curl -X POST http://localhost:4000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePassword123!"
  }'
```

#### Get Products
```bash
curl "http://localhost:4000/api/v1/products?page=1&limit=10&sort=name&order=asc"
```

#### Add Item to Cart (requires authentication)
```bash
curl -X POST http://localhost:4000/api/v1/cart/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "variant_id": "550e8400-e29b-41d4-a716-446655440001",
    "quantity": 2
  }'
```

### 📄 Additional Documentation

For detailed API documentation, database schemas, and development guides, see the `docs/` directory:

- **[API Reference](docs/API.md)** - Complete API endpoint documentation
- **[Database Schema](docs/DATABASE.md)** - Database design and relationships
- **[Development Guide](docs/DEVELOPMENT.md)** - Development setup and guidelines
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment instructions
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to the project

## 🧪 Testing

### 🔍 Test Coverage

We maintain high test coverage across all layers:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run tests with race detection
go test -race ./...
```

### 🧪 Test Types

#### Unit Tests
```bash
# Run unit tests only
go test -short ./...

# Test specific package
go test ./internal/api
```

#### Integration Tests
```bash
# Run integration tests (requires database)
go test -tags=integration ./...
```

#### API Tests
```bash
# Run API endpoint tests
go test ./tests/api/...
```

### 📊 Test Structure

```
tests/
├── unit/           # Unit tests for individual functions
├── integration/    # Integration tests with database
├── api/           # API endpoint tests
├── fixtures/      # Test data and fixtures
└── mocks/         # Generated mocks for testing
```

### 🎯 Testing Best Practices

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test complete workflows with real database
- **API Tests**: Test HTTP endpoints with full request/response cycle
- **Mocking**: Use generated mocks for external dependencies
- **Test Data**: Use factories and fixtures for consistent test data

## 📁 Project Structure

```
server/
├── 📁 cmd/                     # Application entry points
│   ├── 🔄 migrate/            # Database migration CLI tool
│   ├── 🌱 seed/               # Database seeding CLI tool
│   └── 🌐 web/                # Main API server application
│
├── ⚙️ config/                  # Configuration management
│   └── config.go              # Configuration struct and loader
│
├── 📖 docs/                    # Documentation and generated docs
│   ├── 🔄 docs.go             # Generated Swagger documentation
│   ├── 📄 API.md              # API documentation
│   ├── 🗄️ DATABASE.md         # Database schema documentation
│   ├── 🔧 DEVELOPMENT.md      # Development guide
│   ├── 🚀 DEPLOYMENT.md       # Deployment guide
│   └── 🤝 CONTRIBUTING.md     # Contributing guidelines
│
├── 🔒 internal/               # Private application code
│   ├── 🌐 api/                # HTTP handlers and middleware
│   │   ├── handlers/          # Request handlers by domain
│   │   ├── middleware/        # Custom middleware
│   │   └── server.go          # Server setup and routing
│   │
│   ├── 📋 cmd/                # Cobra CLI commands
│   │   ├── migrate.go         # Migration commands
│   │   ├── seed.go            # Seeding commands
│   │   └── root.go            # Root command setup
│   │
│   ├── 🗄️ db/                 # Database layer
│   │   ├── query/             # SQL query definitions
│   │   └── repository/        # Generated SQLC code & models
│   │
│   ├── 🔧 utils/              # Internal utility functions
│   └── ⚙️ worker/             # Background job processing
│       ├── processor.go       # Job processor
│       ├── distributor.go     # Job distributor
│       └── tasks/             # Task definitions
│
├── 🔄 migrations/             # Database migration files
│   ├── 000001_initial.up.sql
│   └── 000001_initial.down.sql
│
├── 📦 pkg/                    # Public reusable packages
│   ├── 🔐 auth/               # Authentication utilities
│   │   ├── jwt.go             # JWT token handling
│   │   ├── paseto.go          # PASETO token handling
│   │   └── password.go        # Password hashing
│   │
│   ├── 💾 cachesrv/           # Cache service abstraction
│   ├── 📊 logger/             # Structured logging
│   ├── 📧 mailer/             # Email service
│   ├── 💳 pmgateway/          # Payment gateway integration
│   └── ☁️ upload/             # File upload service
│
├── 🌱 seeds/                  # Database seed data
│   ├── users.json            # Sample user data
│   ├── products.json         # Sample product data
│   └── categories.json       # Sample category data
│
├── 📂 static/                 # Static assets
│   └── templates/             # Email templates
│       ├── verify-email.html
│       └── order-created.html
│
├── 🧪 tests/                  # Test files
│   ├── unit/                 # Unit tests
│   ├── integration/          # Integration tests
│   ├── api/                  # API tests
│   └── fixtures/             # Test fixtures
│
├── 🐳 volumes/                # Docker volume mounts
├── 📄 docker-compose.yml     # Development containers
├── 🐳 Dockerfile            # Container definition
├── 🔧 Makefile              # Build and task automation
├── 🌍 app.env               # Environment configuration
├── 📋 go.mod                # Go module definition
├── 🔧 sqlc.yaml             # SQLC configuration
└── 📖 README.md             # This file
```

### 🏗️ Architecture Overview

The project follows **Clean Architecture** principles with clear separation of concerns:

- **`cmd/`**: Application entry points and CLI tools
- **`internal/api/`**: HTTP layer (handlers, middleware, routing)
- **`internal/db/`**: Data access layer (repositories, models)
- **`pkg/`**: Reusable packages that could be imported by other projects
- **`config/`**: Configuration management and environment variables
- **`migrations/`**: Database schema versioning
- **`tests/`**: Comprehensive test suite with different test types

## 🔧 Development

### 📋 Available Commands

The project includes a comprehensive Makefile for common development tasks:

#### 🗄️ Database Operations
```bash
# Database Migrations
make migrate-up           # Apply all pending migrations
make migrate-down         # Rollback one migration
make migrate-up-1         # Apply one migration
make migrate-drop         # Drop all tables (⚠️ destructive)
make migrate-version      # Show current migration version

# Create new migration
make create-migration name=add_user_preferences

# Database Seeding
make seed                 # Populate database with sample data
```

#### 🏗️ Build & Run
```bash
# Development
make serve-server         # Start development server
make build-server         # Build server binary
make build-migrate        # Build migration tool

# Code Generation
make sqlc                 # Generate Go code from SQL
make swagger              # Generate API documentation
```

#### 🧪 Testing & Quality
```bash
make test                 # Run all tests
make test-coverage        # Run tests with coverage report
make lint                 # Run golangci-lint
make fmt                  # Format Go code
```

#### 💳 External Services
```bash
make listen-stripe        # Listen to Stripe webhooks (development)
```

### 🔧 Development Tools

#### Live Reloading with Air
```bash
# Install Air
go install github.com/air-verse/air@latest

# Start with live reload
air
```

#### Code Generation
```bash
# Generate mocks for testing
go generate ./...

# Update Go dependencies
go mod tidy
go mod verify
```

### 🌍 Environment Variables

Key environment variables for development:

```env
# Server Configuration
ENV=development
PORT=4000
DOMAIN=localhost

# Database
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
MAX_POOL_SIZE=10

# Cache
REDIS_URL=localhost:6380

# Authentication
SYMMETRIC_KEY=your-32-character-secret-key
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h

# External Services (Development)
CLOUDINARY_URL=cloudinary://key:secret@cloud_name
STRIPE_SECRET_KEY=sk_test_your_test_key
SMTP_USERNAME=your_email@example.com
SMTP_PASSWORD=your_app_password
```

### 🚀 Hot Reloading Setup

1. **Install Air (Go live reload)**
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. **Start with hot reload**
   ```bash
   air
   ```

3. **Database changes auto-apply**
   ```bash
   # Monitor migration files and auto-apply
   ls migrations/*.sql | entr make migrate-up
   ```

## 🚀 Deployment

### 🐳 Docker Deployment

#### Quick Docker Setup
```bash
# Build and start all services
docker-compose up -d

# Build and start with rebuild
docker-compose up -d --build

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

#### Production Docker
```bash
# Build production image
docker build -t eshop-api:latest .

# Run with production configuration
docker run -d \
  --name eshop-api \
  -p 4000:4000 \
  --env-file .env.prod \
  eshop-api:latest
```

### ☁️ Cloud Deployment Options

#### AWS (Recommended for scalability)
- **ECS Fargate**: Serverless containers with automatic scaling
- **RDS PostgreSQL**: Managed database with automated backups
- **ElastiCache Redis**: Managed Redis for caching
- **Application Load Balancer**: High availability and SSL termination

#### Google Cloud Platform
- **Cloud Run**: Serverless container deployment
- **Cloud SQL**: Managed PostgreSQL database
- **Memorystore**: Managed Redis instance

#### DigitalOcean (Cost-effective)
- **App Platform**: Simple container deployment
- **Managed Databases**: PostgreSQL and Redis
- **Load Balancers**: Built-in SSL and health checks

#### Heroku (Quick deployment)
- **Heroku Dynos**: Simple git-based deployment
- **Heroku Postgres**: Managed PostgreSQL
- **Heroku Redis**: Managed Redis instance

### 📋 Pre-deployment Checklist

- [ ] Environment variables configured
- [ ] Database migrations applied
- [ ] SSL certificates installed
- [ ] Health checks configured
- [ ] Monitoring and logging setup
- [ ] Backup strategy implemented
- [ ] Load testing completed

For detailed deployment instructions, see **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)**

## 🤝 Contributing

We welcome contributions from the community! Please read our **[Contributing Guide](docs/CONTRIBUTING.md)** for detailed information on:

- 📋 **Development setup** and prerequisites
- 🔧 **Coding standards** and best practices  
- 🧪 **Testing requirements** and guidelines
- 📝 **Documentation standards**
- 🔄 **Pull request process** and review criteria
- 🐛 **Bug reporting** and feature requests

### 🚀 Quick Contribution Steps

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a feature branch** from main
4. **Make your changes** with tests
5. **Submit a pull request** with clear description

### 🌟 Ways to Contribute

- 🐛 **Report bugs** and suggest fixes
- ✨ **Propose new features** and enhancements
- 📖 **Improve documentation** and examples
- 🧪 **Add tests** and improve coverage
- 🔍 **Review pull requests** and provide feedback
- 🎨 **Improve UI/UX** design and user experience

### 📋 Development Workflow

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/go-eshop.git
cd go-eshop/server

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes and test
go test ./...
make lint

# 4. Commit and push
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name

# 5. Create pull request on GitHub
```

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### 🔓 What this means:
- ✅ **Commercial use** - Use it in your commercial projects
- ✅ **Modification** - Modify the code as needed
- ✅ **Distribution** - Share and distribute the code
- ✅ **Private use** - Use it in private projects
- ✅ **No warranty** - No liability for any damages

## 🙏 Acknowledgments

Special thanks to the amazing open-source projects that make this possible:

### 🔧 Core Technologies
- **[Gin Web Framework](https://github.com/gin-gonic/gin)** - Fast HTTP web framework
- **[PostgreSQL](https://www.postgresql.org/)** - Powerful object-relational database
- **[Redis](https://redis.io/)** - In-memory data structure store
- **[SQLC](https://sqlc.dev/)** - Type-safe SQL code generation

### 🔌 Integrations & Services  
- **[Stripe](https://stripe.com/)** - Online payment processing
- **[Cloudinary](https://cloudinary.com/)** - Image and video management
- **[Zerolog](https://github.com/rs/zerolog)** - Fast structured logging

### 🧪 Development Tools
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit
- **[GoMock](https://github.com/golang/mock)** - Mocking framework
- **[Air](https://github.com/air-verse/air)** - Live reloading for Go apps
- **[golangci-lint](https://golangci-lint.run/)** - Go linting tool

---

<div align="center">

### 🌟 If you find this project helpful, please consider giving it a star! ⭐

**Made with ❤️ by [Thanh Phuoc Nguyen](https://github.com/thanhphuocnguyen)**

</div>