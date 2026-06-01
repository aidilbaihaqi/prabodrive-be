# domain

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder is the **core layer** of Clean Architecture. Contains business entities and interfaces.

## Files

| File | Description |
|------|-------------|
| `user.go` | User entity with business methods |
| `repository.go` | Repository interfaces |
| `errors.go` | Domain-specific errors |

## user.go

User Entity with:
- Struct fields (ID, Email, Password, Name, Role, etc.)
- Business methods (`Validate()`, `IsAdmin()`, `CanAccess()`, etc.)
- No framework dependencies

```go
type User struct {
    ID       uuid.UUID
    Email    string
    Password string
    Name     string
    Role     string
    IsActive bool
}

func (u *User) Validate() error {
    if u.Email == "" {
        return ErrEmailRequired
    }
    return nil
}
```

## repository.go

Repository interfaces to be implemented in the repository layer:

```go
type UserRepository interface {
    Create(user *User) error
    FindByID(id uuid.UUID) (*User, error)
    FindByEmail(email string) (*User, error)
    FindAll(filters UserFilters) ([]*User, int64, error)
    Update(user *User) error
    Delete(id uuid.UUID) error
}
```

## errors.go

Domain-specific errors for business rules:

```go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrEmailRequired    = errors.New("email is required")
    ErrInvalidEmail     = errors.New("invalid email format")
    ErrEmailExists      = errors.New("email already exists")
    ErrUnauthorized     = errors.New("unauthorized")
)
```

## Rules

### ✅ DO
- Pure Go code only
- Business logic in entity methods
- Define interfaces here
- Define custom errors here

### ❌ DON'T
- Import GORM, Gin, or other frameworks
- Import from other layers (repository, usecase, delivery)
- Database-specific code
- HTTP-specific code

## Adding a New Entity

1. Create file `product.go`:
```go
type Product struct {
    ID    uuid.UUID
    Name  string
    Price float64
}

func (p *Product) Validate() error { ... }
```

2. Add interface in `repository.go`:
```go
type ProductRepository interface {
    Create(product *Product) error
    FindByID(id uuid.UUID) (*Product, error)
}
```

3. Add errors in `errors.go`:
```go
var ErrProductNotFound = errors.New("product not found")
```
