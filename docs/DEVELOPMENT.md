# 🛠️ Development Guide

<div align="center">

**Complete development guide for Go eShop platform**

*Everything you need to know to contribute effectively*

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Chi Router](https://img.shields.io/badge/Chi-Router-00ADD8?style=for-the-badge&logo=go)](https://go-chi.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)

[🚀 Quick Setup](#-quick-development-setup) • 
[🏗️ Architecture](#-architecture-deep-dive) • 
[📝 Coding Standards](#-coding-standards--best-practices) • 
[🧪 Testing](#-comprehensive-testing-guide) • 
[🚀 Deployment](#-development--deployment)

</div>

---

## 📋 Table of Contents

- [🚀 Quick Development Setup](#-quick-development-setup)
- [🏗️ Architecture Deep Dive](#-architecture-deep-dive)
- [🔧 Development Workflow](#-development-workflow)
- [📝 Coding Standards & Best Practices](#-coding-standards--best-practices)
- [🗄️ Database Development](#-database-development)
- [🧪 Comprehensive Testing Guide](#-comprehensive-testing-guide)
- [🔒 Security Guidelines](#-security-guidelines)
- [⚡ Performance Optimization](#-performance-optimization)
- [📊 Monitoring & Observability](#-monitoring--observability)
- [🚀 Development & Deployment](#-development--deployment)
- [🤝 Contributing Guidelines](#-contributing-guidelines)

---

## 🚀 Quick Development Setup

### 🎯 **One-Command Development Setup**

```bash
# Clone, setup, and start developing in one command
git clone https://github.com/thanhphuocnguyen/go-eshop.git && \
cd go-eshop/server && \
make dev-setup && \
echo "🎉 Ready to code! Server running at http://localhost:4000"
```

### 📋 **Prerequisites Checklist**

<table>
<tr>
<td valign="top" width="50%">

#### **🔧 Required Tools**
- ✅ **[Go 1.24+](https://golang.org/dl/)** - Modern Go with generics
- ✅ **[Docker](https://docs.docker.com/get-docker/)** - Containerization
- ✅ **[Make](https://www.gnu.org/software/make/)** - Build automation
- ✅ **[Git](https://git-scm.com/)** - Version control

</td>
<td valign="top" width="50%">

#### **⚡ Recommended Tools**  
- 🔥 **[Air](https://github.com/air-verse/air)** - Live reload
- 🔍 **[golangci-lint](https://golangci-lint.run/)** - Linting
- 🐛 **[Delve](https://github.com/go-delve/delve)** - Debugging
- 📊 **[Postman](https://postman.com)** or **[Insomnia](https://insomnia.rest)** - API testing

</td>
</tr>
</table>

### ⚡ **Detailed Setup Steps**

<details>
<summary>🔧 <strong>Manual Development Environment Setup</strong></summary>

#### **1️⃣ Repository Setup**
```bash
# Clone repository
git clone https://github.com/thanhphuocnguyen/go-eshop.git
cd go-eshop/server

# Install Go dependencies
go mod download
go mod tidy

# Verify Go installation
go version  # Should show 1.24+
```

#### **2️⃣ Environment Configuration**
```bash
# Copy environment template
cp app.env.example app.env

# Edit configuration (optional for development)
nano app.env
```

**Essential Development Variables:**
```env
# Server Configuration
ENV=development
PORT=4000
DOMAIN=localhost
DEBUG=true
LOG_LEVEL=debug

# Database (automatically configured by Docker)
DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
MAX_POOL_SIZE=10

# Cache (automatically configured by Docker)  
REDIS_URL=localhost:6380

# Authentication (development keys)
SYMMETRIC_KEY=development-secret-key-32-chars!!
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h
```

#### **3️⃣ Infrastructure Setup**
```bash
# Start PostgreSQL and Redis with Docker
docker-compose up -d postgres redis

# Verify services are running
docker-compose ps

# Check service logs if needed
docker-compose logs postgres redis
```

#### **4️⃣ Database Initialization**
```bash
# Apply database migrations
make migrate-up

# Verify migration status
make migrate-version

# Seed with sample data (recommended for development)
make seed

# Verify data was loaded
make db-shell  # Then run: \dt to see tables
```

#### **5️⃣ Development Server**
```bash
# Option 1: Standard development server
make serve-server

# Option 2: Live reload with Air (recommended)
go install github.com/air-verse/air@latest
air

# Option 3: Debug mode
dlv debug ./cmd/web
```

</details>

### 🎉 **Verification Checklist**

Test your development setup:

```bash
# ✅ Health check
curl http://localhost:4000/health

# ✅ API functionality
curl http://localhost:4000/api/v1/products?limit=5

# ✅ Database connectivity  
curl http://localhost:4000/api/v1/health/db

# ✅ Redis connectivity
curl http://localhost:4000/api/v1/health/cache

# ✅ API documentation
open http://localhost:4000/swagger/index.html
```

**Expected Responses:**
- Health endpoint: `{"status": "healthy"}`
- Database check: `{"database": "connected"}`  
- Redis check: `{"redis": "connected"}`

---

## 🏗️ Architecture Deep Dive

### 🎯 **Clean Architecture Overview**

<details>
<summary>📐 <strong>Architecture Diagram & Principles</strong></summary>

```mermaid
graph TB
    subgraph "🌐 Presentation Layer"
        A[Chi Router] 
        B[Middleware Chain]
        C[HTTP Handlers]
    end
    
    subgraph "🎯 Business Layer"  
        D[Service Layer]
        E[Domain Logic]
        F[Validation]
    end
    
    subgraph "🗄️ Data Layer"
        G[Repository Pattern]
        H[SQLC Generated]
        I[Database Models]
    end
    
    subgraph "🔌 External Layer"
        J[PostgreSQL]
        K[Redis Cache]
        L[Stripe API]
        M[Cloudinary]
        N[Email Service]
    end

    A --> B --> C
    C --> D --> E
    E --> F
    D --> G --> H
    H --> I
    I --> J
    D --> K
    D --> L
    D --> M  
    D --> N

    style A fill:#e3f2fd
    style D fill:#f3e5f5  
    style G fill:#e8f5e8
    style J fill:#fff3e0
```

**🎯 Architecture Principles:**

1. **🔒 Dependency Inversion** - High-level modules don't depend on low-level modules
2. **🔄 Interface Segregation** - Clients shouldn't depend on interfaces they don't use  
3. **🎯 Single Responsibility** - Each component has one reason to change
4. **📐 Separation of Concerns** - Business logic separated from infrastructure

</details>

### 📁 **Project Structure Deep Dive**

<details>
<summary>🗂️ <strong>Directory Structure with Purpose</strong></summary>

```
server/
├── 🚀 cmd/                         # Application Entry Points
│   ├── migrate/                   #   Database migration CLI tool
│   │   └── main.go               #   Migration command implementation
│   ├── seed/                     #   Database seeding CLI tool  
│   │   └── main.go               #   Seed command implementation
│   └── web/                      #   Main API server application
│       └── main.go               #   Server startup and initialization
│
├── ⚙️ config/                      # Configuration Management
│   └── config.go                 #   Environment config loading & validation
│
├── 📚 docs/                        # Documentation Hub
│   ├── docs.go                   #   Swagger documentation generator
│   ├── swagger.json              #   Generated OpenAPI spec (JSON)
│   ├── swagger.yaml              #   Generated OpenAPI spec (YAML)
│   ├── API.md                    #   API endpoint documentation
│   ├── DATABASE.md               #   Database schema & design docs
│   ├── DEPLOYMENT.md             #   Production deployment guide
│   └── CONTRIBUTING.md           #   Contribution guidelines
│
├── 🔐 internal/                    # Private Application Code (Cannot be imported)
│   ├── 🌐 api/                    #   HTTP Layer Implementation
│   │   ├── addresses.go          #     Address management endpoints
│   │   ├── admin*.go             #     Admin panel endpoints (1-4)
│   │   ├── advertises.go         #     Advertisement endpoints
│   │   ├── attributes.go         #     Product attribute endpoints
│   │   ├── auth.go               #     Authentication endpoints
│   │   ├── brands.go             #     Brand management endpoints
│   │   ├── cart.go               #     Shopping cart endpoints
│   │   ├── categories.go         #     Category management endpoints
│   │   ├── checkout.go           #     Checkout process endpoints
│   │   ├── collections.go        #     Product collection endpoints
│   │   ├── discounts.go          #     Discount & coupon endpoints
│   │   ├── images.go             #     Image upload endpoints
│   │   ├── middlewares.go        #     HTTP middleware functions
│   │   ├── orders.go             #     Order management endpoints
│   │   ├── payments.go           #     Payment processing endpoints
│   │   ├── products.go           #     Product catalog endpoints
│   │   ├── ratings.go            #     Product rating endpoints
│   │   ├── routes.go             #     Route definitions & setup
│   │   ├── server.go             #     HTTP server configuration
│   │   ├── users.go              #     User management endpoints
│   │   ├── validators.go         #     Input validation functions
│   │   └── webhooks.go           #     Webhook handlers (Stripe, etc.)
│   │
│   ├── 📋 cmd/                    #   CLI Commands (Cobra)
│   │   ├── migrate.go            #     Migration command definitions
│   │   ├── root.go               #     Root command setup
│   │   └── seed.go               #     Seed command definitions
│   │
│   ├── 🗄️ db/                     #   Database Access Layer
│   │   ├── query/                #     SQL query definitions (SQLC)
│   │   │   ├── users.sql         #       User-related queries
│   │   │   ├── products.sql      #       Product-related queries
│   │   │   ├── orders.sql        #       Order-related queries
│   │   │   └── ...               #       Other domain queries
│   │   ├── repository/           #     Generated SQLC repositories
│   │   │   ├── models.go         #       Database model definitions
│   │   │   ├── querier.go        #       Query interface
│   │   │   ├── users.sql.go      #       Generated user queries
│   │   │   └── ...               #       Other generated files
│   │   └── seeds/                #     Database seed data files
│   │       ├── users.json        #       Sample user data
│   │       ├── products.json     #       Sample product data
│   │       └── ...               #       Other seed data
│   │
│   ├── 📋 dto/                    #   Data Transfer Objects
│   │   ├── address.go            #     Address-related DTOs
│   │   ├── cart.go               #     Shopping cart DTOs
│   │   ├── order.go              #     Order-related DTOs
│   │   ├── product.go            #     Product-related DTOs
│   │   ├── user.go               #     User-related DTOs
│   │   └── common.go             #     Shared DTO definitions
│   │
│   ├── 🎯 models/                 #   Domain Models (Business Logic)
│   │   ├── user.go               #     User domain model
│   │   ├── product.go            #     Product domain model  
│   │   ├── order.go              #     Order domain model
│   │   └── common.go             #     Shared model types
│   │
│   ├── 🔧 utils/                  #   Internal Utility Functions
│   │   ├── helpers.go            #     Common helper functions
│   │   ├── html_template.go      #     HTML template utilities
│   │   ├── pg_helper.go          #     PostgreSQL helpers
│   │   └── random.go             #     Random data generation
│   │
│   └── ⚙️ worker/                 #   Background Job Processing
│       ├── distributor.go        #     Job distribution logic
│       ├── processor.go          #     Job processing implementation
│       ├── scheduler.go          #     Job scheduling logic
│       ├── send_*_task.go        #     Specific task implementations
│       └── task_names.go         #     Task name constants
│
├── 🔄 migrations/                 #   Database Schema Versioning
│   ├── 000001_*.up.sql          #   Initial schema creation
│   ├── 000001_*.down.sql        #   Schema rollback
│   └── ...                      #   Additional migrations
│
├── 📦 pkg/                        #   Public Reusable Packages (Can be imported)
│   ├── 🔐 auth/                   #   Authentication Utilities
│   │   └── pwdhasher.go          #     Password hashing functions
│   │
│   ├── 💾 cache/                  #   Cache Service Abstraction
│   │   ├── cache.go              #     Cache interface definition
│   │   └── redis.go              #     Redis implementation
│   │
│   ├── 🔌 gateways/               #   External Service Gateways
│   │   ├── paypal.go             #     PayPal integration
│   │   └── stripe.go             #     Stripe payment gateway
│   │
│   ├── 📊 logger/                 #   Structured Logging
│   │   └── logger.go             #     Zerolog configuration
│   │
│   ├── 📧 mailer/                 #   Email Service
│   │   └── mailer.go             #     Email sending interface & impl
│   │
│   ├── 💳 payment/                #   Payment Processing
│   │   └── payment.go            #     Payment interface definitions
│   │
│   └── ☁️ upload/                 #   File Upload Service
│       └── upload.go             #     File upload interface & impl
│
├── 🌱 seeds/                      #   Database Seed Data
│   ├── addresses.json            #   Sample address data
│   ├── brands.json               #   Sample brand data
│   ├── categories.json           #   Sample category data
│   ├── products.json             #   Sample product catalog
│   ├── users.json                #   Sample user accounts
│   └── ...                       #   Other sample data
│
├── 📂 static/                     #   Static Assets
│   └── templates/                #   Email HTML Templates
│       ├── verify-email.html     #     Email verification template
│       ├── order-created.html    #     Order confirmation template
│       └── ...                   #     Other email templates
│
├── 🐳 volumes/                    #   Docker Volume Mounts
│   └── mailhog/                  #   Local email testing
│
├── 📄 Configuration Files
│   ├── docker-compose.yml        #   Development containers
│   ├── Dockerfile               #   Container build definition
│   ├── Makefile                 #   Build & task automation
│   ├── app.env                  #   Environment variables
│   ├── go.mod                   #   Go module definition
│   ├── go.sum                   #   Go dependency checksums
│   ├── sqlc.yaml                #   SQLC configuration
│   └── .air.toml                #   Live reload configuration
│
└── 📖 README.md                  #   Project documentation
```

</details>

### 🔗 **Dependency Flow**

```go
// Example: Clean dependency flow
cmd/web/main.go
    ↓
internal/api/server.go (HTTP Layer)
    ↓  
internal/api/users.go (Handler)
    ↓
internal/services/user_service.go (Business Logic) 
    ↓
internal/db/repository (Data Access)
    ↓
PostgreSQL Database
```

---

## 🔧 Development Workflow

### 🛠️ **Essential Development Commands**

<div align="center">

| **Category** | **Command** | **Description** | **Frequency** |
|-------------|-------------|----------------|---------------|
| 🏗️ **Build** | `make build-server` | Build production binary | When deploying |
| 🔄 **Development** | `make serve-server` | Start development server | Daily |
| 🔄 **Development** | `air` | Live reload server | Daily |
| 🗄️ **Database** | `make migrate-up` | Apply database migrations | When schema changes |
| 🗄️ **Database** | `make migrate-down` | Rollback one migration | When fixing migrations |
| 🗄️ **Database** | `make seed` | Load sample data | Initial setup |
| 🧪 **Testing** | `make test` | Run all tests | Before commits |
| 🧪 **Testing** | `make test-coverage` | Run with coverage | Weekly |
| 🔍 **Code Quality** | `make lint` | Run linting checks | Before commits |
| 📊 **Generation** | `make sqlc` | Generate Go from SQL | When queries change |
| 📊 **Generation** | `make swagger` | Generate API docs | When API changes |

</div>

### 🔄 **Development Lifecycle**

<details>
<summary>📋 <strong>Day-to-day Development Workflow</strong></summary>

#### **🌅 Starting Development Session**
```bash
# 1. Pull latest changes
git pull origin main

# 2. Start infrastructure
docker-compose up -d postgres redis

# 3. Apply any new migrations
make migrate-up

# 4. Start development server with live reload
air

# 5. Open API documentation
open http://localhost:4000/swagger/index.html
```

#### **🔄 During Development**
```bash
# Generate code after SQL changes
make sqlc

# Update API docs after handler changes  
make swagger

# Run tests frequently
make test

# Check code quality
make lint && make fmt

# Database operations
make migrate-create name=add_user_preferences  # Create new migration
make migrate-up                                # Apply migrations
make seed                                      # Refresh test data
```

#### **📝 Before Committing**
```bash
# Complete testing suite
make test-coverage

# Code quality checks
make lint
make fmt
make vet

# Build verification
make build-server

# Commit with conventional format
git add .
git commit -m "feat(auth): add OAuth2 social login support"
```

</details>

### 🔧 **Development Tools Setup**

<details>
<summary>🛠️ <strong>Essential Tool Installation & Configuration</strong></summary>

#### **🔥 Air (Live Reload)**
```bash
# Install Air
go install github.com/air-verse/air@latest

# Create .air.toml configuration
cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/web"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "migrations", "volumes"]
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
EOF

# Start with live reload
air
```

#### **🔍 golangci-lint (Code Quality)**
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create .golangci.yml configuration
cat > .golangci.yml << 'EOF'
run:
  timeout: 5m
  tests: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable-all: false
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - misspell
    - nakedret
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - golint
        - gosec
EOF

# Run linting
golangci-lint run
```

#### **🐛 Delve (Debugging)**
```bash
# Install Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the main application
dlv debug ./cmd/web

# Debug with arguments
dlv debug ./cmd/web -- api

# Debug tests
dlv test ./internal/api

# Attach to running process
dlv attach <pid>
```

#### **📊 VS Code Configuration**
```json
// .vscode/settings.json
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
  },
  "go.toolsManagement.checkForUpdates": "local",
  "gopls": {
    "experimentalWorkspaceModule": true
  }
}

// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch API Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/web",
      "env": {
        "ENV": "development"
      },
## 📋 Contributing Guidelines

### 🤝 **How to Contribute**

<details>
<summary>🚀 <strong>Getting Started</strong></summary>

#### **1. 🍴 Fork & Clone**

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/eshop-server.git
cd eshop-server

# Add upstream remote
git remote add upstream https://github.com/thanhphuocnguyen/eshop-server.git

# Verify remotes
git remote -v
```

#### **2. 🔄 Setup Development Environment**

```bash
# Install dependencies
go mod download

# Setup pre-commit hooks
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/vektra/mockery/v2@latest

# Copy environment configuration
cp app.env.example app.env
# Edit app.env with your local settings

# Run database migrations
make migrate-up

# Seed development data
make seed
```

#### **3. 🌿 Create Feature Branch**

```bash
# Create and switch to feature branch
git checkout -b feature/user-authentication

# Or for bug fixes
git checkout -b fix/order-calculation-bug

# Or for documentation
git checkout -b docs/api-documentation
```

</details>

### 📝 **Development Workflow**

<details>
<summary>⚡ <strong>Coding & Testing Process</strong></summary>

#### **💻 Development Cycle**

```bash
# 1. Start development server
make dev

# 2. Run tests continuously during development
make test-watch

# 3. Check code quality before committing
make lint
make test
make test-coverage

# 4. Format code
make fmt

# 5. Generate mocks (if needed)
make mock

# 6. Update API documentation
make docs
```

#### **🧪 Testing Requirements**

All contributions must include appropriate tests:

```go
// ✅ Unit test example
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        setup   func(*testing.T, *UserService)
        want    *User
        wantErr bool
    }{
        {
            name: "successful user creation",
            input: CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "SecurePass123!",
            },
            setup: func(t *testing.T, s *UserService) {
                s.mockRepo.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            want: &User{
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockRepo := mocks.NewMockUserRepository(t)
            service := NewUserService(mockRepo, nil, nil)
            
            if tt.setup != nil {
                tt.setup(t, service)
            }

            // Execute
            got, err := service.CreateUser(context.Background(), tt.input)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want.Username, got.Username)
            assert.Equal(t, tt.want.Email, got.Email)
        })
    }
}
```

#### **🔍 Code Review Checklist**

Before submitting your PR, ensure:

- [ ] **Tests**: All new code has unit tests with >80% coverage
- [ ] **Documentation**: Public APIs are documented with clear examples
- [ ] **Linting**: Code passes all linting checks (`make lint`)
- [ ] **Security**: No secrets or sensitive data in code
- [ ] **Performance**: No obvious performance issues introduced
- [ ] **Error Handling**: Proper error handling and logging
- [ ] **Breaking Changes**: Clearly documented if any

</details>

### 📨 **Pull Request Process**

<details>
<summary>🔄 <strong>PR Guidelines & Standards</strong></summary>

#### **📝 PR Template**

```markdown
## 🎯 Description

Brief description of what this PR does.

## 🚀 Type of Change

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Code refactoring
- [ ] ⚡ Performance improvement

## 🧪 Testing

- [ ] All existing tests pass
- [ ] New tests added for new functionality
- [ ] Manual testing completed

## 📝 Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Code is properly commented
- [ ] Documentation updated (if applicable)
- [ ] No breaking changes (or clearly documented)

## 📸 Screenshots (if applicable)

Add screenshots for UI changes.

## 🔗 Related Issues

Closes #[issue number]
```

#### **📏 PR Size Guidelines**

- **Small PRs** (< 200 lines): Preferred for faster review
- **Medium PRs** (200-500 lines): Acceptable with good description
- **Large PRs** (> 500 lines): Should be split unless absolutely necessary

#### **🎯 Commit Message Convention**

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Format: <type>(<scope>): <description>

# Types:
feat: add user authentication system
fix: resolve order calculation bug
docs: update API documentation
style: format code according to standards
refactor: restructure user service
test: add unit tests for order service
chore: update dependencies

# Examples:
git commit -m "feat(auth): implement JWT token validation"
git commit -m "fix(orders): correct tax calculation logic"
git commit -m "docs(api): add authentication endpoints documentation"
```

</details>

### 🐛 **Issue Reporting**

<details>
<summary>📋 <strong>Bug Reports & Feature Requests</strong></summary>

#### **🐛 Bug Report Template**

```markdown
## 🐛 Bug Report

### 📝 Description
A clear description of what the bug is.

### 🔄 Steps to Reproduce
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

### ✅ Expected Behavior
What you expected to happen.

### ❌ Actual Behavior
What actually happened.

### 📸 Screenshots
If applicable, add screenshots to help explain your problem.

### 🖥️ Environment
- OS: [e.g., macOS, Ubuntu]
- Go version: [e.g., 1.24]
- Project version: [e.g., v1.2.3]

### 📋 Additional Context
Add any other context about the problem here.
```

#### **✨ Feature Request Template**

```markdown
## ✨ Feature Request

### 🎯 Problem Statement
Is your feature request related to a problem? Please describe.

### 💡 Proposed Solution
Describe the solution you'd like.

### 🔄 Alternatives Considered
Describe alternatives you've considered.

### 📝 Additional Context
Add any other context about the feature request here.

### ✅ Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
```

</details>

---

## 📚 Additional Resources

### 🔗 **Useful Links**

- 📖 **[Go Documentation](https://golang.org/doc/)**
- 🎯 **[Chi Router Documentation](https://github.com/go-chi/chi)**
- 🗄️ **[SQLC Documentation](https://docs.sqlc.dev/)**
- 🐘 **[PostgreSQL Documentation](https://www.postgresql.org/docs/)**
- 🔴 **[Redis Documentation](https://redis.io/documentation)**
- 🐳 **[Docker Documentation](https://docs.docker.com/)**

### 📝 **Project-Specific Guides**

- **[API Documentation](API.md)** - Complete API reference
- **[Database Schema](DATABASE.md)** - Database design and migrations
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment instructions
- **[Contributing Guidelines](CONTRIBUTING.md)** - How to contribute to the project

### 🚀 **Development Tools**

- **[Air](https://github.com/cosmtrek/air)** - Live reload for Go development
- **[golangci-lint](https://golangci-lint.run/)** - Go linters aggregator
- **[Mockery](https://github.com/vektra/mockery)** - Mock generation tool
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit

---

## 🎯 Quick Start Checklist

When setting up the development environment, complete these steps:

- [ ] **Clone repository and install dependencies**
- [ ] **Setup PostgreSQL and Redis locally**
- [ ] **Copy and configure environment variables**
- [ ] **Run database migrations**
- [ ] **Seed development data**
- [ ] **Start development server with Air**
- [ ] **Run the test suite**
- [ ] **Explore the API with Swagger UI**

---

<div align="center">

**🎉 Happy Coding! 🎉**

*For questions or support, reach out via [GitHub Issues](../../issues) or contact the development team.*

</div>

## 📝 Coding Standards & Best Practices

### 🎯 **Go Style Guidelines**

We follow **[Go Style Guide](https://github.com/golang/go/wiki/CodeReviewComments)** and **[Effective Go](https://golang.org/doc/effective_go.html)** with additional project-specific conventions.

<details>
<summary>🔤 <strong>Naming Conventions & Code Style</strong></summary>

#### **📛 Naming Standards**

```go
// ✅ Constants: SCREAMING_SNAKE_CASE or PascalCase
const (
    DefaultPageSize     = 10
    MaxPageSize        = 100
    APIVersionPrefix   = "/api/v1"
    USER_ROLE_ADMIN    = "admin"
)

// ✅ Variables: camelCase for private, PascalCase for exported
var (
    httpClient     *http.Client           // private
    dbPool        *pgxpool.Pool          // private
    DefaultConfig  config.Config         // exported
)

// ✅ Functions: PascalCase (exported) or camelCase (private)
func CreateUser(ctx context.Context, user User) error      // exported
func validateUserInput(user User) error                    // private
func (s *Service) ProcessOrder(ctx context.Context) error  // method

// ✅ Structs & Interfaces: PascalCase
type UserService struct {               // struct
    repo       repository.UserRepository
    cache      cache.Cache
    validator  *validator.Validate
}

type UserRepository interface {         // interface
    CreateUser(ctx context.Context, user User) error
    GetUser(ctx context.Context, id uuid.UUID) (*User, error)
}

// ✅ Package naming: lowercase, single word
package auth      // ✅ Good
package userauth  // ✅ Acceptable  
package user_auth // ❌ Avoid underscores
package UserAuth  // ❌ Avoid capital letters
```

#### **📝 Code Organization**

```go
// ✅ File structure order
package main

// 1. Imports (grouped and sorted)
import (
    // Standard library
    "context"
    "fmt"
    "log"
    
    // Third party
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    
    // Local packages
    "github.com/yourorg/project/internal/config"
    "github.com/yourorg/project/pkg/auth"
)

// 2. Constants
const (
    DefaultTimeout = 30 * time.Second
)

// 3. Variables  
var (
    ErrUserNotFound = errors.New("user not found")
)

// 4. Types
type User struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
}

// 5. Functions (exported first, then private)
func NewUserService() *UserService { ... }
func (s *UserService) CreateUser() error { ... }
func validateUser() error { ... }
```

#### **🔧 Function Design**

```go
// ✅ Good function design
func (s *UserService) CreateUser(
    ctx context.Context,           // Context first
    req CreateUserRequest,         // Input parameters
) (*User, error) {                // Return values (value, error)
    // Validation
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Business logic
    hashedPassword, err := s.hashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("password hashing failed: %w", err)
    }
    
    // Database operation
    user, err := s.repo.CreateUser(ctx, CreateUserParams{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
    })
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    return &user, nil
}

// ✅ Error handling pattern
func (s *Service) ProcessPayment(ctx context.Context, amount float64) error {
    // Input validation
    if amount <= 0 {
        return ErrInvalidAmount
    }
    
    // External service call with timeout
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    result, err := s.paymentGateway.Charge(ctx, amount)
    if err != nil {
        // Wrap errors with context
        return fmt.Errorf("payment gateway error: %w", err)
    }
    
    // Log successful operation
    log.Info().
        Float64("amount", amount).
        Str("transaction_id", result.ID).
        Msg("Payment processed successfully")
        
    return nil
}
```

</details>

### 🔒 **Error Handling Standards**

<details>
<summary>⚠️ <strong>Error Handling Best Practices</strong></summary>

#### **🎯 Error Types & Conventions**

```go
// ✅ Domain-specific error types
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// ✅ Sentinel errors for common cases
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUnauthorized     = errors.New("unauthorized")
    ErrForbidden       = errors.New("forbidden")
)

// ✅ Error wrapping with context
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := s.repo.GetUser(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}

// ✅ HTTP error responses
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        h.sendError(w, http.StatusBadRequest, "INVALID_USER_ID", 
            "Invalid user ID format")
        return
    }
    
    user, err := h.userService.GetUser(r.Context(), userID)
    if err != nil {
        switch {
        case errors.Is(err, ErrUserNotFound):
            h.sendError(w, http.StatusNotFound, "USER_NOT_FOUND", 
                "User not found")
        default:
            h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", 
                "Internal server error")
            // Log the actual error for debugging
            log.Error().Err(err).Str("user_id", userID.String()).
                Msg("Failed to get user")
        }
        return
    }
    
    h.sendJSON(w, http.StatusOK, user)
}

// ✅ Error response helper
func (h *Handler) sendError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    errorResp := ErrorResponse{
        Success: false,
        Error: ErrorDetail{
            Code:    code,
            Message: message,
        },
        Timestamp: time.Now(),
    }
    
    json.NewEncoder(w).Encode(errorResp)
}
```

#### **📊 Error Response Format**

```go
// Consistent error response structure
type ErrorResponse struct {
    Success   bool        `json:"success"`
    Error     ErrorDetail `json:"error"`
    Timestamp time.Time   `json:"timestamp"`
}

type ErrorDetail struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}

// Example error responses
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Input validation failed",
        "details": {
            "email": "Invalid email format",
            "password": "Password must be at least 8 characters"
        }
    },
    "timestamp": "2024-01-15T10:30:00Z"
}
```

</details>

### 🌐 **Chi Router Patterns**

<details>
<summary>🛣️ <strong>Chi Router Implementation & Best Practices</strong></summary>

#### **🔧 Router Setup & Configuration**

```go
// ✅ Main router setup
func (s *Server) setupRoutes() http.Handler {
    r := chi.NewRouter()
    
    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(60 * time.Second))
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))
    
    // API versioning
    r.Route("/api/v1", func(r chi.Router) {
        // Public routes
        r.Group(func(r chi.Router) {
            r.Post("/auth/login", s.LoginHandler)
            r.Post("/auth/register", s.RegisterHandler)
            r.Get("/products", s.GetProductsHandler)
            r.Get("/products/{id}", s.GetProductHandler)
        })
        
        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(s.AuthenticateMiddleware)
            r.Get("/users/profile", s.GetProfileHandler)
            r.Put("/users/profile", s.UpdateProfileHandler)
            r.Post("/cart/items", s.AddToCartHandler)
            r.Get("/orders", s.GetOrdersHandler)
        })
        
        // Admin routes
        r.Group(func(r chi.Router) {
            r.Use(s.AuthenticateMiddleware)
            r.Use(s.AuthorizeMiddleware("admin"))
            r.Post("/admin/products", s.CreateProductHandler)
            r.Put("/admin/products/{id}", s.UpdateProductHandler)
            r.Delete("/admin/products/{id}", s.DeleteProductHandler)
        })
    })
    
    // Health checks
    r.Get("/health", s.HealthCheckHandler)
    r.Get("/health/db", s.DatabaseHealthHandler)
    r.Get("/health/cache", s.CacheHealthHandler)
    
    // Static files and documentation
    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"),
    ))
    
    return r
}

// ✅ RESTful resource routing
func (s *Server) setupProductRoutes(r chi.Router) {
    r.Route("/products", func(r chi.Router) {
        r.Get("/", s.GetProductsHandler)        // GET /products
        r.Post("/", s.CreateProductHandler)     // POST /products
        
        r.Route("/{id}", func(r chi.Router) {
            r.Use(s.ProductContextMiddleware)   // Load product into context
            r.Get("/", s.GetProductHandler)     // GET /products/{id}
            r.Put("/", s.UpdateProductHandler)  // PUT /products/{id}
            r.Delete("/", s.DeleteProductHandler) // DELETE /products/{id}
            
            // Nested resources
            r.Route("/images", func(r chi.Router) {
                r.Post("/", s.UploadProductImageHandler)    // POST /products/{id}/images
                r.Delete("/{imageId}", s.DeleteProductImageHandler) // DELETE /products/{id}/images/{imageId}
            })
            
            r.Route("/reviews", func(r chi.Router) {
                r.Get("/", s.GetProductReviewsHandler)      // GET /products/{id}/reviews
                r.Post("/", s.CreateProductReviewHandler)   // POST /products/{id}/reviews
            })
        })
        
        // Special endpoints
        r.Get("/search", s.SearchProductsHandler)           // GET /products/search
        r.Get("/categories/{categoryId}", s.GetProductsByCategoryHandler) // GET /products/categories/{categoryId}
    })
}
```

#### **🔧 Middleware Implementation**

```go
// ✅ Authentication middleware
func (s *Server) AuthenticateMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            s.sendError(w, http.StatusUnauthorized, "MISSING_AUTH", 
                "Authorization header required")
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            s.sendError(w, http.StatusUnauthorized, "INVALID_AUTH_FORMAT", 
                "Invalid authorization format")
            return
        }
        
        payload, err := s.tokenMaker.VerifyToken(tokenString)
        if err != nil {
            s.sendError(w, http.StatusUnauthorized, "INVALID_TOKEN", 
                "Invalid or expired token")
            return
        }
        
        // Add user info to context
        ctx := context.WithValue(r.Context(), "user", payload)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// ✅ Authorization middleware
func (s *Server) AuthorizeMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := r.Context().Value("user").(*auth.Payload)
            if !ok {
                s.sendError(w, http.StatusUnauthorized, "NO_USER_CONTEXT", 
                    "User not found in context")
                return
            }
            
            for _, role := range allowedRoles {
                if user.Role == role {
                    next.ServeHTTP(w, r)
                    return
                }
            }
            
            s.sendError(w, http.StatusForbidden, "INSUFFICIENT_PERMISSIONS", 
                "Access denied")
        })
    }
}

// ✅ Request validation middleware
func (s *Server) ValidateJSONMiddleware(v interface{}) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Header.Get("Content-Type") != "application/json" {
                s.sendError(w, http.StatusBadRequest, "INVALID_CONTENT_TYPE", 
                    "Content-Type must be application/json")
                return
            }
            
            decoder := json.NewDecoder(r.Body)
            decoder.DisallowUnknownFields()
            
            if err := decoder.Decode(v); err != nil {
                s.sendError(w, http.StatusBadRequest, "INVALID_JSON", 
                    "Invalid JSON format")
                return
            }
            
            if err := s.validator.Struct(v); err != nil {
                s.sendValidationError(w, err)
                return
            }
            
            // Add validated data to context
            ctx := context.WithValue(r.Context(), "requestData", v)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

#### **📝 Handler Implementation**

```go
// ✅ Handler pattern with proper error handling
func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse and validate request
    var req dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        s.sendError(w, http.StatusBadRequest, "INVALID_JSON", 
            "Invalid JSON format")
        return
    }
    
    if err := s.validator.Struct(req); err != nil {
        s.sendValidationError(w, err)
        return
    }
    
    // Business logic
    user, err := s.userService.CreateUser(r.Context(), req)
    if err != nil {
        switch {
        case errors.Is(err, ErrUserAlreadyExists):
            s.sendError(w, http.StatusConflict, "USER_EXISTS", 
                "User already exists")
        default:
            log.Error().Err(err).Msg("Failed to create user")
            s.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", 
                "Failed to create user")
        }
        return
    }
    
    // Success response
    response := dto.UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        // ... other fields
    }
    
    s.sendJSON(w, http.StatusCreated, map[string]interface{}{
        "success": true,
        "data":    response,
        "message": "User created successfully",
    })
}

// ✅ URL parameter handling
func (s *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    userIDStr := chi.URLParam(r, "id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        s.sendError(w, http.StatusBadRequest, "INVALID_USER_ID", 
            "Invalid user ID format")
        return
    }
    
    user, err := s.userService.GetUser(r.Context(), userID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            s.sendError(w, http.StatusNotFound, "USER_NOT_FOUND", 
                "User not found")
            return
        }
        log.Error().Err(err).Str("user_id", userID.String()).
            Msg("Failed to get user")
        s.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", 
            "Failed to retrieve user")
        return
    }
    
    s.sendJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "data":    user,
    })
}

// ✅ Query parameter handling with pagination
func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    search := r.URL.Query().Get("search")
    role := r.URL.Query().Get("role")
    
    // Build filter
    filter := dto.UserFilter{
        Search: search,
        Role:   role,
        Page:   page,
        Limit:  limit,
    }
    
    users, total, err := s.userService.GetUsers(r.Context(), filter)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get users")
        s.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", 
            "Failed to retrieve users")
        return
    }
    
    // Calculate pagination
    totalPages := (total + limit - 1) / limit
    
    s.sendJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "data":    users,
        "pagination": map[string]interface{}{
            "page":        page,
            "limit":       limit,
            "total":       total,
            "totalPages":  totalPages,
            "hasNext":     page < totalPages,
            "hasPrev":     page > 1,
        },
    })
}
```

</details>

### 🗄️ **Database Patterns with SQLC**

<details>
<summary>💾 <strong>SQLC Integration & SQL Best Practices</strong></summary>

#### **📝 SQLC Query Development**

```sql
-- internal/db/query/users.sql

-- name: CreateUser :one
INSERT INTO users (
    role_id,
    username, 
    email,
    phone_number,
    first_name,
    last_name,
    hashed_password,
    avatar_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByUsername :one  
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateUser :one
UPDATE users 
SET 
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    email = COALESCE(sqlc.narg('email'), email),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE 
    deleted_at IS NULL
    AND ($1::text IS NULL OR username ILIKE '%' || $1 || '%' OR first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%')
    AND ($2::uuid IS NULL OR role_id = $2)
    AND ($3::bool IS NULL OR email_verified = $3)
ORDER BY 
    CASE WHEN $4 = 'name' THEN first_name END ASC,
    CASE WHEN $4 = 'email' THEN email END ASC,
    CASE WHEN $4 = 'created' THEN created_at END DESC,
    created_at DESC
LIMIT $5 OFFSET $6;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE 
    deleted_at IS NULL
    AND ($1::text IS NULL OR username ILIKE '%' || $1 || '%' OR first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%')
    AND ($2::uuid IS NULL OR role_id = $2)
    AND ($3::bool IS NULL OR email_verified = $3);

-- name: GetUsersWithProfiles :many
SELECT 
    u.id,
    u.username,
    u.email,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.created_at,
    r.name as role_name,
    COUNT(o.id) as order_count
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, r.name
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;
```

#### **🔧 Repository Implementation**

```go
// ✅ Repository interface
type UserRepository interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    GetUser(ctx context.Context, id uuid.UUID) (User, error)
    GetUserByEmail(ctx context.Context, email string) (User, error)
    UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
    ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
    CountUsers(ctx context.Context, arg CountUsersParams) (int64, error)
    SoftDeleteUser(ctx context.Context, id uuid.UUID) error
}

// ✅ Service layer implementation
type UserService struct {
    repo      UserRepository
    cache     cache.Cache
    validator *validator.Validate
}

func NewUserService(repo UserRepository, cache cache.Cache) *UserService {
    return &UserService{
        repo:      repo,
        cache:     cache,
        validator: validator.New(),
    }
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Validation
    if err := s.validator.Struct(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Check if user already exists
    _, err := s.repo.GetUserByEmail(ctx, req.Email)
    if err == nil {
        return nil, ErrUserAlreadyExists
    }
    if !errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user in database
    user, err := s.repo.CreateUser(ctx, CreateUserParams{
        RoleID:         uuid.MustParse("default-role-id"),
        Username:       req.Username,
        Email:          req.Email,
        PhoneNumber:    req.PhoneNumber,
        FirstName:      req.FirstName,
        LastName:       req.LastName,
        HashedPassword: string(hashedPassword),
        AvatarUrl:      sql.NullString{},
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Convert to response DTO
    response := &dto.UserResponse{
        ID:          user.ID,
        Username:    user.Username,
        Email:       user.Email,
        PhoneNumber: user.PhoneNumber,
        FirstName:   user.FirstName,
        LastName:    user.LastName,
        AvatarURL:   user.AvatarUrl.String,
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
    }
    
    return response, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%s", id.String())
    if cached := s.cache.Get(ctx, cacheKey); cached != nil {
        var user dto.UserResponse
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // Get from database
    user, err := s.repo.GetUser(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    response := &dto.UserResponse{
        ID:          user.ID,
        Username:    user.Username,
        Email:       user.Email,
        PhoneNumber: user.PhoneNumber,
        FirstName:   user.FirstName,
        LastName:    user.LastName,
        AvatarURL:   user.AvatarUrl.String,
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
    }
    
    // Cache the result
    if data, err := json.Marshal(response); err == nil {
        s.cache.Set(ctx, cacheKey, string(data), time.Hour)
    }
    
    return response, nil
}
```

#### **🔄 Transaction Management**

```go
// ✅ Transaction wrapper
type TxRepository struct {
    *Queries
    db *sql.DB
}

func (r *Repository) WithTx(tx *sql.Tx) *TxRepository {
    return &TxRepository{
        Queries: r.Queries.WithTx(tx),
        db:      r.db,
    }
}

// ✅ Service method with transaction
func (s *UserService) CreateUserWithProfile(
    ctx context.Context, 
    userReq dto.CreateUserRequest,
    profileReq dto.CreateProfileRequest,
) (*dto.UserResponse, error) {
    // Start transaction
    tx, err := s.repo.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer tx.Rollback() // Will be no-op if tx.Commit() is called
    
    // Get repository with transaction
    txRepo := s.repo.WithTx(tx)
    
    // Create user
    user, err := txRepo.CreateUser(ctx, CreateUserParams{
        // ... user parameters
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Create profile
    profile, err := txRepo.CreateProfile(ctx, CreateProfileParams{
        UserID: user.ID,
        Bio:    profileReq.Bio,
        // ... other profile parameters
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create profile: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return &dto.UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Profile:  &dto.ProfileResponse{
            Bio: profile.Bio,
            // ... other fields
        },
    }, nil
}
```

</details>

---
```

## 🗄️ Database Development

### 📊 **Database Design Principles**

<details>
<summary>🏗️ <strong>Schema Design & Best Practices</strong></summary>

#### **🎯 Table Design Standards**

```sql
-- ✅ Good table design
CREATE TABLE users (
    -- Primary key: UUID for distributed systems
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign keys: Explicit naming
    role_id UUID NOT NULL REFERENCES roles(id),
    
    -- Columns: snake_case, descriptive names
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    
    -- Optional fields: Use NULL appropriately
    avatar_url TEXT,
    bio TEXT,
    date_of_birth DATE,
    
    -- Metadata: Always include timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ, -- Soft delete support
    
    -- Constraints
    CONSTRAINT users_username_length CHECK (LENGTH(username) >= 3),
    CONSTRAINT users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT users_phone_format CHECK (phone_number ~* '^\+?[1-9]\d{1,14}$')
);

-- ✅ Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at); -- For soft delete queries

-- ✅ Composite indexes for common queries
CREATE INDEX idx_users_role_active ON users(role_id, deleted_at) 
WHERE deleted_at IS NULL;
```

#### **🔄 Migration Best Practices**

```sql
-- ✅ Migration file: 000002_add_user_preferences.up.sql
-- Add column with default value (safe for large tables)
ALTER TABLE users 
ADD COLUMN preferences JSONB DEFAULT '{}' NOT NULL;

-- Add index for JSONB queries
CREATE INDEX idx_users_preferences ON users USING GIN(preferences);

-- Update existing rows (if needed)
UPDATE users 
SET preferences = '{
    "notifications": {
        "email": true,
        "sms": false,
        "push": true
    },
    "privacy": {
        "profile_visibility": "public"
    }
}'::jsonb 
WHERE preferences = '{}'::jsonb;

-- Add constraint after data migration
ALTER TABLE users 
ADD CONSTRAINT users_preferences_valid 
CHECK (jsonb_typeof(preferences) = 'object');

-- ✅ Migration file: 000002_add_user_preferences.down.sql
-- Remove constraint first
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferences_valid;

-- Remove index
DROP INDEX IF EXISTS idx_users_preferences;

-- Remove column
ALTER TABLE users DROP COLUMN IF EXISTS preferences;
```

#### **📋 Relationship Design**

```sql
-- ✅ One-to-many relationship
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    status order_status_enum NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    currency_code CHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ✅ Many-to-many relationship with junction table
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);

-- ✅ Self-referencing relationship (categories)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES categories(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for hierarchical queries
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
```

</details>

### 🔄 **Migration Management**

<details>
<summary>📝 <strong>Database Migration Workflow</strong></summary>

#### **🛠️ Creating Migrations**

```bash
# Create new migration
make migrate-create name=add_product_variants

# This creates two files:
# migrations/000003_add_product_variants.up.sql
# migrations/000003_add_product_variants.down.sql

# Migration naming conventions:
make migrate-create name=create_users_table        # Initial table creation
make migrate-create name=add_email_to_users        # Adding columns
make migrate-create name=remove_legacy_fields      # Removing columns  
make migrate-create name=create_orders_index       # Adding indexes
make migrate-create name=update_user_constraints   # Constraint changes
```

#### **📝 Migration File Structure**

```sql
-- ✅ migrations/000003_add_product_variants.up.sql
-- Create product variants table
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    -- Variant attributes
    sku VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    
    -- Pricing
    price DECIMAL(10,2) NOT NULL,
    compare_at_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    
    -- Inventory
    track_quantity BOOLEAN DEFAULT true,
    quantity INTEGER DEFAULT 0,
    allow_backorder BOOLEAN DEFAULT false,
    
    -- Physical attributes
    weight_grams INTEGER,
    length_cm DECIMAL(5,2),
    width_cm DECIMAL(5,2),
    height_cm DECIMAL(5,2),
    
    -- Metadata
    position INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT product_variants_price_positive CHECK (price >= 0),
    CONSTRAINT product_variants_quantity_non_negative CHECK (quantity >= 0)
);

-- Indexes
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_sku ON product_variants(sku);
CREATE INDEX idx_product_variants_active ON product_variants(is_active);

-- Variant attributes (EAV pattern for flexibility)
CREATE TABLE variant_attribute_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(variant_id, attribute_id)
);

CREATE INDEX idx_variant_attribute_values_variant ON variant_attribute_values(variant_id);
CREATE INDEX idx_variant_attribute_values_attribute ON variant_attribute_values(attribute_id);

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_product_variants_updated_at 
    BEFORE UPDATE ON product_variants 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- ✅ migrations/000003_add_product_variants.down.sql
-- Drop triggers first
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON product_variants;

-- Drop dependent tables
DROP TABLE IF EXISTS variant_attribute_values;

-- Drop main table
DROP TABLE IF EXISTS product_variants;

-- Drop function if no longer needed
-- (Be careful - other tables might use this function)
-- DROP FUNCTION IF EXISTS update_updated_at_column();
```

#### **🔧 Migration Commands**

```bash
# Apply migrations
make migrate-up              # Apply all pending migrations
make migrate-up-1            # Apply one migration
make migrate-up-to 5         # Apply up to specific version

# Rollback migrations
make migrate-down            # Rollback one migration
make migrate-down-1          # Same as migrate-down
make migrate-down-to 3       # Rollback to specific version

# Migration status
make migrate-version         # Show current version
make migrate-status          # Show detailed status

# Dangerous operations (be careful!)
make migrate-drop            # Drop all tables (⚠️ DATA LOSS)
make migrate-force 5         # Force set version without running migration
```

</details>

### 🌱 **Database Seeding**

<details>
<summary>📋 <strong>Seed Data Management</strong></summary>

#### **📁 Seed Data Structure**

```json
// seeds/users.json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "role_id": "admin-role-uuid",
    "username": "admin",
    "email": "admin@eshop.local", 
    "phone_number": "+1234567890",
    "first_name": "System",
    "last_name": "Administrator",
    "password": "SecureAdmin123!",
    "email_verified": true,
    "is_active": true
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001", 
    "role_id": "customer-role-uuid",
    "username": "customer1",
    "email": "customer1@example.com",
    "phone_number": "+1234567891",
    "first_name": "John",
    "last_name": "Doe", 
    "password": "Customer123!",
    "email_verified": true,
    "is_active": true
  }
]

// seeds/categories.json
[
  {
    "id": "cat-electronics-uuid",
    "parent_id": null,
    "name": "Electronics",
    "slug": "electronics",
    "description": "Electronic devices and accessories",
    "sort_order": 1,
    "is_active": true
  },
  {
    "id": "cat-smartphones-uuid",
    "parent_id": "cat-electronics-uuid", 
    "name": "Smartphones",
    "slug": "smartphones",
    "description": "Mobile phones and smartphones",
    "sort_order": 1,
    "is_active": true
  }
]

// seeds/products.json  
[
  {
    "id": "prod-iphone-uuid",
    "category_id": "cat-smartphones-uuid",
    "brand_id": "brand-apple-uuid",
    "name": "iPhone 15 Pro",
    "slug": "iphone-15-pro",
    "description": "Latest iPhone with advanced features",
    "short_description": "Premium smartphone with Pro camera system",
    "base_price": 999.99,
    "compare_at_price": 1099.99,
    "cost_price": 750.00,
    "track_quantity": true,
    "quantity": 50,
    "weight_grams": 187,
    "is_active": true,
    "is_featured": true,
    "meta_title": "iPhone 15 Pro - Premium Smartphone",
    "meta_description": "Buy the latest iPhone 15 Pro with advanced camera and performance"
  }
]
```

#### **🔧 Seed Implementation**

```go
// internal/db/seeds/seeder.go
package seeds

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "path/filepath"
    
    "github.com/yourorg/eshop/internal/db/repository"
    "github.com/yourorg/eshop/pkg/auth"
)

type Seeder struct {
    repo        repository.Repository
    passwordHasher auth.PasswordHasher
    seedDir     string
}

func NewSeeder(repo repository.Repository, passwordHasher auth.PasswordHasher, seedDir string) *Seeder {
    return &Seeder{
        repo:           repo,
        passwordHasher: passwordHasher,
        seedDir:        seedDir,
    }
}

func (s *Seeder) SeedAll(ctx context.Context) error {
    // Order matters due to foreign key constraints
    seedFunctions := []struct{
        name string
        fn   func(context.Context) error
    }{
        {"roles", s.seedRoles},
        {"users", s.seedUsers}, 
        {"categories", s.seedCategories},
        {"brands", s.seedBrands},
        {"attributes", s.seedAttributes},
        {"products", s.seedProducts},
        {"product_variants", s.seedProductVariants},
        {"shipping_zones", s.seedShippingZones},
        {"shipping_methods", s.seedShippingMethods},
        {"discounts", s.seedDiscounts},
    }
    
    for _, seed := range seedFunctions {
        fmt.Printf("Seeding %s...\n", seed.name)
        if err := seed.fn(ctx); err != nil {
            return fmt.Errorf("failed to seed %s: %w", seed.name, err)
        }
    }
    
    fmt.Println("✅ All seed data loaded successfully")
    return nil
}

func (s *Seeder) seedUsers(ctx context.Context) error {
    var users []SeedUser
    if err := s.loadJSONFile("users.json", &users); err != nil {
        return err
    }
    
    for _, user := range users {
        // Check if user already exists
        existing, err := s.repo.GetUserByEmail(ctx, user.Email)
        if err == nil {
            fmt.Printf("User %s already exists, skipping\n", existing.Email)
            continue
        }
        
        // Hash password
        hashedPassword, err := s.passwordHasher.HashPassword(user.Password)
        if err != nil {
            return fmt.Errorf("failed to hash password for user %s: %w", user.Email, err)
        }
        
        // Create user
        _, err = s.repo.CreateUser(ctx, repository.CreateUserParams{
            ID:             user.ID,
            RoleID:         user.RoleID,
            Username:       user.Username,
            Email:          user.Email,
            PhoneNumber:    user.PhoneNumber,
            FirstName:      user.FirstName,
            LastName:       user.LastName,
            HashedPassword: hashedPassword,
            EmailVerified:  user.EmailVerified,
            IsActive:       user.IsActive,
        })
        if err != nil {
            return fmt.Errorf("failed to create user %s: %w", user.Email, err)
        }
        
        fmt.Printf("✅ Created user: %s (%s)\n", user.Email, user.Username)
    }
    
    return nil
}

func (s *Seeder) loadJSONFile(filename string, v interface{}) error {
    filePath := filepath.Join(s.seedDir, filename)
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read seed file %s: %w", filename, err)
    }
    
    if err := json.Unmarshal(data, v); err != nil {
        return fmt.Errorf("failed to parse seed file %s: %w", filename, err)
    }
    
    return nil
}

// Seed data structures
type SeedUser struct {
    ID            uuid.UUID `json:"id"`
    RoleID        uuid.UUID `json:"role_id"`
    Username      string    `json:"username"`
    Email         string    `json:"email"`
    PhoneNumber   string    `json:"phone_number"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    Password      string    `json:"password"`
    EmailVerified bool      `json:"email_verified"`
    IsActive      bool      `json:"is_active"`
}
```

#### **⚡ Seed Commands**

```bash
# Seed all data
make seed

# Seed specific data type
make seed-users
make seed-products  
make seed-categories

# Clear and re-seed (development only)
make seed-fresh     # Drops all data and re-seeds

# Production seed (minimal data)
make seed-prod     # Only essential data for production
```

</details>

---

## 🧪 Comprehensive Testing Guide

### 🎯 **Testing Philosophy & Strategy**

<div align="center">

**Test Pyramid: Unit → Integration → E2E**

```
       🔺 E2E Tests (Few)
      ----/\----
     /  Complex  \
    /   Flows    \
   ----------------
  🔷 Integration (Some)  
 /   API + Database   \
/      Service       \
-----------------------
🟦   Unit Tests (Many)   
  Functions & Methods
```

</div>

<details>
<summary>🧪 <strong>Testing Strategy Deep Dive</strong></summary>

#### **📊 Test Types & Coverage Targets**

| **Test Type** | **Scope** | **Coverage Target** | **Execution Time** |
|--------------|-----------|-------------------|-------------------|
| 🔬 **Unit Tests** | Functions, methods | 80%+ | < 100ms |
| 🔗 **Integration Tests** | API endpoints, DB operations | 70%+ | < 1s |
| 🌐 **End-to-End Tests** | Complete user flows | 60%+ | < 10s |
| ⚡ **Performance Tests** | Load, stress testing | Key endpoints | < 30s |

#### **🎯 Testing Principles**

1. **🔒 Deterministic** - Tests should always produce the same result
2. **🔄 Isolated** - Tests shouldn't depend on each other
3. **⚡ Fast** - Quick feedback cycle for developers
4. **📝 Readable** - Clear test names and structure
5. **🎯 Focused** - One concept per test
6. **🔧 Maintainable** - Easy to update when code changes

</details>

### 🔬 **Unit Testing**

<details>
<summary>🧪 <strong>Unit Testing Patterns & Examples</strong></summary>

#### **📋 Test File Organization**

```
internal/
├── services/
│   ├── user_service.go
│   └── user_service_test.go      # Unit tests for user service
├── utils/
│   ├── validation.go  
│   └── validation_test.go        # Unit tests for validation
└── api/
    ├── users.go
    └── users_test.go             # Handler unit tests
```

#### **✅ Service Layer Unit Tests**

```go
// internal/services/user_service_test.go
package services

import (
    "context"
    "testing"
    "time"
    
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/google/uuid"
    
    "github.com/yourorg/eshop/internal/dto"
    "github.com/yourorg/eshop/internal/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name        string
        input       dto.CreateUserRequest
        setupMocks  func(*mocks.MockUserRepository, *mocks.MockPasswordHasher)
        expectedErr error
        expectUser  bool
    }{
        {
            name: "successful_user_creation",
            input: dto.CreateUserRequest{
                Username:    "testuser",
                Email:       "test@example.com",
                PhoneNumber: "+1234567890",
                FirstName:   "Test",
                LastName:    "User",
                Password:    "SecurePass123!",
            },
            setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher) {
                // Mock password hashing
                hasher.EXPECT().
                    HashPassword("SecurePass123!").
                    Return("hashed_password", nil).
                    Times(1)
                
                // Mock user creation
                repo.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(repository.User{
                        ID:        uuid.New(),
                        Username:  "testuser", 
                        Email:     "test@example.com",
                        CreatedAt: time.Now(),
                    }, nil).
                    Times(1)
                    
                // Mock email check (user doesn't exist)
                repo.EXPECT().
                    GetUserByEmail(gomock.Any(), "test@example.com").
                    Return(repository.User{}, sql.ErrNoRows).
                    Times(1)
            },
            expectedErr: nil,
            expectUser:  true,
        },
        {
            name: "user_already_exists",
            input: dto.CreateUserRequest{
                Username: "existinguser",
                Email:    "existing@example.com",
                Password: "password123",
            },
            setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher) {
                // Mock existing user check
                repo.EXPECT().
                    GetUserByEmail(gomock.Any(), "existing@example.com").
                    Return(repository.User{
                        ID:    uuid.New(),
                        Email: "existing@example.com",
                    }, nil).
                    Times(1)
            },
            expectedErr: ErrUserAlreadyExists,
            expectUser:  false,
        },
        {
            name: "password_hashing_failure",
            input: dto.CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com", 
                Password: "password123",
            },
            setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher) {
                // Mock user doesn't exist
                repo.EXPECT().
                    GetUserByEmail(gomock.Any(), "test@example.com").
                    Return(repository.User{}, sql.ErrNoRows).
                    Times(1)
                    
                // Mock password hashing failure
                hasher.EXPECT().
                    HashPassword("password123").
                    Return("", errors.New("hashing failed")).
                    Times(1)
            },
            expectedErr: errors.New("failed to hash password"),
            expectUser:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := mocks.NewMockUserRepository(ctrl)
            mockHasher := mocks.NewMockPasswordHasher(ctrl)
            
            service := NewUserService(mockRepo, mockHasher)
            
            // Configure mocks
            tt.setupMocks(mockRepo, mockHasher)
            
            // Execute
            user, err := service.CreateUser(context.Background(), tt.input)
            
            // Assert
            if tt.expectedErr != nil {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedErr.Error())
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                if tt.expectUser {
                    require.NotNil(t, user)
                    assert.Equal(t, tt.input.Username, user.Username)
                    assert.Equal(t, tt.input.Email, user.Email)
                    assert.NotZero(t, user.ID)
                    assert.NotZero(t, user.CreatedAt)
                }
            }
        })
    }
}

func TestUserService_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo, nil)
    
    userID := uuid.New()
    
    t.Run("user_found", func(t *testing.T) {
        expectedUser := repository.User{
            ID:        userID,
            Username:  "testuser",
            Email:     "test@example.com",
            CreatedAt: time.Now(),
        }
        
        mockRepo.EXPECT().
            GetUser(gomock.Any(), userID).
            Return(expectedUser, nil).
            Times(1)
            
        user, err := service.GetUser(context.Background(), userID)
        
        assert.NoError(t, err)
        require.NotNil(t, user)
        assert.Equal(t, expectedUser.ID, user.ID)
        assert.Equal(t, expectedUser.Username, user.Username)
    })
    
    t.Run("user_not_found", func(t *testing.T) {
        mockRepo.EXPECT().
            GetUser(gomock.Any(), userID).
            Return(repository.User{}, sql.ErrNoRows).
            Times(1)
            
        user, err := service.GetUser(context.Background(), userID)
        
        assert.Error(t, err)
        assert.Equal(t, ErrUserNotFound, err)
        assert.Nil(t, user)
    })
}
```

#### **🔧 Utility Function Unit Tests**

```go
// internal/utils/validation_test.go
package utils

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid_email", "user@example.com", false},
        {"valid_email_subdomain", "user@mail.example.com", false},
        {"valid_email_numbers", "user123@example.co.uk", false},
        {"invalid_no_at", "userexample.com", true},
        {"invalid_no_domain", "user@", true},
        {"invalid_no_tld", "user@example", true},
        {"invalid_spaces", "user @example.com", true},
        {"invalid_empty", "", true},
        {"invalid_only_at", "@", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err, "ValidateEmail(%q) should return error", tt.email)
            } else {
                assert.NoError(t, err, "ValidateEmail(%q) should not return error", tt.email)
            }
        })
    }
}

func TestGenerateSlug(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple_text", "Hello World", "hello-world"},
        {"with_numbers", "Product 123", "product-123"},
        {"with_special_chars", "Hello, World!", "hello-world"},
        {"multiple_spaces", "Hello   World", "hello-world"},
        {"unicode_chars", "Café & Restaurant", "cafe-restaurant"},
        {"already_slug", "hello-world", "hello-world"},
        {"empty_string", "", ""},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := GenerateSlug(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

// Benchmark tests for performance-critical functions
func BenchmarkValidateEmail(b *testing.B) {
    email := "test@example.com"
    for i := 0; i < b.N; i++ {
        ValidateEmail(email)
    }
}

func BenchmarkGenerateSlug(b *testing.B) {
    text := "This is a long product name with many words"
    for i := 0; i < b.N; i++ {
        GenerateSlug(text)
    }
}
```

</details>

### 🔗 **Integration Testing**

<details>
<summary>🌐 <strong>API Integration & Database Testing</strong></summary>

#### **🏗️ Test Database Setup**

```go
// tests/helpers/test_db.go
package helpers

import (
    "context"
    "database/sql"
    "fmt"
    "testing"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
    "github.com/stretchr/testify/require"
    
    "github.com/yourorg/eshop/config"
    "github.com/yourorg/eshop/internal/db/repository"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
    // Load test configuration
    cfg := config.Config{
        DBUrl: "postgresql://postgres:postgres@localhost:5433/eshop_test?sslmode=disable",
    }
    
    // Connect to database
    db, err := sql.Open("postgres", cfg.DBUrl)
    require.NoError(t, err)
    
    // Verify connection
    err = db.Ping()
    require.NoError(t, err)
    
    // Run migrations
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    require.NoError(t, err)
    
    m, err := migrate.NewWithDatabaseInstance(
        "file://../../migrations",
        "postgres",
        driver,
    )
    require.NoError(t, err)
    
    err = m.Up()
    require.NoError(t, err)
    
    // Return cleanup function
    cleanup := func() {
        // Clean up data after test
        tables := []string{
            "order_items", "orders", "cart_items", "carts",
            "product_variants", "products", "categories",
            "user_addresses", "users", "roles",
        }
        
        for _, table := range tables {
            _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
            if err != nil {
                t.Logf("Failed to truncate table %s: %v", table, err)
            }
        }
        
        db.Close()
    }
    
    return db, cleanup
}

func CreateTestUser(t *testing.T, db *sql.DB) repository.User {
    repo := repository.New(db)
    
    user, err := repo.CreateUser(context.Background(), repository.CreateUserParams{
        Username:       "testuser",
        Email:          "test@example.com",
        PhoneNumber:    "+1234567890",
        FirstName:      "Test",
        LastName:       "User",
        HashedPassword: "hashed_password",
        RoleID:         GetDefaultRoleID(t, db),
    })
    require.NoError(t, err)
    
    return user
}
```

#### **🌐 API Integration Tests**

```go
// tests/integration/api/users_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/yourorg/eshop/internal/api"
    "github.com/yourorg/eshop/internal/dto"
    "github.com/yourorg/eshop/tests/helpers"
)

func TestUserAPI_Integration(t *testing.T) {
    // Setup test database
    db, cleanup := helpers.SetupTestDB(t)
    defer cleanup()
    
    // Setup test server
    server := api.NewServer(api.Config{
        DB: db,
        // ... other config
    })
    
    router := chi.NewRouter()
    server.SetupRoutes(router)
    
    testServer := httptest.NewServer(router)
    defer testServer.Close()
    
    t.Run("POST /api/v1/users - Create User", func(t *testing.T) {
        // Prepare request
        userReq := dto.CreateUserRequest{
            Username:    "newuser",
            Email:       "newuser@example.com",
            PhoneNumber: "+1234567890",
            FirstName:   "New",
            LastName:    "User",
            Password:    "SecurePass123!",
        }
        
        reqBody, err := json.Marshal(userReq)
        require.NoError(t, err)
        
        // Make request
        resp, err := http.Post(
            testServer.URL+"/api/v1/users",
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        require.NoError(t, err)
        defer resp.Body.Close()
        
        // Verify response
        assert.Equal(t, http.StatusCreated, resp.StatusCode)
        assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
        
        var response struct {
            Success bool            `json:"success"`
            Data    dto.UserResponse `json:"data"`
            Message string          `json:"message"`
        }
        
        err = json.NewDecoder(resp.Body).Decode(&response)
        require.NoError(t, err)
        
        assert.True(t, response.Success)
        assert.Equal(t, userReq.Username, response.Data.Username)
        assert.Equal(t, userReq.Email, response.Data.Email)
        assert.NotZero(t, response.Data.ID)
        assert.NotZero(t, response.Data.CreatedAt)
    })
    
    t.Run("POST /api/v1/users - Duplicate Email", func(t *testing.T) {
        // Create user first
        existingUser := helpers.CreateTestUser(t, db)
        
        // Try to create user with same email
        userReq := dto.CreateUserRequest{
            Username:    "anotheruser",
            Email:       existingUser.Email, // Same email
            PhoneNumber: "+1234567891",
            FirstName:   "Another",
            LastName:    "User",
            Password:    "SecurePass123!",
        }
        
        reqBody, err := json.Marshal(userReq)
        require.NoError(t, err)
        
        resp, err := http.Post(
            testServer.URL+"/api/v1/users",
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusConflict, resp.StatusCode)
        
        var errorResp struct {
            Success bool   `json:"success"`
            Error   struct {
                Code    string `json:"code"`
                Message string `json:"message"`
            } `json:"error"`
        }
        
        err = json.NewDecoder(resp.Body).Decode(&errorResp)
        require.NoError(t, err)
        
        assert.False(t, errorResp.Success)
        assert.Equal(t, "USER_EXISTS", errorResp.Error.Code)
    })
    
    t.Run("GET /api/v1/users/{id} - Get User", func(t *testing.T) {
        // Create test user
        user := helpers.CreateTestUser(t, db)
        
        // Make request
        resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", testServer.URL, user.ID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var response struct {
            Success bool            `json:"success"`
            Data    dto.UserResponse `json:"data"`
        }
        
        err = json.NewDecoder(resp.Body).Decode(&response)
        require.NoError(t, err)
        
        assert.True(t, response.Success)
        assert.Equal(t, user.ID, response.Data.ID)
        assert.Equal(t, user.Username, response.Data.Username)
        assert.Equal(t, user.Email, response.Data.Email)
    })
    
    t.Run("GET /api/v1/users - List Users with Pagination", func(t *testing.T) {
        // Create multiple test users
        for i := 0; i < 5; i++ {
            helpers.CreateTestUser(t, db, helpers.UserOptions{
                Username: fmt.Sprintf("user%d", i),
                Email:    fmt.Sprintf("user%d@example.com", i),
            })
        }
        
        // Test pagination
        resp, err := http.Get(testServer.URL + "/api/v1/users?page=1&limit=3")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var response struct {
            Success    bool              `json:"success"`
            Data       []dto.UserResponse `json:"data"`
            Pagination struct {
                Page       int  `json:"page"`
                Limit      int  `json:"limit"`
                Total      int  `json:"total"`
                TotalPages int  `json:"totalPages"`
                HasNext    bool `json:"hasNext"`
            } `json:"pagination"`
        }
        
        err = json.NewDecoder(resp.Body).Decode(&response)
        require.NoError(t, err)
        
        assert.True(t, response.Success)
        assert.Len(t, response.Data, 3) // Should return 3 users
        assert.Equal(t, 1, response.Pagination.Page)
        assert.Equal(t, 3, response.Pagination.Limit)
        assert.True(t, response.Pagination.Total >= 5)
        assert.True(t, response.Pagination.HasNext)
    })
}

// Authentication integration tests
func TestAuthAPI_Integration(t *testing.T) {
    db, cleanup := helpers.SetupTestDB(t)
    defer cleanup()
    
    server := setupTestServer(t, db)
    
    t.Run("Login Flow", func(t *testing.T) {
        // Create test user
        user := helpers.CreateTestUser(t, db)
        
        // Login request
        loginReq := dto.LoginRequest{
            Username: user.Username,
            Password: "password123", // Known password
        }
        
        reqBody, _ := json.Marshal(loginReq)
        resp, err := http.Post(
            server.URL+"/api/v1/auth/login",
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var loginResp struct {
            Success bool   `json:"success"`
            Data    struct {
                AccessToken  string `json:"access_token"`
                RefreshToken string `json:"refresh_token"`
                User         dto.UserResponse `json:"user"`
            } `json:"data"`
        }
        
        err = json.NewDecoder(resp.Body).Decode(&loginResp)
        require.NoError(t, err)
        
        assert.True(t, loginResp.Success)
        assert.NotEmpty(t, loginResp.Data.AccessToken)
        assert.NotEmpty(t, loginResp.Data.RefreshToken)
        assert.Equal(t, user.ID, loginResp.Data.User.ID)
        
        // Test protected endpoint with token
        req, _ := http.NewRequest("GET", server.URL+"/api/v1/users/profile", nil)
        req.Header.Set("Authorization", "Bearer "+loginResp.Data.AccessToken)
        
        client := &http.Client{}
        resp, err = client.Do(req)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })
}
```

</details>

### ⚡ **Performance Testing**

<details>
<summary>🚀 <strong>Load Testing & Benchmarks</strong></summary>

#### **📊 Benchmark Tests**

```go
// tests/performance/api_benchmark_test.go
package performance

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/yourorg/eshop/internal/dto"
    "github.com/yourorg/eshop/tests/helpers"
)

func BenchmarkUserAPI_CreateUser(b *testing.B) {
    db, cleanup := helpers.SetupTestDB(b)
    defer cleanup()
    
    server := helpers.SetupTestServer(b, db)
    testServer := httptest.NewServer(server)
    defer testServer.Close()
    
    userReq := dto.CreateUserRequest{
        Username:    "benchuser",
        Email:       "bench@example.com",
        PhoneNumber: "+1234567890",
        FirstName:   "Bench",
        LastName:    "User", 
        Password:    "SecurePass123!",
    }
    
    reqBody, _ := json.Marshal(userReq)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Use unique email for each request
            userReq.Email = fmt.Sprintf("bench%d@example.com", rand.Int())
            reqBody, _ := json.Marshal(userReq)
            
            resp, err := http.Post(
                testServer.URL+"/api/v1/users",
                "application/json",
                bytes.NewBuffer(reqBody),
            )
            if err != nil {
                b.Errorf("Request failed: %v", err)
            }
            resp.Body.Close()
            
            if resp.StatusCode != http.StatusCreated {
                b.Errorf("Expected 201, got %d", resp.StatusCode)
            }
        }
    })
}

func BenchmarkUserAPI_GetUser(b *testing.B) {
    db, cleanup := helpers.SetupTestDB(b)
    defer cleanup()
    
    // Create test users
    userIDs := make([]uuid.UUID, 100)
    for i := 0; i < 100; i++ {
        user := helpers.CreateTestUser(b, db, helpers.UserOptions{
            Username: fmt.Sprintf("benchuser%d", i),
            Email:    fmt.Sprintf("bench%d@example.com", i),
        })
        userIDs[i] = user.ID
    }
    
    server := helpers.SetupTestServer(b, db)
    testServer := httptest.NewServer(server)
    defer testServer.Close()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Get random user
            userID := userIDs[rand.Intn(len(userIDs))]
            
            resp, err := http.Get(
                fmt.Sprintf("%s/api/v1/users/%s", testServer.URL, userID),
            )
            if err != nil {
                b.Errorf("Request failed: %v", err)
            }
            resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                b.Errorf("Expected 200, got %d", resp.StatusCode)
            }
        }
    })
}

// Database operation benchmarks
func BenchmarkRepository_CreateUser(b *testing.B) {
    db, cleanup := helpers.SetupTestDB(b)
    defer cleanup()
    
    repo := repository.New(db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := repo.CreateUser(context.Background(), repository.CreateUserParams{
            Username:       fmt.Sprintf("benchuser%d", i),
            Email:          fmt.Sprintf("bench%d@example.com", i),
            PhoneNumber:    "+1234567890",
            FirstName:      "Bench",
            LastName:       "User",
            HashedPassword: "hashed_password",
            RoleID:         uuid.New(),
        })
        if err != nil {
            b.Errorf("CreateUser failed: %v", err)
        }
    }
}

func BenchmarkRepository_GetUser(b *testing.B) {
    db, cleanup := helpers.SetupTestDB(b)
    defer cleanup()
    
    repo := repository.New(db)
    
    // Create test users
    userIDs := make([]uuid.UUID, 1000)
    for i := 0; i < 1000; i++ {
        user := helpers.CreateTestUser(b, db)
        userIDs[i] = user.ID
    }
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            userID := userIDs[rand.Intn(len(userIDs))]
            _, err := repo.GetUser(context.Background(), userID)
            if err != nil {
                b.Errorf("GetUser failed: %v", err)
            }
        }
    })
}
```

#### **🔥 Load Testing with Go**

```go
// tests/load/load_test.go
package load

import (
    "context"
    "fmt"
    "sync"
    "testing"
    "time"
    
    "github.com/yourorg/eshop/tests/helpers"
)

func TestUserAPI_LoadTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }
    
    db, cleanup := helpers.SetupTestDB(t)
    defer cleanup()
    
    server := helpers.SetupTestServer(t, db)
    testServer := httptest.NewServer(server)
    defer testServer.Close()
    
    // Load test configuration
    const (
        concurrentUsers = 50
        requestsPerUser = 100
        testDuration    = 30 * time.Second
    )
    
    // Metrics collection
    var (
        totalRequests  int64
        successCount   int64
        errorCount     int64
        responseTimes  []time.Duration
        mutex          sync.Mutex
    )
    
    // Start load test
    ctx, cancel := context.WithTimeout(context.Background(), testDuration)
    defer cancel()
    
    var wg sync.WaitGroup
    
    for i := 0; i < concurrentUsers; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            client := &http.Client{Timeout: 5 * time.Second}
            
            for j := 0; j < requestsPerUser; j++ {
                select {
                case <-ctx.Done():
                    return
                default:
                }
                
                start := time.Now()
                
                // Make request
                resp, err := client.Get(
                    fmt.Sprintf("%s/api/v1/products?page=1&limit=20", testServer.URL),
                )
                
                duration := time.Since(start)
                
                // Record metrics
                mutex.Lock()
                totalRequests++
                responseTimes = append(responseTimes, duration)
                if err != nil || resp.StatusCode != http.StatusOK {
                    errorCount++
                } else {
                    successCount++
                }
                mutex.Unlock()
                
                if resp != nil {
                    resp.Body.Close()
                }
                
                // Small delay between requests
                time.Sleep(10 * time.Millisecond)
            }
        }(i)
    }
    
    wg.Wait()
    
    // Calculate statistics
    var totalResponseTime time.Duration
    for _, rt := range responseTimes {
        totalResponseTime += rt
    }
    
    avgResponseTime := totalResponseTime / time.Duration(len(responseTimes))
    successRate := float64(successCount) / float64(totalRequests) * 100
    requestsPerSecond := float64(totalRequests) / testDuration.Seconds()
    
    // Report results
    t.Logf("Load Test Results:")
    t.Logf("  Duration: %v", testDuration)
    t.Logf("  Concurrent Users: %d", concurrentUsers)
    t.Logf("  Total Requests: %d", totalRequests)
    t.Logf("  Successful Requests: %d", successCount)
    t.Logf("  Failed Requests: %d", errorCount)
    t.Logf("  Success Rate: %.2f%%", successRate)
    t.Logf("  Requests/Second: %.2f", requestsPerSecond)
    t.Logf("  Average Response Time: %v", avgResponseTime)
    
    // Assertions
    if successRate < 95.0 {
        t.Errorf("Success rate too low: %.2f%% (expected >= 95%%)", successRate)
    }
    
    if avgResponseTime > 1*time.Second {
        t.Errorf("Average response time too high: %v (expected <= 1s)", avgResponseTime)
    }
    
    if requestsPerSecond < 100 {
        t.Errorf("Throughput too low: %.2f req/s (expected >= 100)", requestsPerSecond)
    }
}
```

</details>

---

## 🔒 Security Guidelines

### 🛡️ **Security-First Development**

<details>
<summary>🔐 <strong>Authentication & Authorization</strong></summary>

#### **🔑 JWT/PASETO Token Management**

```go
// pkg/auth/token_maker.go
type TokenMaker interface {
    // CreateToken creates a new token for a specific username and duration
    CreateToken(username string, role string, duration time.Duration) (string, error)
    
    // VerifyToken checks if the token is valid or not
    VerifyToken(token string) (*Payload, error)
}

// PASETO implementation (recommended for production)
type PasetoMaker struct {
    paseto       *paseto.V2
    symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
    if len(symmetricKey) != chacha20poly1305.KeySize {
        return nil, fmt.Errorf("invalid key size: must be exactly %d characters", 
            chacha20poly1305.KeySize)
    }
    
    maker := &PasetoMaker{
        paseto:       paseto.NewV2(),
        symmetricKey: []byte(symmetricKey),
    }
    
    return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, error) {
    payload, err := NewPayload(username, role, duration)
    if err != nil {
        return "", err
    }
    
    token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
    return token, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
    payload := &Payload{}
    
    err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    err = payload.Valid()
    if err != nil {
        return nil, err
    }
    
    return payload, nil
}

// Token payload structure
type Payload struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    IssuedAt  time.Time `json:"issued_at"`
    ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, role string, duration time.Duration) (*Payload, error) {
    tokenID, err := uuid.NewRandom()
    if err != nil {
        return nil, err
    }
    
    payload := &Payload{
        ID:        tokenID,
        Username:  username,
        Role:      role,
        IssuedAt:  time.Now(),
        ExpiredAt: time.Now().Add(duration),
    }
    
    return payload, nil
}

func (payload *Payload) Valid() error {
    if time.Now().After(payload.ExpiredAt) {
        return ErrExpiredToken
    }
    return nil
}
```

#### **🔐 Password Security**

```go
// pkg/auth/password.go
type PasswordHasher interface {
    HashPassword(password string) (string, error)
    CheckPassword(password, hashedPassword string) error
}

type BcryptHasher struct {
    cost int
}

func NewBcryptHasher() PasswordHasher {
    return &BcryptHasher{
        cost: bcrypt.DefaultCost, // Currently 10
    }
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
    // Validate password strength before hashing
    if err := ValidatePasswordStrength(password); err != nil {
        return "", err
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    
    return string(hashedPassword), nil
}

func (h *BcryptHasher) CheckPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Password strength validation
func ValidatePasswordStrength(password string) error {
    if len(password) < 8 {
        return ErrPasswordTooShort
    }
    
    if len(password) > 128 {
        return ErrPasswordTooLong
    }
    
    var (
        hasUpper   bool
        hasLower   bool
        hasNumber  bool
        hasSpecial bool
    )
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if !hasUpper {
        return ErrPasswordMissingUpper
    }
    if !hasLower {
        return ErrPasswordMissingLower
    }
    if !hasNumber {
        return ErrPasswordMissingNumber
    }
    if !hasSpecial {
        return ErrPasswordMissingSpecial
    }
    
    // Check against common passwords
    if isCommonPassword(password) {
        return ErrPasswordTooCommon
    }
    
    return nil
}
```

#### **🔒 Rate Limiting**

```go
// internal/api/middleware/rate_limit.go
func RateLimitMiddleware(store *redis.Client, requests int, window time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client identifier (IP + User ID if authenticated)
            clientID := getClientIdentifier(r)
            
            // Check rate limit
            allowed, remaining, resetTime, err := checkRateLimit(store, clientID, requests, window)
            if err != nil {
                log.Error().Err(err).Str("client_id", clientID).Msg("Rate limit check failed")
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            
            // Set rate limit headers
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
            w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
            
            if !allowed {
                w.Header().Set("Retry-After", fmt.Sprintf("%.0f", time.Until(resetTime).Seconds()))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func getClientIdentifier(r *http.Request) string {
    // Try to get authenticated user ID
    if userCtx := r.Context().Value("user"); userCtx != nil {
        if payload, ok := userCtx.(*auth.Payload); ok {
            return fmt.Sprintf("user:%s", payload.Username)
        }
    }
    
    // Fall back to IP address
    ip := getRealIP(r)
    return fmt.Sprintf("ip:%s", ip)
}

func getRealIP(r *http.Request) string {
    // Check X-Forwarded-For header (load balancer)
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        ips := strings.Split(xff, ",")
        return strings.TrimSpace(ips[0])
    }
    
    // Check X-Real-IP header (reverse proxy)
    if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
        return xrip
    }
    
    // Fall back to RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    return ip
}

func checkRateLimit(store *redis.Client, clientID string, limit int, window time.Duration) (bool, int, time.Time, error) {
    ctx := context.Background()
    key := fmt.Sprintf("rate_limit:%s", clientID)
    
    // Use sliding window with Redis
    now := time.Now()
    windowStart := now.Add(-window)
    
    pipe := store.TxPipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
    
    // Count current requests
    countCmd := pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{
        Score:  float64(now.UnixNano()),
        Member: now.UnixNano(),
    })
    
    // Set expiration
    pipe.Expire(ctx, key, window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, 0, time.Time{}, err
    }
    
    count := int(countCmd.Val())
    remaining := limit - count
    if remaining < 0 {
        remaining = 0
    }
    
    resetTime := now.Add(window)
    allowed := count <= limit
    
    return allowed, remaining, resetTime, nil
}
```

</details>

### 🔍 **Input Validation & Sanitization**

<details>
<summary>✅ <strong>Comprehensive Input Validation</strong></summary>

#### **📋 Request Validation**

```go
// internal/api/validators.go
import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
    
    // Register custom validators
    validate.RegisterValidation("password", validatePassword)
    validate.RegisterValidation("username", validateUsername)
    validate.RegisterValidation("phone", validatePhone)
    validate.RegisterValidation("slug", validateSlug)
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    return ValidatePasswordStrength(password) == nil
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    
    // Username rules: 3-30 chars, alphanumeric + underscore/dash
    if len(username) < 3 || len(username) > 30 {
        return false
    }
    
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
    return matched
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    
    // E.164 format validation
    matched, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phone)
    return matched
}

func validateSlug(fl validator.FieldLevel) bool {
    slug := fl.Field().String()
    
    // URL-safe slug validation
    matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, slug)
    return matched
}

// DTO with validation tags
type CreateUserRequest struct {
    Username    string `json:"username" validate:"required,username"`
    Email       string `json:"email" validate:"required,email"`
    PhoneNumber string `json:"phone_number" validate:"required,phone"`
    FirstName   string `json:"first_name" validate:"required,min=1,max=50"`
    LastName    string `json:"last_name" validate:"required,min=1,max=50"`
    Password    string `json:"password" validate:"required,password"`
    DateOfBirth string `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

type CreateProductRequest struct {
    Name        string    `json:"name" validate:"required,min=3,max=255"`
    Slug        string    `json:"slug" validate:"required,slug,min=3,max=255"`
    Description string    `json:"description" validate:"required,min=10,max=5000"`
    Price       float64   `json:"price" validate:"required,gt=0"`
    CategoryID  uuid.UUID `json:"category_id" validate:"required"`
    BrandID     uuid.UUID `json:"brand_id,omitempty"`
    Tags        []string  `json:"tags,omitempty" validate:"omitempty,dive,min=2,max=50"`
}

// Validation helper functions
func ValidateRequest(req interface{}) error {
    if err := validate.Struct(req); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        return formatValidationError(validationErrors)
    }
    return nil
}

func formatValidationError(validationErrors validator.ValidationErrors) error {
    errors := make(map[string]string)
    
    for _, err := range validationErrors {
        field := strings.ToLower(err.Field())
        
        switch err.Tag() {
        case "required":
            errors[field] = fmt.Sprintf("%s is required", field)
        case "email":
            errors[field] = "Invalid email format"
        case "min":
            errors[field] = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
        case "max":
            errors[field] = fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
        case "gt":
            errors[field] = fmt.Sprintf("%s must be greater than %s", field, err.Param())
        case "password":
            errors[field] = "Password must be at least 8 characters with uppercase, lowercase, number, and special character"
        case "username":
            errors[field] = "Username must be 3-30 characters, alphanumeric with underscore/dash"
        case "phone":
            errors[field] = "Phone number must be in E.164 format (e.g., +1234567890)"
        case "slug":
            errors[field] = "Slug must be lowercase alphanumeric with hyphens"
        default:
            errors[field] = fmt.Sprintf("%s is invalid", field)
        }
    }
    
    return ValidationError{Errors: errors}
}

type ValidationError struct {
    Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
    return "Validation failed"
}
```

#### **🧹 Data Sanitization**

```go
// internal/utils/sanitize.go
import (
    "github.com/microcosm-cc/bluemonday"
    "html"
    "strings"
)

var htmlSanitizer *bluemonday.Policy

func init() {
    // Create strict HTML sanitizer
    htmlSanitizer = bluemonday.StrictPolicy()
}

// Sanitize HTML content
func SanitizeHTML(input string) string {
    return htmlSanitizer.Sanitize(input)
}

// Sanitize plain text (escape HTML entities)
func SanitizePlainText(input string) string {
    return html.EscapeString(strings.TrimSpace(input))
}

// Sanitize user input for database storage
func SanitizeUserInput(input string) string {
    // Trim whitespace
    input = strings.TrimSpace(input)
    
    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")
    
    // Normalize line endings
    input = strings.ReplaceAll(input, "\r\n", "\n")
    input = strings.ReplaceAll(input, "\r", "\n")
    
    // Remove control characters except tab and newline
    var cleaned strings.Builder
    for _, r := range input {
        if r == '\t' || r == '\n' || (r >= 32 && r != 127) {
            cleaned.WriteRune(r)
        }
    }
    
    return cleaned.String()
}

// Sanitize filename for safe file storage
func SanitizeFilename(filename string) string {
    // Replace unsafe characters
    unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
    for _, char := range unsafe {
        filename = strings.ReplaceAll(filename, char, "_")
    }
    
    // Trim dots and spaces from start/end
    filename = strings.Trim(filename, ". ")
    
    // Limit length
    if len(filename) > 255 {
        ext := filepath.Ext(filename)
        name := strings.TrimSuffix(filename, ext)
        maxNameLen := 255 - len(ext)
        if maxNameLen > 0 {
            filename = name[:maxNameLen] + ext
        } else {
            filename = "file" + ext
        }
    }
    
    return filename
}

// Sanitize SQL input (additional layer beyond parameterized queries)
func SanitizeSQLInput(input string) string {
    // Remove SQL comment patterns
    input = regexp.MustCompile(`--.*`).ReplaceAllString(input, "")
    input = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(input, "")
    
    // Remove common SQL injection patterns
    dangerous := []string{"'", "\"", ";", "/*", "*/", "xp_", "sp_"}
    for _, pattern := range dangerous {
        input = strings.ReplaceAll(input, pattern, "")
    }
    
    return strings.TrimSpace(input)
}
```

</details>

### 🔐 **Data Protection & Privacy**

<details>
<summary>🛡️ <strong>Encryption & Data Privacy</strong></summary>

#### **🔒 Field-Level Encryption**

```go
// pkg/encryption/field_encryption.go
type FieldEncryptor interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}

type AESEncryptor struct {
    key []byte
}

func NewAESEncryptor(key string) (*AESEncryptor, error) {
    keyBytes := []byte(key)
    if len(keyBytes) != 32 { // AES-256
        return nil, fmt.Errorf("key must be 32 bytes for AES-256")
    }
    
    return &AESEncryptor{key: keyBytes}, nil
}

func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
    if plaintext == "" {
        return "", nil
    }
    
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }
    
    // Generate random IV
    iv := make([]byte, aes.BlockSize)
    if _, err := rand.Read(iv); err != nil {
        return "", err
    }
    
    // Encrypt
    stream := cipher.NewCFBEncrypter(block, iv)
    ciphertext := make([]byte, len(plaintext))
    stream.XORKeyStream(ciphertext, []byte(plaintext))
    
    // Combine IV + ciphertext and encode as base64
    combined := append(iv, ciphertext...)
    return base64.StdEncoding.EncodeToString(combined), nil
}

func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
    if ciphertext == "" {
        return "", nil
    }
    
    // Decode from base64
    combined, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    if len(combined) < aes.BlockSize {
        return "", fmt.Errorf("ciphertext too short")
    }
    
    // Split IV and ciphertext
    iv := combined[:aes.BlockSize]
    encryptedData := combined[aes.BlockSize:]
    
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }
    
    // Decrypt
    stream := cipher.NewCFBDecrypter(block, iv)
    plaintext := make([]byte, len(encryptedData))
    stream.XORKeyStream(plaintext, encryptedData)
    
    return string(plaintext), nil
}

// Usage in models
type User struct {
    ID              uuid.UUID `json:"id" db:"id"`
    Username        string    `json:"username" db:"username"`
    Email           string    `json:"email" db:"email"`
    PhoneNumber     string    `json:"phone_number" db:"phone_number"`         // Encrypted
    DateOfBirth     string    `json:"date_of_birth,omitempty" db:"date_of_birth"` // Encrypted
    SSN             string    `json:"-" db:"ssn"`                            // Encrypted, never in JSON
    HashedPassword  string    `json:"-" db:"hashed_password"`                // Hashed, never in JSON
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Encryption helper methods
func (u *User) EncryptSensitiveFields(encryptor FieldEncryptor) error {
    var err error
    
    if u.PhoneNumber != "" {
        u.PhoneNumber, err = encryptor.Encrypt(u.PhoneNumber)
        if err != nil {
            return fmt.Errorf("failed to encrypt phone number: %w", err)
        }
    }
    
    if u.DateOfBirth != "" {
        u.DateOfBirth, err = encryptor.Encrypt(u.DateOfBirth)
        if err != nil {
            return fmt.Errorf("failed to encrypt date of birth: %w", err)
        }
    }
    
    if u.SSN != "" {
        u.SSN, err = encryptor.Encrypt(u.SSN)
        if err != nil {
            return fmt.Errorf("failed to encrypt SSN: %w", err)
        }
    }
    
    return nil
}

func (u *User) DecryptSensitiveFields(encryptor FieldEncryptor) error {
    var err error
    
    if u.PhoneNumber != "" {
        u.PhoneNumber, err = encryptor.Decrypt(u.PhoneNumber)
        if err != nil {
            return fmt.Errorf("failed to decrypt phone number: %w", err)
        }
    }
    
    if u.DateOfBirth != "" {
        u.DateOfBirth, err = encryptor.Decrypt(u.DateOfBirth)
        if err != nil {
            return fmt.Errorf("failed to decrypt date of birth: %w", err)
        }
    }
    
    if u.SSN != "" {
        u.SSN, err = encryptor.Decrypt(u.SSN)
        if err != nil {
            return fmt.Errorf("failed to decrypt SSN: %w", err)
        }
    }
    
    return nil
}
```

#### **🔒 Secure Headers Middleware**

```go
// internal/api/middleware/security.go
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Prevent MIME type sniffing
            w.Header().Set("X-Content-Type-Options", "nosniff")
            
            // Prevent clickjacking
            w.Header().Set("X-Frame-Options", "DENY")
            
            // XSS protection (legacy browsers)
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            
            // Strict Transport Security (HTTPS only)
            if r.TLS != nil {
                w.Header().Set("Strict-Transport-Security", 
                    "max-age=31536000; includeSubDomains; preload")
            }
            
            // Content Security Policy
            csp := strings.Join([]string{
                "default-src 'self'",
                "script-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com",
                "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
                "font-src 'self' https://fonts.gstatic.com",
                "img-src 'self' data: https:",
                "connect-src 'self'",
                "frame-ancestors 'none'",
                "base-uri 'self'",
                "form-action 'self'",
            }, "; ")
            w.Header().Set("Content-Security-Policy", csp)
            
            // Referrer Policy
            w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
            
            // Feature Policy
            w.Header().Set("Permissions-Policy", 
                "geolocation=(), microphone=(), camera=()")
            
            next.ServeHTTP(w, r)
        })
    }
}
```

</details>

---

## ⚡ Performance Optimization

### 🚀 **Application Performance**

<details>
<summary>📊 <strong>Performance Optimization Strategies</strong></summary>

#### **🗄️ Database Query Optimization**

```go
// Efficient pagination with cursor-based approach
func (r *Repository) GetProductsWithCursor(
    ctx context.Context, 
    cursor string, 
    limit int,
) ([]Product, string, error) {
    query := `
        SELECT id, name, price, created_at
        FROM products 
        WHERE ($1::timestamptz IS NULL OR created_at < $1)
        ORDER BY created_at DESC, id DESC
        LIMIT $2 + 1  -- Get one extra to determine if there's a next page
    `
    
    var cursorTime *time.Time
    if cursor != "" {
        t, err := time.Parse(time.RFC3339, cursor)
        if err != nil {
            return nil, "", fmt.Errorf("invalid cursor: %w", err)
        }
        cursorTime = &t
    }
    
    rows, err := r.db.QueryContext(ctx, query, cursorTime, limit)
    if err != nil {
        return nil, "", err
    }
    defer rows.Close()
    
    var products []Product
    for rows.Next() {
        var p Product
        err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)
        if err != nil {
            return nil, "", err
        }
        products = append(products, p)
    }
    
    // Determine next cursor
    var nextCursor string
    if len(products) > limit {
        // Remove the extra item
        products = products[:limit]
        nextCursor = products[limit-1].CreatedAt.Format(time.RFC3339)
    }
    
    return products, nextCursor, nil
}

// Batch operations for efficiency
func (r *Repository) CreateProductsBatch(ctx context.Context, products []CreateProductParams) error {
    if len(products) == 0 {
        return nil
    }
    
    // Use COPY for large batch inserts
    if len(products) > 100 {
        return r.copyProducts(ctx, products)
    }
    
    // Use multi-row INSERT for smaller batches
    valueStrings := make([]string, len(products))
    valueArgs := make([]interface{}, 0, len(products)*6)
    
    for i, p := range products {
        valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", 
            i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
        valueArgs = append(valueArgs, p.Name, p.Slug, p.Description, 
            p.Price, p.CategoryID, p.BrandID)
    }
    
    query := fmt.Sprintf(`
        INSERT INTO products (name, slug, description, price, category_id, brand_id)
        VALUES %s
    `, strings.Join(valueStrings, ","))
    
    _, err := r.db.ExecContext(ctx, query, valueArgs...)
    return err
}

// Optimized search with full-text search
func (r *Repository) SearchProducts(ctx context.Context, query string, limit, offset int) ([]Product, error) {
    sqlQuery := `
        SELECT 
            p.id, p.name, p.slug, p.price, p.image_url,
            ts_rank(p.search_vector, plainto_tsquery('english', $1)) as rank
        FROM products p
        WHERE p.search_vector @@ plainto_tsquery('english', $1)
            AND p.is_active = true
        ORDER BY rank DESC, p.created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.QueryContext(ctx, sqlQuery, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var products []Product
    for rows.Next() {
        var p Product
        var rank float64
        err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Price, &p.ImageURL, &rank)
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    
    return products, nil
}
```

#### **⚡ Caching Strategies**

```go
// Multi-layer caching implementation
type CacheService struct {
    redis      *redis.Client
    localCache *cache.Cache[string, interface{}]
    mutex      sync.RWMutex
}

func NewCacheService(redisClient *redis.Client) *CacheService {
    // Local cache for frequently accessed data
    localCache := cache.New[string, interface{}](
        cache.WithTTL[string, interface{}](5 * time.Minute),
        cache.WithEvictionCallback[string, interface{}](func(key string, value interface{}) {
            log.Debug().Str("key", key).Msg("Local cache entry evicted")
        }),
    )
    
    return &CacheService{
        redis:      redisClient,
        localCache: localCache,
    }
}

// Get with fallback strategy: Local -> Redis -> Database
func (c *CacheService) Get(ctx context.Context, key string, fallback func() (interface{}, error)) (interface{}, error) {
    // Try local cache first
    if value, found := c.localCache.Get(key); found {
        return value, nil
    }
    
    // Try Redis cache
    redisValue, err := c.redis.Get(ctx, key).Result()
    if err == nil {
        var value interface{}
        if err := json.Unmarshal([]byte(redisValue), &value); err == nil {
            // Store in local cache
            c.localCache.Set(key, value)
            return value, nil
        }
    }
    
    // Fallback to database
    if fallback != nil {
        value, err := fallback()
        if err != nil {
            return nil, err
        }
        
        // Store in both caches
        c.Set(ctx, key, value, time.Hour)
        return value, nil
    }
    
    return nil, ErrCacheNotFound
}

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // Store in local cache
    c.localCache.Set(key, value)
    
    // Store in Redis
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.redis.Set(ctx, key, data, ttl).Err()
}

// Cache warm-up for frequently accessed data
func (s *ProductService) WarmupCache(ctx context.Context) error {
    // Get popular products
    popularProducts, err := s.repo.GetPopularProducts(ctx, 100)
    if err != nil {
        return err
    }
    
    // Cache them
    for _, product := range popularProducts {
        cacheKey := fmt.Sprintf("product:%s", product.ID)
        s.cache.Set(ctx, cacheKey, product, time.Hour)
    }
    
    // Get featured categories
    categories, err := s.repo.GetFeaturedCategories(ctx)
    if err != nil {
        return err
    }
    
    s.cache.Set(ctx, "categories:featured", categories, 30*time.Minute)
    
    return nil
}
```

#### **🔄 Connection Pooling**

```go
// Database connection pool optimization
func setupDatabasePool(databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database URL: %w", err)
    }
    
    // Pool configuration
    config.MaxConns = 30                    // Maximum connections
    config.MinConns = 5                     // Minimum connections to maintain
    config.MaxConnLifetime = time.Hour      // Max time a connection can be reused
    config.MaxConnIdleTime = 30 * time.Minute // Max time a connection can be idle
    config.HealthCheckPeriod = 5 * time.Minute // Health check frequency
    
    // Connection configuration
    config.ConnConfig.ConnectTimeout = 10 * time.Second
    config.ConnConfig.PreferSimpleProtocol = true // Faster for simple queries
    
    // Statement timeout
    config.ConnConfig.RuntimeParams["statement_timeout"] = "30s"
    config.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = "60s"
    
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }
    
    return pool, nil
}

// Redis connection pool
func setupRedisPool(redisURL string) *redis.Client {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to parse Redis URL")
    }
    
    // Pool settings
    opt.PoolSize = 20                    // Number of socket connections
    opt.MinIdleConns = 5                 // Minimum idle connections
    opt.PoolTimeout = 30 * time.Second   // Pool timeout
    opt.IdleTimeout = 5 * time.Minute    // Close connections after this time
    opt.IdleCheckFrequency = time.Minute // Frequency to check idle connections
    
    // Connection timeouts
    opt.DialTimeout = 5 * time.Second
    opt.ReadTimeout = 3 * time.Second
    opt.WriteTimeout = 3 * time.Second
    
    client := redis.NewClient(opt)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to Redis")
    }
    
    return client
}
```

</details>

### 🔧 **Code-Level Optimizations**

<details>
<summary>⚡ <strong>Go Performance Best Practices</strong></summary>

#### **🧠 Memory Management**

```go
// Object pooling for frequently allocated objects
var (
    // Pool for request/response objects
    requestPool = sync.Pool{
        New: func() interface{} {
            return &dto.ProductSearchRequest{}
        },
    }
    
    // Pool for string builders
    stringBuilderPool = sync.Pool{
        New: func() interface{} {
            return &strings.Builder{}
        },
    }
    
    // Pool for byte slices
    byteSlicePool = sync.Pool{
        New: func() interface{} {
            return make([]byte, 0, 1024) // Initial capacity
        },
    }
)

// Using object pools in handlers
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
    // Get request object from pool
    req := requestPool.Get().(*dto.ProductSearchRequest)
    defer func() {
        // Reset and return to pool
        *req = dto.ProductSearchRequest{}
        requestPool.Put(req)
    }()
    
    // Parse request
    if err := json.NewDecoder(r.Body).Decode(req); err != nil {
        h.sendError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    // Process request...
}

// Efficient string building
func buildProductDescription(product Product, variants []ProductVariant) string {
    sb := stringBuilderPool.Get().(*strings.Builder)
    defer func() {
        sb.Reset()
        stringBuilderPool.Put(sb)
    }()
    
    sb.WriteString(product.Description)
    sb.WriteString("\n\nAvailable variants:\n")
    
    for _, variant := range variants {
        sb.WriteString("- ")
        sb.WriteString(variant.Name)
        sb.WriteString(" ($")
        sb.WriteString(fmt.Sprintf("%.2f", variant.Price))
        sb.WriteString(")\n")
    }
    
    return sb.String()
}

// Pre-allocating slices with known capacity
func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error) {
    // Estimate capacity based on average category size
    estimatedSize := 50
    products := make([]Product, 0, estimatedSize)
    
    // Query products...
    rows, err := s.repo.GetProductsByCategoryQuery(ctx, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var product Product
        if err := rows.Scan(/* fields */); err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    
    return products, nil
}
```

#### **⚡ Concurrent Processing**

```go
// Worker pool for concurrent processing
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultQueue chan Result
    wg         sync.WaitGroup
}

type Job struct {
    ID   string
    Data interface{}
}

type Result struct {
    JobID string
    Data  interface{}
    Error error
}

func NewWorkerPool(workers int, jobBuffer int) *WorkerPool {
    return &WorkerPool{
        workers:     workers,
        jobQueue:    make(chan Job, jobBuffer),
        resultQueue: make(chan Result, jobBuffer),
    }
}

func (wp *WorkerPool) Start(processor func(Job) Result) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(processor)
    }
}

func (wp *WorkerPool) worker(processor func(Job) Result) {
    defer wp.wg.Done()
    
    for job := range wp.jobQueue {
        result := processor(job)
        wp.resultQueue <- result
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.resultQueue
}

func (wp *WorkerPool) Close() {
    close(wp.jobQueue)
    wp.wg.Wait()
    close(wp.resultQueue)
}

// Example usage: Parallel image processing
func (s *ProductService) ProcessProductImages(ctx context.Context, productID uuid.UUID, images []Image) error {
    pool := NewWorkerPool(5, len(images)) // 5 workers
    pool.Start(func(job Job) Result {
        image := job.Data.(Image)
        
        // Process image (resize, optimize, upload)
        processedURL, err := s.imageProcessor.Process(ctx, image)
        return Result{
            JobID: job.ID,
            Data:  processedURL,
            Error: err,
        }
    })
    defer pool.Close()
    
    // Submit jobs
    for i, image := range images {
        pool.Submit(Job{
            ID:   fmt.Sprintf("image_%d", i),
            Data: image,
        })
    }
    
    // Collect results
    var processedImages []string
    for i := 0; i < len(images); i++ {
        result := <-pool.Results()
        if result.Error != nil {
            return fmt.Errorf("failed to process image %s: %w", result.JobID, result.Error)
        }
        processedImages = append(processedImages, result.Data.(string))
    }
    
    // Save processed image URLs
    return s.repo.UpdateProductImages(ctx, productID, processedImages)
}

// Context-aware concurrent operations
func (s *OrderService) ProcessOrdersConcurrently(ctx context.Context, orderIDs []uuid.UUID) error {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // Limit concurrent operations
    semaphore := make(chan struct{}, 5) // Max 5 concurrent operations
    errChan := make(chan error, len(orderIDs))
    
    var wg sync.WaitGroup
    
    for _, orderID := range orderIDs {
        wg.Add(1)
        go func(id uuid.UUID) {
            defer wg.Done()
            
            // Acquire semaphore
            select {
            case semaphore <- struct{}{}:
            case <-ctx.Done():
                errChan <- ctx.Err()
                return
            }
            defer func() { <-semaphore }() // Release semaphore
            
            if err := s.processOrder(ctx, id); err != nil {
                errChan <- err
            }
        }(orderID)
    }
    
    // Wait for completion
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    // Collect errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

#### **📊 Profiling & Monitoring**

```go
// Performance monitoring middleware
func PerformanceMonitoringMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            next.ServeHTTP(wrapped, r)
            
            duration := time.Since(start)
            
            // Log performance metrics
            log.Info().
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Int("status", wrapped.statusCode).
                Dur("duration", duration).
                Int64("size", wrapped.size).
                Msg("HTTP request completed")
                
            // Record metrics (Prometheus, StatsD, etc.)
            recordHTTPMetrics(r.Method, r.URL.Path, wrapped.statusCode, duration)
            
            // Alert on slow requests
            if duration > 2*time.Second {
                log.Warn().
                    Str("method", r.Method).
                    Str("path", r.URL.Path).
                    Dur("duration", duration).
                    Msg("Slow request detected")
            }
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    size, err := rw.ResponseWriter.Write(b)
    rw.size += int64(size)
    return size, err
}

// Memory usage monitoring
func MonitorMemoryUsage() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        log.Info().
            Uint64("alloc_mb", m.Alloc/1024/1024).
            Uint64("total_alloc_mb", m.TotalAlloc/1024/1024).
            Uint64("sys_mb", m.Sys/1024/1024).
            Uint32("num_gc", m.NumGC).
            Uint64("heap_objects", m.HeapObjects).
            Msg("Memory usage")
            
        // Force garbage collection if memory usage is high
        if m.Alloc > 500*1024*1024 { // 500MB
            log.Warn().Msg("High memory usage, forcing GC")
            runtime.GC()
        }
    }
}

// Database query performance monitoring
func (r *Repository) QueryWithMetrics(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    
    duration := time.Since(start)
    
    // Log slow queries
    if duration > 100*time.Millisecond {
        log.Warn().
            Str("query", query).
            Interface("args", args).
            Dur("duration", duration).
            Msg("Slow query detected")
    }
    
    // Record query metrics
    recordDatabaseMetrics("query", duration, err)
    
    return rows, err
}
```

</details>

---

## Contributing Guidelines

### Git Workflow

1. **Branching Strategy**
   ```bash
   # Feature branch
   git checkout -b feature/user-profile-update
   
   # Bug fix branch
   git checkout -b bugfix/order-calculation-error
   
   # Hotfix branch
   git checkout -b hotfix/security-patch
   ```

2. **Commit Messages**
   ```bash
   # Format: type(scope): description
   git commit -m "feat(auth): add JWT token refresh endpoint"
   git commit -m "fix(orders): correct tax calculation logic"
   git commit -m "docs(api): update authentication documentation"
   ```

3. **Pull Request Process**
   - Create feature branch
   - Implement changes with tests
   - Update documentation
   - Submit PR with clear description
   - Address review feedback
   - Merge after approval

### Code Review Checklist

- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Security best practices followed
- [ ] Performance considerations addressed
- [ ] Error handling is appropriate
- [ ] Logging is adequate

### Pre-commit Hooks

```bash
# Install pre-commit hooks
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: gofmt
        language: system
        args: [-w, -s]
        files: \.go$
      
      - id: go-lint
        name: golangci-lint
        entry: golangci-lint
        language: system
        args: [run]
        files: \.go$
```

This development guide provides a comprehensive foundation for working with the e-commerce platform. Regular updates and team feedback will help maintain its relevance and usefulness.