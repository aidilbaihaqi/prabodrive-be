# repository

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains data access layer implementation.

## Structure

```
repository/
├── models/
│   └── user.go              # GORM model
└── user_repository.go        # Repository implementation
```

## Files

### models/user.go

GORM model separate from domain entity:

```go
type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Password  string         `gorm:"not null"`
    Name      string         `gorm:"not null"`
    Role      string         `gorm:"default:user"`
    IsActive  bool           `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### user_repository.go

Implementation of `domain.UserRepository` interface:

```go
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    model := r.toModel(user)
    return r.db.Create(model).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
    var model models.User
    err := r.db.Where("id = ?", id).First(&model).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.ErrUserNotFound
    }
    return r.toDomain(&model), err
}
```

## Mapping Functions

Each repository should have helpers for mapping:

```go
// Domain -> Model (for write operations)
func (r *userRepository) toModel(user *domain.User) *models.User

// Model -> Domain (for read operations)
func (r *userRepository) toDomain(model *models.User) *domain.User
```

## Error Handling

Convert database errors to domain errors:

```go
func (r *userRepository) handleError(err error) error {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return domain.ErrUserNotFound
    }
    return err
}
```

## Adding a New Repository

1. Create model in `models/product.go`
2. Create `product_repository.go`
3. Implement interface from domain

## Best Practices

- ✅ Use model separate from domain entity
- ✅ Mapping with helper functions
- ✅ Convert DB errors to domain errors
- ✅ Use Preload() for relations
- ❌ Don't expose GORM outside repository
- ❌ No business logic in repository
