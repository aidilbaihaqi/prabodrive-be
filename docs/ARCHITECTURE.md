# Architecture Documentation

This document explains the Clean Architecture implementation in this Go backend template.

## Overview

Clean Architecture is a software design philosophy that separates concerns into layers, with dependencies pointing inward toward the core business logic.

## Layer Breakdown

### 1. Domain Layer (`internal/domain/`)

The **innermost layer** containing:
- **Entities**: Core business objects (e.g., `User`, `Product`)
- **Repository Interfaces**: Contracts for data access
- **Domain Errors**: Business-specific error definitions

**Rules:**
- ✅ Pure Go code, no external dependencies
- ✅ Business logic resides in entity methods
- ❌ No framework imports (Gin, GORM, etc.)
- ❌ No knowledge of outer layers

```go
// domain/user.go
type User struct {
    ID       uuid.UUID
    Email    string
    Password string
    Name     string
}

func (u *User) Validate() error {
    if u.Email == "" {
        return ErrEmailRequired
    }
    return nil
}
```

### 2. Usecase Layer (`internal/usecase/`)

Contains **application business logic**:
- One usecase per operation (Single Responsibility)
- Orchestrates data flow between layers
- Input/Output DTOs for each usecase

**Rules:**
- ✅ Depends on Domain layer (interfaces)
- ✅ Contains transaction management
- ❌ No HTTP/framework concepts
- ❌ No direct database access

```go
// usecase/user/create.go
type CreateUserUsecase struct {
    userRepo domain.UserRepository
}

func (uc *CreateUserUsecase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Validate
    // 2. Hash password
    // 3. Create entity
    // 4. Save via repository
    // 5. Return output
}
```

### 3. Repository Layer (`internal/repository/`)

**Data access implementation**:
- Implements domain interfaces
- Contains ORM models (separate from domain)
- Handles database operations

**Rules:**
- ✅ Implements domain interfaces
- ✅ Maps between domain entities and DB models
- ❌ Domain layer must not import this

```go
// repository/user_repository.go
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(user *domain.User) error {
    model := r.toModel(user)  // Domain -> DB Model
    return r.db.Create(model).Error
}
```

### 4. Delivery Layer (`internal/delivery/`)

**Transport/Presentation layer**:
- HTTP handlers (Gin)
- Request/Response DTOs
- Route registration

**Rules:**
- ✅ Depends on Usecase layer
- ✅ Maps HTTP requests to usecase inputs
- ✅ Handles HTTP responses
- ❌ No business logic here

```go
// delivery/http/handler/user_handler.go
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    result, err := h.createUserUC.Execute(input)
    // Handle response
}
```

## Dependency Flow

```
Delivery → Usecase → Domain ← Repository
```

- **Delivery** depends on **Usecase**
- **Usecase** depends on **Domain** (interfaces)
- **Repository** implements **Domain** interfaces
- **Domain** has no dependencies

## Folder Structure

```
internal/
├── domain/                 # Core business layer
│   ├── user.go            # User entity
│   ├── repository.go      # Repository interfaces
│   └── errors.go          # Domain errors
│
├── usecase/               # Application layer
│   └── user/
│       ├── create.go      # One file per operation
│       ├── get.go
│       ├── list.go
│       └── update.go
│
├── repository/            # Data layer
│   ├── models/           # GORM models
│   │   └── user.go
│   └── user_repository.go
│
├── delivery/              # Presentation layer
│   └── http/
│       ├── handler/      # HTTP handlers
│       ├── routes/       # Route registration
│       ├── request/      # Request DTOs
│       └── response/     # Response helpers
│
├── middleware/            # HTTP middleware
│   ├── auth.go
│   └── middleware.go
│
├── infrastructure/        # External services
│   └── database/
│       └── postgres.go
│
└── shared/               # Shared utilities
    ├── utils/
    ├── token/
    └── constants/
```

## Benefits

1. **Testability**: Easy to mock dependencies
2. **Maintainability**: Changes isolated to specific layers
3. **Flexibility**: Easy to swap implementations
4. **Clarity**: Clear separation of concerns

## Anti-Patterns to Avoid

❌ **Direct DB access in handlers**
```go
// WRONG
func (h *Handler) Get(c *gin.Context) {
    h.db.Where("id = ?", id).First(&user)
}
```

❌ **HTTP concepts in usecase**
```go
// WRONG
func (uc *Usecase) Execute(c *gin.Context) {
    // Using Gin context in usecase
}
```

❌ **Business logic in repository**
```go
// WRONG
func (r *repo) Create(user *domain.User) error {
    if user.Age < 18 {  // Business rule shouldn't be here
        return errors.New("too young")
    }
}
```
