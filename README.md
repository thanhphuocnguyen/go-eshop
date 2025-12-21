# 🛍️ Go eShop - Modern E-Commerce Backend Platform

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)
[![Chi Router](https://img.shields.io/badge/Chi-Router-00ADD8?style=for-the-badge&logo=go)](https://go-chi.io/)
[![Docker](https://img.shields.io/badge/Docker-20+-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**A production-ready, scalable e-commerce backend API built with Go** ⚡

*Enterprise-grade microservices architecture with modern development practices*

[🚀 Quick Start](#-quick-start) • 
[📖 API Docs](#-api-documentation) • 
[🏗️ Architecture](#-architecture) • 
[🐳 Deploy](#-deployment) • 
[🤝 Contribute](#-contributing)

</div>

---

## ✨ Why Go eShop?

A **battle-tested, enterprise-ready** e-commerce platform designed for **high performance**, **scalability**, and **maintainability**. Built with modern Go practices and clean architecture principles.

### 🎯 **Perfect For:**
- 🚀 **Startups** needing a robust foundation
- 🏢 **Enterprises** requiring scalable architecture  
- 👨‍💻 **Developers** learning modern Go patterns
- 🛒 **E-commerce** businesses of any size

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

## 🌟 Key Features

<table>
<tr>
<td valign="top" width="33%">

### 🛒 **Customer Experience**
- 🔐 **JWT/PASETO Auth** with refresh tokens
- 🔍 **Advanced Search** with filters & sorting
- 🛍️ **Smart Cart** with persistence
- 💳 **Multi-Payment** (Stripe, PayPal)
- 📦 **Order Tracking** with real-time updates
- ⭐ **Review System** with verification
- 🎯 **Personalization** & recommendations
- 📱 **Mobile-First** API design

</td>
<td valign="top" width="33%">

### 🎛️ **Admin & Management**
- 📊 **Analytics Dashboard** with insights
- 📝 **Product Management** with variants
- 📋 **Order Processing** & fulfillment
- 👥 **User Management** & RBAC
- 💰 **Payment Processing** & refunds
- 🎨 **Content Management** system
- 📈 **Inventory Control** with alerts
- ⚙️ **System Configuration** panel

</td>
<td valign="top" width="33%">

### 🔧 **Technical Excellence**
- 🚀 **High Performance** with caching
- 🔒 **Enterprise Security** standards
- 📱 **API-First** design with OpenAPI
- 🔄 **Background Jobs** with Redis/Asynq
- 📧 **Email Service** with templates  
- ☁️ **Cloud Integration** (Cloudinary)
- 🔍 **Observability** & monitoring
- 🐳 **Containerized** deployment

</td>
</tr>
</table>

## 🏗️ Architecture & Tech Stack

### 🎯 Clean Architecture Overview

```mermaid
graph TB
    A[🌐 HTTP Layer] -->|"requests"| B[🎯 Business Logic]
    B -->|"data ops"| C[🗄️ Data Layer] 
    B -->|"async tasks"| D[🔄 Background Jobs]
    
    A --> A1[Chi Router]
    A --> A2[Middlewares] 
    A --> A3[Handlers]
    
    B --> B1[Services]
    B --> B2[Validators]
    B --> B3[DTOs]
    
    C --> C1[PostgreSQL]
    C --> C2[Redis Cache]
    C --> C3[SQLC ORM]
    
    D --> D1[Asynq Jobs]
    D --> D2[Email Service]
    D --> D3[Notifications]

    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#fff3e0
```

### 🛠️ **Core Technologies**

<div align="center">

| **Component** | **Technology** | **Why We Chose It** |
|---------------|---------------|-------------------|
| 🌐 **Web Framework** | [Chi Router](https://go-chi.io/) | Lightweight, fast, idiomatic HTTP routing |
| 🗄️ **Database** | [PostgreSQL 14+](https://postgresql.org/) | ACID compliance, JSON support, performance |
| ⚡ **Cache** | [Redis 6+](https://redis.io/) | In-memory speed, pub/sub, job queues |
| 🔐 **Auth** | JWT/PASETO | Stateless, secure, token-based authentication |
| 💳 **Payments** | [Stripe](https://stripe.com/) | Global coverage, robust API, webhooks |
| ☁️ **Media** | [Cloudinary](https://cloudinary.com/) | CDN, optimization, transformations |
| 🔄 **Jobs** | [Asynq](https://github.com/hibiken/asynq) | Redis-backed, reliable, retries |
| 📊 **Docs** | [Swagger/OpenAPI](https://swagger.io/) | Interactive API documentation |

</div>

### 🔧 **Development Stack**

```bash
# Backend Core
Go 1.24+              # Modern Go with generics
Chi v5                # HTTP router and middleware  
SQLC                  # Type-safe SQL code generation
golang-migrate        # Database version control
Zerolog               # Structured JSON logging

# Data & Cache
PostgreSQL 14+        # Primary database
Redis 6+              # Cache and session store
Asynq                 # Background job processing

# External APIs
Stripe                # Payment processing
Cloudinary            # Image hosting & CDN
SMTP                  # Email delivery

# Development
Air                   # Live reload
golangci-lint         # Code linting
Testify               # Testing framework
Docker Compose        # Local development
```

## 🚀 Quick Start

> **💡 Get up and running in under 5 minutes!**

### 🎯 **One-Command Setup** (Recommended)

```bash
# Clone and start everything with Docker
git clone https://github.com/thanhphuocnguyen/go-eshop.git
cd go-eshop/server
make dev-setup
```

This command will:
- ✅ Start PostgreSQL & Redis containers
- ✅ Run database migrations  
- ✅ Seed sample data
- ✅ Start the API server
- ✅ Open API docs in browser

### 📋 **Prerequisites**

<table>
<tr>
<td>

**🔧 Required**
- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) (usually pre-installed)

</td>
<td>

**⚡ Optional (for advanced setup)**
- [PostgreSQL 14+](https://postgresql.org/download/)
- [Redis 6+](https://redis.io/docs/getting-started/) 
- [Air](https://github.com/air-verse/air) (live reload)

</td>
</tr>
</table>

### ⚡ **Manual Installation Steps**

<details>
<summary>🔽 <strong>Expand for manual setup instructions</strong></summary>

#### 1️⃣ **Clone Repository**
```bash
git clone https://github.com/thanhphuocnguyen/go-eshop.git
cd go-eshop/server
```

#### 2️⃣ **Install Dependencies** 
```bash
go mod download
go mod tidy
```

#### 3️⃣ **Environment Configuration**
```bash
# Copy environment template
cp app.env.example app.env

# Edit with your settings (optional for development)
nano app.env
```

#### 4️⃣ **Start Infrastructure**
```bash
# Option A: Using Docker (Recommended)
docker-compose up -d postgres redis

# Option B: Local installation
# macOS: brew services start postgresql redis
# Ubuntu: sudo systemctl start postgresql redis
```

#### 5️⃣ **Setup Database**
```bash
# Run migrations
make migrate-up

# Add sample data (optional)
make seed
```

#### 6️⃣ **Start Development Server**
```bash
# Method 1: Using Make
make serve-server

# Method 2: With live reload
air

# Method 3: Direct command
go run ./cmd/web
```

</details>

### 🎉 **Verification**

Your API should now be running! Test with these commands:

```bash
# Health check
curl http://localhost:4000/health

# API documentation
open http://localhost:4000/swagger/index.html

# Test endpoint
curl http://localhost:4000/api/v1/products?limit=5
```

### ⚙️ **Configuration**

<details>
<summary>🔧 <strong>Essential Environment Variables</strong></summary>

Update your `app.env` file:

```env
# 🌐 Server Configuration
ENV=development
PORT=4000
DOMAIN=localhost

# 🗄️ Database 
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
MAX_POOL_SIZE=10

# ⚡ Redis Cache
REDIS_URL=localhost:6380

# 🔐 Authentication (generate secure keys)
SYMMETRIC_KEY=your-32-character-secret-key-here!!
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h

# 💳 Stripe (optional for development)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...

# ☁️ Cloudinary (optional)
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name

# 📧 Email (optional)
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

</details>

## 📚 API Documentation

<div align="center">

### 🌐 **Interactive API Explorer**

**[📖 Swagger UI →](http://localhost:4000/swagger/index.html)**  
*Comprehensive API documentation with live testing*

</div>

### 🔐 **Authentication**

All protected endpoints require a Bearer token in the Authorization header:

```bash
Authorization: Bearer YOUR_JWT_TOKEN
```

### 🚀 **API Examples**

<details>
<summary>🔑 <strong>Authentication Endpoints</strong></summary>

#### **User Registration**
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

#### **User Login**
```bash
curl -X POST http://localhost:4000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePassword123!"
  }'
```

#### **Refresh Token**
```bash
curl -X POST http://localhost:4000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token_here"
  }'
```

</details>

<details>
<summary>🛍️ <strong>Product & Shopping Endpoints</strong></summary>

#### **Get Products** (with filtering)
```bash
curl "http://localhost:4000/api/v1/products?page=1&limit=10&sort=name&order=asc&category=electronics&min_price=10&max_price=1000"
```

#### **Get Product Details**
```bash
curl "http://localhost:4000/api/v1/products/550e8400-e29b-41d4-a716-446655440000"
```

#### **Search Products**
```bash
curl "http://localhost:4000/api/v1/products/search?q=laptop&category=electronics"
```

#### **Add to Cart** (authenticated)
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

#### **Get Cart** (authenticated)
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:4000/api/v1/cart
```

</details>

<details>
<summary>📦 <strong>Order Management</strong></summary>

#### **Create Order** (authenticated)
```bash
curl -X POST http://localhost:4000/api/v1/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address_id": "address-uuid",
    "shipping_method_id": "method-uuid",
    "payment_method": "stripe",
    "payment_data": {
      "stripe_token": "tok_visa"
    }
  }'
```

#### **Get User Orders** (authenticated)
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:4000/api/v1/orders?page=1&limit=10&status=completed"
```

#### **Track Order** (authenticated)
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:4000/api/v1/orders/order-uuid/tracking
```

</details>

### 📄 **Documentation Resources**

| **Resource** | **Description** | **Link** |
|-------------|----------------|----------|
| 📖 **API Reference** | Complete endpoint documentation | [docs/API.md](docs/API.md) |
| 🗄️ **Database Schema** | Database design and relationships | [docs/DATABASE.md](docs/DATABASE.md) |
| 🔧 **Development Guide** | Setup and development workflow | [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) |
| 🚀 **Deployment Guide** | Production deployment instructions | [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) |
| 🤝 **Contributing Guide** | How to contribute to the project | [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) |

### 📊 **Response Format**

All API responses follow a consistent structure:

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Request processed successfully",
  "pagination": {  // For paginated responses
    "page": 1,
    "limit": 10,
    "total": 150,
    "pages": 15
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Email format is invalid"
      }
    ]
  }
}
```

## 🧪 Testing

### � **Test Coverage & Quality**

<div align="center">

![Test Coverage](https://img.shields.io/badge/Coverage-85%25+-green?style=for-the-badge)
![Test Types](https://img.shields.io/badge/Tests-Unit%20%7C%20Integration%20%7C%20API-blue?style=for-the-badge)

</div>

We maintain **high test coverage** across all application layers with comprehensive testing strategies.

### 🚀 **Quick Testing**

```bash
# Run all tests
make test

# Run with coverage report  
make test-coverage

# Run tests with race detection
make test-race

# Run specific test types
make test-unit          # Unit tests only
make test-integration   # Integration tests  
make test-api          # API endpoint tests
```

### 🧪 **Test Categories**

<table>
<tr>
<td valign="top">

#### **🔬 Unit Tests**
```bash
# Run unit tests
go test -short ./...

# Test specific package
go test ./internal/api
go test ./pkg/auth

# With verbose output
go test -v ./internal/utils
```

**Focus Areas:**
- ✅ Business logic functions
- ✅ Utility functions  
- ✅ Data transformations
- ✅ Validation logic

</td>
<td valign="top">

#### **🔗 Integration Tests**  
```bash
# Run integration tests
go test -tags=integration ./...

# Database integration
go test -tags=db ./internal/db

# External service integration  
go test -tags=external ./pkg/...
```

**Focus Areas:**
- ✅ Database operations
- ✅ Redis caching
- ✅ Email service
- ✅ Payment gateways

</td>
</tr>
<tr>
<td valign="top">

#### **🌐 API Tests**
```bash
# Run API tests
go test ./tests/api/...

# Specific endpoint tests
go test ./tests/api/auth
go test ./tests/api/products
```

**Focus Areas:**
- ✅ HTTP endpoints
- ✅ Authentication flows
- ✅ Request/response formats
- ✅ Error handling

</td>
<td valign="top">

#### **⚡ Performance Tests**
```bash
# Benchmark tests
go test -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof ./...

# CPU profiling  
go test -cpuprofile=cpu.prof ./...
```

**Focus Areas:**
- ✅ Response times
- ✅ Memory usage
- ✅ Concurrent requests
- ✅ Database queries

</td>
</tr>
</table>

### 📁 **Test Structure**

```
tests/
├── 🔬 unit/                    # Unit tests for individual functions
│   ├── auth/                  #   Authentication logic tests
│   ├── utils/                 #   Utility function tests  
│   └── validators/            #   Input validation tests
│
├── 🔗 integration/             # Integration tests with external deps
│   ├── database/              #   Database operation tests
│   ├── cache/                 #   Redis cache tests
│   └── email/                 #   Email service tests
│
├── 🌐 api/                     # API endpoint tests  
│   ├── auth/                  #   Authentication endpoint tests
│   ├── products/              #   Product API tests
│   ├── cart/                  #   Shopping cart tests
│   └── orders/                #   Order management tests
│
├── 📊 fixtures/                # Test data and fixtures
│   ├── users.json            #   Sample user data
│   ├── products.json         #   Sample product data
│   └── orders.json           #   Sample order data
│
├── 🎭 mocks/                   # Generated mocks for testing
│   ├── db_mock.go            #   Database mock interfaces
│   ├── cache_mock.go         #   Cache service mocks
│   └── payment_mock.go       #   Payment gateway mocks
│
└── 🔧 helpers/                 # Test utility functions
    ├── test_db.go            #   Test database setup
    ├── test_server.go        #   Test server helpers
    └── assertions.go         #   Custom assertion helpers
```

### 🛠️ **Testing Best Practices**

<details>
<summary>📋 <strong>Testing Guidelines & Standards</strong></summary>

#### **✅ Writing Effective Tests**

1. **Test Structure (AAA Pattern)**
   ```go
   func TestUserRegistration(t *testing.T) {
       // Arrange
       user := &dto.CreateUserRequest{
           Username: "testuser",
           Email:    "test@example.com",
           Password: "SecurePass123!",
       }
       
       // Act
       result, err := userService.CreateUser(context.Background(), user)
       
       // Assert
       assert.NoError(t, err)
       assert.NotNil(t, result)
       assert.Equal(t, user.Email, result.Email)
   }
   ```

2. **Table-Driven Tests**
   ```go
   func TestValidateEmail(t *testing.T) {
       tests := []struct {
           name    string
           email   string
           wantErr bool
       }{
           {"valid email", "user@example.com", false},
           {"invalid format", "invalid-email", true},
           {"empty email", "", true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               err := ValidateEmail(tt.email)
               if (err != nil) != tt.wantErr {
                   t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
               }
           })
       }
   }
   ```

3. **Mock Dependencies**
   ```go
   //go:generate mockgen -source=user.go -destination=mocks/user_mock.go
   
   func TestGetUser(t *testing.T) {
       ctrl := gomock.NewController(t)
       defer ctrl.Finish()
       
       mockDB := mocks.NewMockDatabase(ctrl)
       userService := NewUserService(mockDB)
       
       mockDB.EXPECT().
           GetUser(gomock.Any(), "user123").
           Return(&models.User{ID: "user123"}, nil)
       
       user, err := userService.GetUser(context.Background(), "user123")
       assert.NoError(t, err)
       assert.Equal(t, "user123", user.ID)
   }
   ```

#### **🗄️ Database Testing**
```go
// Test database setup
func setupTestDB(t *testing.T) *sql.DB {
    db := testhelpers.CreateTestDB(t)
    
    // Run migrations
    err := runMigrations(db)
    require.NoError(t, err)
    
    return db
}

// Cleanup after tests
func teardownTestDB(t *testing.T, db *sql.DB) {
    _, err := db.Exec("TRUNCATE TABLE users, products, orders CASCADE")
    require.NoError(t, err)
    db.Close()
}
```

#### **🌐 API Testing**  
```go
func TestLoginEndpoint(t *testing.T) {
    // Setup test server
    server := setupTestServer(t)
    defer server.Close()
    
    // Prepare request
    loginReq := &dto.LoginRequest{
        Username: "testuser",
        Password: "password123",
    }
    
    body, _ := json.Marshal(loginReq)
    
    // Make request
    resp, err := http.Post(
        server.URL+"/api/v1/auth/login",
        "application/json",
        bytes.NewBuffer(body),
    )
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    var response dto.LoginResponse
    err = json.NewDecoder(resp.Body).Decode(&response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response.AccessToken)
}
```

</details>

### 📈 **Coverage Reports**

```bash
# Generate detailed coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Coverage by package
go tool cover -func=coverage.out

# Set coverage threshold
go test -coverprofile=coverage.out -covermode=count ./...
go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//' | awk '{if($1>=80) print "✅ Coverage OK ("$1"%)"; else print "❌ Coverage below threshold ("$1"%)"}'
```

### 🔄 **Continuous Testing**

```bash
# Watch mode for tests (using entr)
find . -name "*.go" | entr -r go test ./...

# Run tests on file changes (using watchdog)
watchmedo shell-command \
    --patterns="*.go" \
    --recursive \
    --command="go test ./..." \
    .
```

## 📁 Project Structure

<details>
<summary>🗂️ <strong>Explore the codebase structure</strong></summary>

```
server/
├── 📁 cmd/                     # 🚀 Application Entry Points
│   ├── migrate/               #    Database migration CLI
│   ├── seed/                  #    Database seeding CLI  
│   └── web/                   #    Main API server
│
├── ⚙️ config/                  # 🔧 Configuration Management
│   └── config.go              #    App configuration & env loading
│
├── 📖 docs/                    # 📚 Documentation Hub
│   ├── docs.go                #    Generated Swagger docs
│   ├── API.md                 #    API documentation
│   ├── DATABASE.md            #    Database schema docs
│   ├── DEVELOPMENT.md         #    Development guide
│   ├── DEPLOYMENT.md          #    Deployment instructions
│   └── CONTRIBUTING.md        #    Contributing guidelines
│
├── 🔐 internal/               # 🏠 Private Application Code
│   ├── 🌐 api/                #    HTTP handlers & routing
│   │   ├── handlers/          #        Domain-specific handlers
│   │   ├── middleware/        #        Custom middleware functions
│   │   ├── routes.go          #        Route definitions
│   │   └── server.go          #        Server setup & config
│   │
│   ├── 🗄️ db/                 #    Data Access Layer
│   │   ├── query/             #        SQL query definitions  
│   │   ├── repository/        #        Generated SQLC repositories
│   │   └── seeds/             #        Database seed data
│   │
│   ├── 📋 dto/                #    Data Transfer Objects
│   │   ├── user.go            #        User-related DTOs
│   │   ├── product.go         #        Product-related DTOs
│   │   └── order.go           #        Order-related DTOs
│   │
│   ├── 🎯 models/             #    Domain Models
│   │   ├── user.go            #        User domain model
│   │   ├── product.go         #        Product domain model
│   │   └── common.go          #        Shared model types
│   │
│   ├── 🔧 utils/              #    Internal Utilities
│   │   ├── helpers.go         #        Common helper functions
│   │   ├── validators.go      #        Input validation
│   │   └── random.go          #        Random data generation
│   │
│   └── ⚙️ worker/             #    Background Jobs
│       ├── processor.go       #        Job processor implementation
│       ├── distributor.go     #        Job distribution logic
│       └── tasks/             #        Task definitions
│
├── 🔄 migrations/             # 📊 Database Migrations
│   ├── 000001_initial.up.sql #    Initial schema creation
│   └── 000001_initial.down.sql#   Schema rollback
│
├── 📦 pkg/                    # 🌍 Public Reusable Packages  
│   ├── 🔐 auth/               #    Authentication utilities
│   │   ├── jwt.go             #        JWT token handling
│   │   ├── paseto.go          #        PASETO token handling  
│   │   └── pwdhasher.go       #        Password hashing
│   │
│   ├── 💾 cache/              #    Cache service abstraction
│   │   ├── cache.go           #        Cache interface
│   │   └── redis.go           #        Redis implementation
│   │
│   ├── 📧 mailer/             #    Email service
│   │   ├── mailer.go          #        Email interface
│   │   └── smtp.go            #        SMTP implementation
│   │
│   ├── 💳 gateways/           #    Payment gateway integration
│   │   ├── stripe.go          #        Stripe integration
│   │   └── paypal.go          #        PayPal integration
│   │
│   ├── 📊 logger/             #    Structured logging
│   │   └── logger.go          #        Zerolog configuration
│   │
│   └── ☁️ upload/             #    File upload service
│       ├── cloudinary.go      #        Cloudinary integration
│       └── local.go           #        Local file storage
│
├── 🌱 seeds/                  # 📋 Sample Data
│   ├── users.json            #    Sample user accounts
│   ├── products.json         #    Sample product catalog  
│   ├── categories.json       #    Sample categories
│   └── brands.json           #    Sample brand data
│
├── 📂 static/                 # 🎨 Static Assets
│   └── templates/             #    Email HTML templates
│       ├── verify-email.html  #        Email verification
│       ├── order-created.html #        Order confirmation
│       └── password-reset.html#        Password reset
│
├── 🧪 tests/                  # ✅ Test Suite
│   ├── unit/                 #    Unit tests
│   ├── integration/          #    Integration tests  
│   ├── api/                  #    API endpoint tests
│   ├── fixtures/             #    Test fixtures & data
│   └── mocks/                #    Generated test mocks
│
├── 🐳 volumes/                # 💾 Docker Volume Mounts
├── 📄 docker-compose.yml     # 🐳 Development containers
├── 🐳 Dockerfile            # 📦 Container definition
├── 🔧 Makefile              # 🛠️ Build & task automation  
├── 🌍 app.env               # ⚙️ Environment configuration
├── 📋 go.mod                # 📦 Go module definition
├── 🔧 sqlc.yaml             # 🗄️ SQLC configuration
└── 📖 README.md             # 📚 Project documentation
```

</details>

### 🏗️ **Architecture Principles**

**Clean Architecture** with clear separation of concerns:

- **`cmd/`** → Application entry points and CLI tools
- **`internal/api/`** → HTTP layer (handlers, middleware, routing)  
- **`internal/db/`** → Data access layer (repositories, models)
- **`pkg/`** → Reusable packages (could be imported by other projects)
- **`config/`** → Configuration management and environment variables
- **`migrations/`** → Database schema versioning
- **`tests/`** → Comprehensive test suite with different test types

### 🎯 **Key Design Patterns**

| **Pattern** | **Implementation** | **Benefits** |
|-------------|-------------------|-------------|
| **Repository** | `internal/db/repository/` | Data access abstraction |
| **Service Layer** | `internal/api/handlers/` | Business logic separation |
| **Dependency Injection** | Constructor pattern | Testability & flexibility |
| **Middleware Chain** | `internal/api/middleware/` | Cross-cutting concerns |
| **Worker Pattern** | `internal/worker/` | Background job processing |

## 🔧 Development

### 📋 **Available Commands**

<div align="center">

| **Category** | **Command** | **Description** |
|-------------|-------------|----------------|
| 🏗️ **Build** | `make build-server` | Build production binary |
| 🏗️ **Build** | `make build-migrate` | Build migration tool |
| 🔄 **Development** | `make serve-server` | Start development server |
| 🔄 **Development** | `make dev-setup` | Complete development setup |
| 🗄️ **Database** | `make migrate-up` | Apply all migrations |
| 🗄️ **Database** | `make migrate-down` | Rollback one migration |
| 🗄️ **Database** | `make seed` | Load sample data |
| 🧪 **Testing** | `make test` | Run all tests |
| 🧪 **Testing** | `make test-coverage` | Run with coverage |
| 🔍 **Code Quality** | `make lint` | Run linting checks |
| 🔍 **Code Quality** | `make fmt` | Format code |
| 📊 **Generate** | `make sqlc` | Generate Go from SQL |
| 📊 **Generate** | `make swagger` | Generate API docs |

</div>

### 🚀 **Development Workflow**

<details>
<summary>🛠️ <strong>Setup Development Environment</strong></summary>

#### **1. Quick Setup**
```bash
# One command to rule them all
make dev-setup
```

#### **2. Manual Setup** 
```bash
# Start infrastructure
docker-compose up -d postgres redis

# Setup database
make migrate-up
make seed

# Start development server
make serve-server
```

#### **3. Live Reloading**
```bash
# Install Air for live reloading
go install github.com/air-verse/air@latest

# Start with hot reload
air

# Or use custom config
air -c .air.toml
```

</details>

<details>
<summary>🗄️ <strong>Database Operations</strong></summary>

#### **Migration Management**
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_user_preferences

# Apply migrations
make migrate-up        # Apply all pending
make migrate-up-1      # Apply one migration
make migrate-version   # Show current version

# Rollback migrations  
make migrate-down      # Rollback one migration
make migrate-down-to 2 # Rollback to version 2
make migrate-drop      # ⚠️ Drop all tables
```

#### **Database Development**
```bash
# Generate Go code from SQL
make sqlc

# Reset database
make db-reset          # Drop, migrate, seed

# Database shell
make db-shell         # Connect to PostgreSQL
make redis-shell      # Connect to Redis
```

</details>

<details>
<summary>🧪 <strong>Testing & Quality Assurance</strong></summary>

#### **Testing Commands**
```bash
# Run tests
make test             # All tests
make test-short       # Unit tests only  
make test-integration # Integration tests
make test-race        # Race condition detection
make test-coverage    # With coverage report

# Specific packages
go test ./internal/api/...
go test ./pkg/auth/...

# With verbose output
go test -v ./...
```

#### **Code Quality**
```bash
# Linting and formatting
make lint            # Run golangci-lint
make fmt             # Format Go code
make vet             # Go vet analysis

# Security scanning
make security-check  # gosec security scanner

# Dependency analysis
go mod tidy         # Clean dependencies
go mod verify       # Verify dependencies
```

</details>

### ⚡ **Development Tools**

#### **🔥 Hot Reloading with Air**
Create `.air.toml` for custom configuration:
```toml
# Air configuration
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/web"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "migrations", "volumes"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta" 
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

#### **🔍 VS Code Configuration**
Create `.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.testFlags": ["-v"],
  "go.testEnvVars": {
    "ENV": "test"
  },
  "go.buildFlags": ["-v"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

### 🌍 **Environment Management**

<details>
<summary>⚙️ <strong>Environment Configurations</strong></summary>

#### **Development (`app.env`)**
```env
ENV=development
PORT=4000
DOMAIN=localhost
DEBUG=true
LOG_LEVEL=debug

# Database
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
MAX_POOL_SIZE=10

# Cache  
REDIS_URL=localhost:6380

# Auth
SYMMETRIC_KEY=development-secret-key-32-chars!!
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h
```

#### **Testing (`app.test.env`)**
```env
ENV=test
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop_test?sslmode=disable
REDIS_URL=localhost:6380
SYMMETRIC_KEY=test-secret-key-for-testing-32!!
```

#### **Production (`app.prod.env`)**
```env
ENV=production  
PORT=4000
DOMAIN=yourdomain.com
DEBUG=false
LOG_LEVEL=info

# Use environment variables for sensitive data
DB_URL=${DATABASE_URL}
REDIS_URL=${REDIS_URL}
SYMMETRIC_KEY=${JWT_SECRET_KEY}

# External services
STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
CLOUDINARY_URL=${CLOUDINARY_URL}
```

</details>

### 🔄 **CI/CD Integration**

<details>
<summary>🚀 <strong>GitHub Actions Workflow</strong></summary>

Create `.github/workflows/ci.yml`:

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: eshop_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
          
      redis:
        image: redis:6-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.24'
        
    - name: Cache dependencies
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
          
    - name: Download dependencies
      run: go mod download
      
    - name: Run migrations
      run: make migrate-up
      env:
        DB_URL: postgresql://postgres:postgres@localhost:5432/eshop_test?sslmode=disable
        
    - name: Run tests
      run: make test-coverage
      env:
        ENV: test
        DB_URL: postgresql://postgres:postgres@localhost:5432/eshop_test?sslmode=disable
        REDIS_URL: localhost:6379
        
    - name: Run linting
      run: make lint
      
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        
    - name: Build binary
      run: make build-server
```

</details>

## 🚀 Deployment

<div align="center">

### 🌐 **Ready for Any Cloud Platform**

[![AWS](https://img.shields.io/badge/AWS-ECS%20%7C%20Fargate-orange?style=for-the-badge&logo=amazonwebservices)](https://aws.amazon.com/)
[![GCP](https://img.shields.io/badge/GCP-Cloud%20Run-blue?style=for-the-badge&logo=googlecloud)](https://cloud.google.com/)
[![Azure](https://img.shields.io/badge/Azure-Container%20Apps-blue?style=for-the-badge&logo=microsoftazure)](https://azure.microsoft.com/)
[![Docker](https://img.shields.io/badge/Docker-Self--Hosted-blue?style=for-the-badge&logo=docker)](https://www.docker.com/)

</div>

### 🐳 **Quick Docker Deployment**

<details>
<summary>🚀 <strong>Local Docker Setup</strong></summary>

#### **Complete Stack Deployment**
```bash
# Start all services (API + Database + Redis)
docker-compose up -d

# Build and restart with latest changes
docker-compose up -d --build

# View real-time logs
docker-compose logs -f api

# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ data loss)
docker-compose down -v
```

#### **Production Docker Build**
```bash
# Build optimized production image
docker build -t go-eshop:latest .

# Run production container
docker run -d \
  --name eshop-api \
  -p 4000:4000 \
  --env-file .env.prod \
  --restart unless-stopped \
  go-eshop:latest

# Health check
curl http://localhost:4000/health
```

#### **Multi-stage Docker Optimization**
```dockerfile
# Our Dockerfile uses multi-stage builds for:
# ✅ Smaller image size (~20MB vs 1GB+)
# ✅ Security (no source code in final image)  
# ✅ Faster deployments
# ✅ Optimized for production
```

</details>

### ☁️ **Cloud Platform Deployments**

<div align="center">

| **Platform** | **Service** | **Best For** | **Deploy Time** |
|-------------|-------------|-------------|----------------|
| 🟠 **AWS** | ECS Fargate | Enterprise scale | 5-10 min |
| 🔵 **Google Cloud** | Cloud Run | Serverless needs | 2-5 min |
| 🟦 **Azure** | Container Apps | Microsoft ecosystem | 5-10 min |
| 🟢 **DigitalOcean** | App Platform | Cost-effective | 3-7 min |
| 🟣 **Heroku** | Containers | Quick prototyping | 1-3 min |

</div>

<details>
<summary>🟠 <strong>AWS Deployment (ECS Fargate)</strong></summary>

#### **Prerequisites**
- AWS CLI configured
- Docker image in ECR
- RDS PostgreSQL instance
- ElastiCache Redis cluster

#### **Deploy Script**
```bash
#!/bin/bash
# deploy-aws.sh

# Build and push image
docker build -t go-eshop .
docker tag go-eshop:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/go-eshop:latest
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/go-eshop:latest

# Update ECS service
aws ecs update-service \
  --cluster eshop-cluster \
  --service eshop-service \
  --force-new-deployment \
  --region $AWS_REGION
```

#### **Infrastructure as Code (Terraform)**
```hcl
# main.tf
resource "aws_ecs_service" "eshop" {
  name            = "eshop-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.eshop.arn
  desired_count   = 2
  
  launch_type = "FARGATE"
  
  network_configuration {
    security_groups = [aws_security_group.eshop.id]
    subnets         = aws_subnet.private[*].id
  }
  
  load_balancer {
    target_group_arn = aws_lb_target_group.eshop.arn
    container_name   = "eshop"
    container_port   = 4000
  }
}
```

</details>

<details>
<summary>🔵 <strong>Google Cloud Run Deployment</strong></summary>

#### **One-Command Deploy**
```bash
# Deploy to Cloud Run
gcloud run deploy eshop-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars ENV=production \
  --set-env-vars PORT=8080 \
  --memory 1Gi \
  --cpu 1000m \
  --max-instances 100
```

#### **Cloud Run YAML Configuration**
```yaml
# cloud-run.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: eshop-api
  annotations:
    run.googleapis.com/ingress: all
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "100"
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containers:
      - image: gcr.io/PROJECT-ID/eshop-api
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: production
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: eshop-secrets
              key: database-url
        resources:
          limits:
            memory: 1Gi
            cpu: 1000m
```

</details>

<details>
<summary>🟢 <strong>DigitalOcean App Platform</strong></summary>

#### **App Specification File**
```yaml
# .do/app.yaml
name: eshop-api
services:
- name: api
  source_dir: /
  github:
    repo: yourusername/go-eshop
    branch: main
  build_command: go build -o bin/server ./cmd/web
  run_command: ./bin/server
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  http_port: 4000
  health_check:
    http_path: /health
  env:
  - key: ENV
    value: production
  - key: PORT
    value: "4000"

databases:
- name: eshop-db
  engine: PG
  version: "14"
  production: true
  
- name: eshop-cache  
  engine: REDIS
  version: "7"
  production: false
```

#### **Deploy Command**
```bash
# Deploy to DigitalOcean
doctl apps create --spec .do/app.yaml

# Update existing app
doctl apps update $APP_ID --spec .do/app.yaml
```

</details>

### 📋 **Pre-deployment Checklist**

<div align="center">

| **Category** | **Task** | **Status** |
|-------------|----------|------------|
| 🔧 **Environment** | Environment variables configured | ✅ |
| 🗄️ **Database** | Migrations applied | ✅ |
| 🔐 **Security** | SSL certificates installed | ✅ |
| 📊 **Monitoring** | Health checks configured | ✅ |
| 📧 **Services** | External services tested | ✅ |
| 💾 **Backup** | Backup strategy implemented | ✅ |
| 🚀 **Performance** | Load testing completed | ✅ |
| 📝 **Documentation** | Deployment docs updated | ✅ |

</div>

### 🔍 **Post-deployment Verification**

```bash
# Health check
curl https://your-domain.com/health

# API functionality test
curl https://your-domain.com/api/v1/products?limit=5

# Database connectivity
curl https://your-domain.com/api/v1/health/db

# Cache connectivity  
curl https://your-domain.com/api/v1/health/cache

# Check API documentation
open https://your-domain.com/swagger/index.html
```

### 📊 **Production Monitoring**

<details>
<summary>📈 <strong>Monitoring & Observability Setup</strong></summary>

#### **Health Check Endpoints**
```bash
GET /health                 # Basic health status
GET /api/v1/health/db      # Database connectivity
GET /api/v1/health/cache   # Redis connectivity  
GET /api/v1/health/deep    # Comprehensive health check
```

#### **Prometheus Metrics** (Optional)
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
- job_name: 'eshop-api'
  static_configs:
  - targets: ['your-api-url:4000']
  metrics_path: /metrics
```

#### **Grafana Dashboard**
- 📊 Request rate & latency
- 🗄️ Database connection pool status  
- ⚡ Cache hit rates
- 🚨 Error rate monitoring
- 📈 Resource utilization

</details>

### 🔒 **Security Considerations**

- **🔐 HTTPS Only**: Force SSL/TLS for all connections
- **🛡️ Rate Limiting**: Implement per-IP rate limiting  
- **🔑 JWT Security**: Use short-lived access tokens
- **📧 Input Validation**: Validate all user inputs
- **🗄️ SQL Injection**: Use parameterized queries (SQLC)
- **🔒 CORS**: Configure appropriate CORS policies
- **📝 Audit Logging**: Log all security-relevant events

For detailed deployment instructions, see **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)**

## 🤝 Contributing

<div align="center">

### 🌟 **We Welcome Contributors!**

*Help us build the best e-commerce platform in Go*

[![Contributors](https://img.shields.io/github/contributors/thanhphuocnguyen/go-eshop?style=for-the-badge)](https://github.com/thanhphuocnguyen/go-eshop/graphs/contributors)
[![Issues](https://img.shields.io/github/issues/thanhphuocnguyen/go-eshop?style=for-the-badge)](https://github.com/thanhphuocnguyen/go-eshop/issues)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge)](https://github.com/thanhphuocnguyen/go-eshop/pulls)

</div>

### 🚀 **Quick Contribution Guide**

1. **🍴 Fork** the repository on GitHub
2. **📥 Clone** your fork locally  
3. **🌿 Create** a feature branch from `main`
4. **💻 Make** your changes with tests
5. **📤 Submit** a pull request with clear description

### 💡 **Ways to Contribute**

<table>
<tr>
<td valign="top" width="25%">

#### **🐛 Bug Reports**
- Found a bug? 
- Report with details
- Include reproduction steps
- Add screenshots if applicable

[🔗 Report Bug →](https://github.com/thanhphuocnguyen/go-eshop/issues/new?template=bug_report.md)

</td>
<td valign="top" width="25%">

#### **✨ Feature Requests**
- Have a cool idea?
- Suggest new features
- Explain the use case
- Discuss implementation approach

[🔗 Request Feature →](https://github.com/thanhphuocnguyen/go-eshop/issues/new?template=feature_request.md)

</td>
<td valign="top" width="25%">

#### **📖 Documentation**
- Improve existing docs
- Add code examples  
- Fix typos & errors
- Translate content

[🔗 Docs Issues →](https://github.com/thanhphuocnguyen/go-eshop/labels/documentation)

</td>
<td valign="top" width="25%">

#### **🧪 Testing**
- Add test cases
- Improve coverage
- Performance testing
- Security testing

[🔗 Testing Issues →](https://github.com/thanhphuocnguyen/go-eshop/labels/testing)

</td>
</tr>
</table>

### 🔧 **Development Workflow**

<details>
<summary>👨‍💻 <strong>Setup Development Environment</strong></summary>

#### **1. Fork & Clone**
```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/go-eshop.git
cd go-eshop/server

# Add upstream remote
git remote add upstream https://github.com/thanhphuocnguyen/go-eshop.git
```

#### **2. Create Feature Branch**  
```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

#### **3. Setup Environment**
```bash
# Install dependencies
go mod download

# Setup development environment  
make dev-setup

# Run tests to ensure everything works
make test
```

</details>

<details>
<summary>📋 <strong>Contribution Standards</strong></summary>

#### **✅ Code Quality Standards**

1. **Go Code Style**
   ```bash
   # Format code
   gofmt -s -w .
   
   # Run linters
   golangci-lint run
   
   # Vet code
   go vet ./...
   ```

2. **Commit Message Convention**
   ```
   type(scope): description
   
   feat(auth): add OAuth2 authentication
   fix(api): resolve product search pagination bug  
   docs(readme): update deployment instructions
   test(user): add user registration test cases
   refactor(db): optimize query performance
   ```

3. **Testing Requirements**
   ```bash
   # All changes must include tests
   go test ./...
   
   # Maintain coverage above 80%
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out | grep total
   ```

#### **📝 Pull Request Guidelines**

**PR Title Format:**
```
[Type] Brief description of changes

feat: Add product recommendation system
fix: Resolve cart checkout validation error
docs: Update API documentation for orders
```

**PR Description Template:**
```markdown
## 📋 Summary
Brief description of what this PR does.

## 🔄 Changes Made
- List of specific changes
- Another change

## 🧪 Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] API tests pass
- [ ] Manual testing completed

## 📖 Documentation
- [ ] Code comments added
- [ ] Documentation updated
- [ ] API docs updated (if applicable)

## 🔗 Related Issues
Closes #123
References #456
```

</details>

### 👥 **Community Guidelines**

#### **🤝 Code of Conduct**
- **Be respectful** and inclusive
- **Be constructive** in feedback
- **Help others** learn and grow
- **Follow** our code of conduct

#### **💬 Communication Channels**

| **Platform** | **Purpose** | **Link** |
|-------------|-------------|----------|
| 💬 **GitHub Issues** | Bug reports & feature requests | [Issues →](https://github.com/thanhphuocnguyen/go-eshop/issues) |
| 🔄 **GitHub Discussions** | General discussions & questions | [Discussions →](https://github.com/thanhphuocnguyen/go-eshop/discussions) |
| 📧 **Email** | Private matters & security issues | [Contact →](mailto:thanhphuocnguyen@example.com) |

### 🏆 **Recognition**

Contributors are recognized in multiple ways:

- **📝 Contributors List** - Listed in README and repository
- **🎖️ Contribution Stats** - GitHub profile contribution graph  
- **🌟 Feature Credits** - Credited in feature announcements
- **📦 Release Notes** - Mentioned in release changelogs

### 🚧 **Current Focus Areas**

<div align="center">

| **Priority** | **Area** | **Description** | **Difficulty** |
|-------------|----------|----------------|---------------|
| 🔥 **High** | Performance | Optimize database queries | 🟡 Medium |
| 🔥 **High** | Security | Add rate limiting middleware | 🟢 Easy |
| 📊 **Medium** | Analytics | Add metrics collection | 🟡 Medium |
| 📱 **Medium** | Mobile API | Mobile-optimized endpoints | 🔴 Hard |
| 🎨 **Low** | UI/UX | Admin dashboard improvements | 🟡 Medium |

</div>

### 📚 **Resources for New Contributors**

- **[Development Guide](docs/DEVELOPMENT.md)** - Detailed setup instructions
- **[API Documentation](docs/API.md)** - Complete API reference  
- **[Database Schema](docs/DATABASE.md)** - Database design overview
- **[Architecture Overview](#-architecture--tech-stack)** - System architecture

### 🎯 **Good First Issues**

New to the project? Start with these beginner-friendly issues:

- 🏷️ [`good first issue`](https://github.com/thanhphuocnguyen/go-eshop/labels/good%20first%20issue) label
- 📖 Documentation improvements
- 🧪 Adding test cases  
- 🐛 Simple bug fixes
- 🔧 Code refactoring

---

<div align="center">

**Thank you for contributing! Every contribution makes this project better. 🙏**

</div>

## 📄 License

<div align="center">

### 🔓 **MIT License** - Freedom to Innovate

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

*This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.*

</div>

### ✅ **What You Can Do**

<div align="center">

| **Permission** | **Description** |
|---------------|----------------|
| ✅ **Commercial Use** | Use in commercial projects and products |
| ✅ **Modification** | Modify the code to suit your needs |
| ✅ **Distribution** | Distribute original or modified versions |
| ✅ **Private Use** | Use privately for personal projects |
| ✅ **Patent Use** | Use any patents covered by the license |

</div>

### 📋 **Requirements**

- 📄 **License Notice** - Include copyright and license notice
- 📝 **Copyright Notice** - Include original copyright notice

### ⚠️ **Limitations**

- 🚫 **No Liability** - Authors not liable for damages
- 🚫 **No Warranty** - No warranties provided

---

## 🙏 Acknowledgments

<div align="center">

### 💝 **Built with Love and Amazing Open Source Projects**

</div>

### 🔧 **Core Technologies**

<table>
<tr>
<td align="center" width="200px">

**🌐 Chi Router**  
*Fast, idiomatic HTTP routing*

[🔗 go-chi.io](https://go-chi.io/)

</td>
<td align="center" width="200px">

**🗄️ PostgreSQL**  
*Advanced open-source database*

[🔗 postgresql.org](https://postgresql.org/)

</td>
<td align="center" width="200px">

**⚡ Redis**  
*In-memory data structure store*

[🔗 redis.io](https://redis.io/)

</td>
<td align="center" width="200px">

**🔧 SQLC**  
*Type-safe SQL in Go*

[🔗 sqlc.dev](https://sqlc.dev/)

</td>
</tr>
</table>

### 🔌 **Integrations & Services**

<div align="center">

| **Service** | **Purpose** | **Why We Love It** |
|------------|-------------|-------------------|
| 💳 **[Stripe](https://stripe.com/)** | Payment processing | Global reach, developer-friendly API |
| ☁️ **[Cloudinary](https://cloudinary.com/)** | Image & video management | Powerful transformations, CDN |
| 📧 **[Asynq](https://github.com/hibiken/asynq)** | Background jobs | Reliable, Redis-based job queue |
| 📊 **[Zerolog](https://github.com/rs/zerolog)** | Structured logging | Fast, zero-allocation logging |

</div>

### 🧪 **Development & Testing**

- **[Testify](https://github.com/stretchr/testify)** - Powerful testing toolkit
- **[Air](https://github.com/air-verse/air)** - Live reload for Go applications  
- **[golangci-lint](https://golangci-lint.run/)** - Comprehensive Go linter
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migrations

### 🎨 **Documentation & APIs**

- **[Swagger/OpenAPI](https://swagger.io/)** - API documentation standard
- **[Mermaid](https://mermaid.js.org/)** - Diagrams and flowcharts  

### 👥 **Community & Inspiration**

Special thanks to:

- **Go Community** - For building an amazing language and ecosystem
- **Open Source Contributors** - For sharing knowledge and code
- **E-commerce Platforms** - For inspiration and best practices
- **API Design Communities** - For REST and GraphQL insights

### 🌟 **Contributors**

<div align="center">

*This project exists thanks to all the people who contribute*

[![Contributors](https://contrib.rocks/image?repo=thanhphuocnguyen/go-eshop)](https://github.com/thanhphuocnguyen/go-eshop/graphs/contributors)

**Want to see your name here? [Contribute to the project!](#-contributing)**

</div>

---

<div align="center">

## 🌟 **Show Your Support**

### If this project helped you, please consider:

[![Star this repo](https://img.shields.io/github/stars/thanhphuocnguyen/go-eshop?style=social)](https://github.com/thanhphuocnguyen/go-eshop/stargazers)
[![Fork this repo](https://img.shields.io/github/forks/thanhphuocnguyen/go-eshop?style=social)](https://github.com/thanhphuocnguyen/go-eshop/network/members)
[![Follow @thanhphuocnguyen](https://img.shields.io/github/followers/thanhphuocnguyen?style=social&label=Follow)](https://github.com/thanhphuocnguyen)

---

### 💬 **Let's Connect**

[![Portfolio](https://img.shields.io/badge/Portfolio-FF5722?style=for-the-badge&logo=todoist&logoColor=white)](https://thanhphuocnguyen.dev)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://linkedin.com/in/thanhphuocnguyen)
[![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/thanhphuocnguyen)
[![Email](https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:thanhphuocnguyen@example.com)

---

**⭐ Star this repository if it helped you build something amazing!**

*Made with ❤️ by [Thanh Phuoc Nguyen](https://github.com/thanhphuocnguyen) and the [community](https://github.com/thanhphuocnguyen/go-eshop/graphs/contributors)*

</div>