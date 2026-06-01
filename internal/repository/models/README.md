# models

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains GORM models for database.

## Files

| File | Description |
|------|-------------|
| `user.go` | GORM model for `users` table |

## user.go

```go
type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Password  string         `gorm:"not null"`
    Name      string         `gorm:"not null"`
    Role      string         `gorm:"default:user"`
    IsActive  bool           `gorm:"default:true"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Difference: Domain Entity vs GORM Model

| Aspect | Domain Entity | GORM Model |
|--------|---------------|------------|
| Location | `domain/user.go` | `repository/models/user.go` |
| Tags | None | GORM tags |
| Methods | Business logic | Table mapping |
| Import | Pure Go | GORM |

## Why Separate?

1. **Clean Architecture** - Domain must not depend on framework
2. **Flexibility** - Can change ORM without changing domain
3. **Testing** - Domain is easy to test without database

## Adding a New Model

1. Create file `product.go`:
```go
type Product struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string    `gorm:"not null"`
    Price     float64   `gorm:"not null"`
    Stock     int       `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (Product) TableName() string {
    return "products"
}
```

2. Create migration in `migrations/`
