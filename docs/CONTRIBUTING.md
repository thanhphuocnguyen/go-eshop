# Contributing Guide

## Welcome Contributors!

Thank you for your interest in contributing to the e-commerce platform! This guide will help you get started and ensure your contributions align with our project standards.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Process](#development-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Review Process](#review-process)

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, sex characteristics, gender identity and expression, level of experience, education, socio-economic status, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Unacceptable behavior includes:**
- The use of sexualized language or imagery
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Go 1.24+** installed
- **PostgreSQL 14+** for database
- **Redis 6+** for caching
- **Docker & Docker Compose** for local development
- **Git** for version control
- **Make** for build automation

### Development Setup

1. **Fork the Repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/go-eshop.git
   cd go-eshop/server
   ```

2. **Set Up Remote**
   ```bash
   git remote add upstream https://github.com/thanhphuocnguyen/go-eshop.git
   git remote -v
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Environment Setup**
   ```bash
   cp app.env.example app.env
   # Edit app.env with your local configuration
   ```

5. **Start Services**
   ```bash
   docker-compose up -d postgres redis
   ```

6. **Database Migration**
   ```bash
   make migrate-up
   make seed  # Optional: Add sample data
   ```

7. **Run Tests**
   ```bash
   go test ./...
   ```

8. **Start Development Server**
   ```bash
   make serve-server
   ```

### Development Tools

Install recommended development tools:

```bash
# Air for live reloading
go install github.com/air-verse/air@latest

# golangci-lint for code linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
```

## Development Process

### Git Workflow

We follow the [GitHub Flow](https://guides.github.com/introduction/flow/) for development:

1. **Create Feature Branch**
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Implement your feature
   - Write/update tests
   - Update documentation

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(scope): add new feature description"
   ```

4. **Keep Branch Updated**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create Pull Request on GitHub
   ```

### Branch Naming Conventions

Use descriptive branch names with prefixes:

- `feature/user-authentication` - New features
- `bugfix/order-calculation-error` - Bug fixes
- `hotfix/security-vulnerability` - Critical fixes
- `docs/api-documentation` - Documentation updates
- `refactor/database-layer` - Code refactoring
- `test/integration-tests` - Test improvements

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

**Examples:**
```
feat(auth): add JWT token refresh endpoint

Implement automatic token refresh functionality to improve
user experience and reduce authentication errors.

Closes #123
```

```
fix(orders): correct tax calculation for international orders

- Fix incorrect tax rate application
- Add validation for country tax codes
- Update test cases

Fixes #456
```

## Coding Standards

### Go Style Guide

We follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

#### Naming Conventions

```go
// Constants: Use PascalCase
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// Variables: Use camelCase
var (
    httpClient *http.Client
    dbConn     *sql.DB
)

// Functions: Exported functions use PascalCase
func CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Implementation
}

// Private functions use camelCase
func validateEmail(email string) bool {
    // Implementation
}

// Structs: Use PascalCase
type UserService struct {
    repo   UserRepository
    logger *slog.Logger
}

// Interfaces: Use PascalCase, often ending with -er
type UserRepository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUser(ctx context.Context, id uuid.UUID) (*User, error)
}
```

#### Function Structure

```go
// Good function structure
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // 1. Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }

    // 2. Business logic
    hashedPassword, err := s.hashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("password hashing failed: %w", err)
    }

    // 3. Database operation
    user := &User{
        Email:          req.Email,
        HashedPassword: hashedPassword,
        CreatedAt:      time.Now(),
    }

    if err := s.repo.CreateUser(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 4. Return result
    return user, nil
}
```

#### Error Handling

```go
// Always handle errors explicitly
user, err := repo.GetUser(ctx, userID)
if err != nil {
    if errors.Is(err, ErrUserNotFound) {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return nil, fmt.Errorf("database error: %w", err)
}

// Define custom error types for domain errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidEmail     = errors.New("invalid email format")
    ErrPasswordTooWeak  = errors.New("password does not meet requirements")
)

// Use error wrapping for context
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) error {
    if err := s.repo.UpdateUser(ctx, id, req); err != nil {
        return fmt.Errorf("failed to update user %s: %w", id, err)
    }
    return nil
}
```

#### Context Usage

```go
// Always accept context as the first parameter
func (s *OrderService) ProcessOrder(ctx context.Context, orderID uuid.UUID) error {
    // Check for cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Use context for database operations
    order, err := s.repo.GetOrder(ctx, orderID)
    if err != nil {
        return err
    }

    // Pass context to downstream services
    if err := s.paymentService.ProcessPayment(ctx, order.PaymentInfo); err != nil {
        return err
    }

    return nil
}
```

### Database Guidelines

#### SQL Style

```sql
-- Table names: snake_case, plural
CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Always use explicit column lists
INSERT INTO users (username, email, hashed_password) 
VALUES ($1, $2, $3) 
RETURNING id, created_at;

-- Use meaningful aliases
SELECT 
    u.username,
    p.name as product_name,
    oi.quantity
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id;
```

#### SQLC Queries

```sql
-- name: CreateUser :one
INSERT INTO users (
    role_id, username, email, phone_number,
    first_name, last_name, hashed_password
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 AND locked = false 
LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;
```

### API Guidelines

#### Request/Response DTOs

```go
// Request DTOs with validation tags
type CreateUserRequest struct {
    Username    string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Email       string `json:"email" binding:"required,email"`
    PhoneNumber string `json:"phone_number" binding:"required,e164"`
    FirstName   string `json:"first_name" binding:"required,min=2,max=50"`
    LastName    string `json:"last_name" binding:"required,min=2,max=50"`
    Password    string `json:"password" binding:"required,min=8"`
}

func (r CreateUserRequest) Validate() error {
    // Additional custom validation logic
    if !isValidPassword(r.Password) {
        return ErrPasswordTooWeak
    }
    return nil
}

// Response DTOs - only include necessary fields
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Handler Structure

```go
// @Summary Create a new user
// @Description Create a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration data"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users [post]
func (s *Server) CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest

    // 1. Parse and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: ErrorDetail{
                Code:    "VALIDATION_ERROR",
                Message: "Invalid request format",
                Details: err.Error(),
            },
        })
        return
    }

    // 2. Additional validation
    if err := req.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: ErrorDetail{
                Code:    "VALIDATION_ERROR",
                Message: err.Error(),
            },
        })
        return
    }

    // 3. Business logic
    user, err := s.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        if errors.Is(err, ErrUserAlreadyExists) {
            c.JSON(http.StatusConflict, ErrorResponse{
                Error: ErrorDetail{
                    Code:    "USER_EXISTS",
                    Message: "User with this email already exists",
                },
            })
            return
        }

        s.logger.Error("Failed to create user", "error", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: ErrorDetail{
                Code:    "INTERNAL_ERROR",
                Message: "Failed to create user",
            },
        })
        return
    }

    // 4. Success response
    c.JSON(http.StatusCreated, toUserResponse(user))
}
```

## Testing Requirements

### Test Structure

All contributions must include appropriate tests:

#### Unit Tests

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name        string
        request     CreateUserRequest
        mockSetup   func(*mocks.MockUserRepository)
        want        *User
        wantErr     bool
        expectedErr error
    }{
        {
            name: "successful user creation",
            request: CreateUserRequest{
                Username:  "testuser",
                Email:     "test@example.com",
                FirstName: "Test",
                LastName:  "User",
                Password:  "SecurePass123!",
            },
            mockSetup: func(repo *mocks.MockUserRepository) {
                repo.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            want: &User{
                Username:  "testuser",
                Email:     "test@example.com",
                FirstName: "Test",
                LastName:  "User",
            },
            wantErr: false,
        },
        {
            name: "duplicate email error",
            request: CreateUserRequest{
                Email: "existing@example.com",
            },
            mockSetup: func(repo *mocks.MockUserRepository) {
                repo.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(ErrUserAlreadyExists)
            },
            wantErr:     true,
            expectedErr: ErrUserAlreadyExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockRepo := mocks.NewMockUserRepository(ctrl)
            if tt.mockSetup != nil {
                tt.mockSetup(mockRepo)
            }

            service := NewUserService(mockRepo, slog.Default())

            got, err := service.CreateUser(context.Background(), tt.request)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErr != nil {
                    assert.True(t, errors.Is(err, tt.expectedErr))
                }
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, got)
            assert.Equal(t, tt.want.Username, got.Username)
            assert.Equal(t, tt.want.Email, got.Email)
        })
    }
}
```

#### Integration Tests

```go
func TestCreateUserIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    repo := repository.New(db)
    service := NewUserService(repo, slog.Default())

    // Test data
    req := CreateUserRequest{
        Username:    "integration_test_user",
        Email:       "integration@test.com",
        PhoneNumber: "+1234567890",
        FirstName:   "Integration",
        LastName:    "Test",
        Password:    "SecurePass123!",
    }

    // Execute test
    user, err := service.CreateUser(context.Background(), req)

    // Assertions
    require.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, req.Username, user.Username)
    assert.Equal(t, req.Email, user.Email)
    assert.NotEmpty(t, user.CreatedAt)

    // Verify in database
    dbUser, err := repo.GetUser(context.Background(), user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.Username, dbUser.Username)
}
```

#### API Tests

```go
func TestCreateUserAPI(t *testing.T) {
    server := setupTestServer(t)
    defer server.Cleanup()

    tests := []struct {
        name         string
        requestBody  CreateUserRequest
        expectedCode int
        checkResult  func(t *testing.T, body []byte)
    }{
        {
            name: "valid user creation",
            requestBody: CreateUserRequest{
                Username:    "apitest",
                Email:       "api@test.com",
                PhoneNumber: "+1234567890",
                FirstName:   "API",
                LastName:    "Test",
                Password:    "SecurePass123!",
            },
            expectedCode: http.StatusCreated,
            checkResult: func(t *testing.T, body []byte) {
                var response UserResponse
                err := json.Unmarshal(body, &response)
                require.NoError(t, err)
                assert.Equal(t, "apitest", response.Username)
                assert.NotEmpty(t, response.ID)
            },
        },
        {
            name: "invalid email format",
            requestBody: CreateUserRequest{
                Username: "test",
                Email:    "invalid-email",
                Password: "SecurePass123!",
            },
            expectedCode: http.StatusBadRequest,
            checkResult: func(t *testing.T, body []byte) {
                var response ErrorResponse
                err := json.Unmarshal(body, &response)
                require.NoError(t, err)
                assert.Equal(t, "VALIDATION_ERROR", response.Error.Code)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            server.Router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedCode, w.Code)
            
            if tt.checkResult != nil {
                tt.checkResult(t, w.Body.Bytes())
            }
        })
    }
}
```

### Test Coverage

Maintain minimum test coverage:
- **Unit Tests**: 80% coverage minimum
- **Integration Tests**: Critical paths covered
- **API Tests**: All endpoints covered

Run coverage analysis:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Code Documentation

```go
// Package user provides user management functionality for the e-commerce platform.
// It handles user registration, authentication, profile management, and related operations.
package user

// UserService handles user-related business logic.
// It coordinates between the repository layer and the API handlers.
type UserService struct {
    repo   UserRepository
    logger *slog.Logger
    hasher PasswordHasher
}

// CreateUser creates a new user account after validating the input data.
// It returns the created user with generated ID and timestamps.
//
// The function performs the following steps:
// 1. Validates the input request
// 2. Checks if the email is already registered
// 3. Hashes the password using bcrypt
// 4. Creates the user record in the database
// 5. Returns the created user
//
// Returns ErrUserAlreadyExists if a user with the same email exists.
// Returns validation errors for invalid input data.
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Implementation
}

// validateEmail checks if the provided email address is valid.
// It uses a regular expression to validate the email format.
func validateEmail(email string) bool {
    // Implementation
}
```

### API Documentation

Use Swagger annotations for all endpoints:

```go
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile information
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse "User profile data"
// @Failure 401 {object} ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/profile [get]
func (s *Server) GetUserProfileHandler(c *gin.Context) {
    // Implementation
}
```

### Database Documentation

Document database changes in migration files:

```sql
-- Migration: 000003_add_user_preferences.up.sql
-- Description: Add user preferences table to store user-specific settings
-- Author: Your Name
-- Date: 2024-01-15

-- Add preferences column to users table
ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';

-- Add index for JSON queries
CREATE INDEX idx_users_preferences_gin ON users USING GIN(preferences);

-- Add comments for documentation
COMMENT ON COLUMN users.preferences IS 'User preferences stored as JSON (theme, notifications, etc.)';
```

## Submitting Changes

### Pull Request Process

1. **Ensure Quality**
   ```bash
   # Run tests
   go test ./...
   
   # Run linter
   golangci-lint run
   
   # Check formatting
   gofmt -d .
   
   # Generate documentation
   make swagger
   ```

2. **Update Documentation**
   - Update API documentation (Swagger)
   - Update README if needed
   - Add migration documentation
   - Update CHANGELOG.md

3. **Create Pull Request**
   - Use clear, descriptive title
   - Fill out the PR template
   - Link related issues
   - Add appropriate labels

### Pull Request Template

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## How Has This Been Tested?
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
- [ ] API tests

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Screenshots (if applicable)

## Additional Notes
Any additional information, concerns, or context.
```

### Pre-submission Checklist

- [ ] **Code Quality**
  - [ ] Code follows style guidelines
  - [ ] No linting errors
  - [ ] Proper error handling
  - [ ] Meaningful variable/function names

- [ ] **Testing**
  - [ ] Unit tests added/updated
  - [ ] Integration tests added/updated
  - [ ] All tests pass
  - [ ] Test coverage maintained

- [ ] **Documentation**
  - [ ] Code documentation updated
  - [ ] API documentation updated
  - [ ] Database changes documented
  - [ ] README updated if needed

- [ ] **Security**
  - [ ] No sensitive data in code
  - [ ] Input validation implemented
  - [ ] Authentication/authorization checked
  - [ ] SQL injection prevention

## Review Process

### Review Criteria

Reviewers will check for:

1. **Code Quality**
   - Adherence to style guide
   - Proper error handling
   - Performance considerations
   - Security implications

2. **Testing**
   - Adequate test coverage
   - Test quality and maintainability
   - Edge cases covered

3. **Documentation**
   - Code comments where needed
   - API documentation updated
   - Database changes documented

4. **Architecture**
   - Follows project architecture
   - Proper separation of concerns
   - No circular dependencies

### Addressing Feedback

- Respond to all review comments
- Make requested changes promptly
- Ask for clarification if needed
- Update tests and documentation as required

### Approval Process

1. **Code Review**: At least one maintainer approval
2. **Automated Checks**: All CI checks must pass
3. **Final Review**: Additional review for breaking changes
4. **Merge**: Maintainer will merge the PR

## Questions and Support

### Getting Help

- **GitHub Discussions**: For general questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Discord/Slack**: For real-time chat (if available)

### Reporting Issues

When reporting bugs, please include:

1. **Environment Information**
   - Go version
   - Operating system
   - Database version
   - Relevant environment variables

2. **Steps to Reproduce**
   - Clear, step-by-step instructions
   - Sample code if applicable
   - Expected vs. actual behavior

3. **Additional Context**
   - Error messages
   - Logs
   - Screenshots if relevant

### Feature Requests

For new features, please:

1. Check existing issues first
2. Provide clear use case description
3. Explain expected behavior
4. Consider implementation approach
5. Discuss with maintainers before starting work

## Recognition

We value all contributions, big and small. Contributors will be:

- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Recognized in project documentation

Thank you for contributing to our e-commerce platform! Your efforts help make this project better for everyone.