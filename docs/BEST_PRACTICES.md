# Best Practices

## 1. Domain Layer

### ✅ DO: Keep entities pure
```go
type Product struct {
    ID    uuid.UUID
    Name  string
    Price float64
    Stock int
}

func (p *Product) CanBeSold() bool {
    return p.Stock > 0 && p.Price > 0
}
```

### ❌ DON'T: Import frameworks
```go
// WRONG
import "gorm.io/gorm"

type Product struct {
    gorm.Model  // Don't embed GORM in domain
}
```

## 2. Usecase Layer

### ✅ DO: One usecase per operation
```go
type CreateUserUsecase struct {}
type UpdateUserUsecase struct {}
type DeleteUserUsecase struct {}
```

### ❌ DON'T: Combine multiple operations
```go
// WRONG
type UserUsecase struct {
    // Create, Update, Delete all in one - hard to test
}
```

## 3. Repository Layer

### ✅ DO: Map domain ↔ model
```go
func (r *repo) toDomain(m *models.User) *domain.User {
    return &domain.User{
        ID:    m.ID,
        Email: m.Email,
    }
}

func (r *repo) toModel(u *domain.User) *models.User {
    return &models.User{
        ID:    u.ID,
        Email: u.Email,
    }
}
```

### ✅ DO: Handle errors properly
```go
func (r *repo) FindByID(id uuid.UUID) (*domain.User, error) {
    var model models.User
    err := r.db.Where("id = ?", id).First(&model).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domain.ErrUserNotFound  // Use domain error
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return r.toDomain(&model), nil
}
```

## 4. Error Handling

### ✅ DO: Use domain errors
```go
// domain/errors.go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrEmailExists  = errors.New("email already exists")
)

// In handler
if errors.Is(err, domain.ErrUserNotFound) {
    response.NotFound(c, err.Error())
    return
}
```

## 5. Testing

### ✅ DO: Use interfaces for mocking
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    uc := user.NewCreateUserUsecase(mockRepo)
    // Test...
}
```

## 6. API Design

### ✅ DO: Use consistent response format
```go
// Success
{
    "success": true,
    "data": {...}
}

// Error
{
    "success": false,
    "error": {
        "code": "NOT_FOUND",
        "message": "User not found"
    }
}
```

### ✅ DO: Validate input at handler level
```go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required,min=2"`
}
```

## 7. Security

### ✅ DO: Hash passwords with bcrypt
```go
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,
)
```

### ✅ DO: Use environment variables for secrets
```go
jwtSecret := os.Getenv("JWT_ACCESS_SECRET")
```

### ❌ DON'T: Log sensitive data
```go
// WRONG
log.Printf("User login: %s, password: %s", email, password)
```

## 8. Database

### ✅ DO: Use migrations
```sql
-- migrations/000001_init.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(254) UNIQUE NOT NULL
);
```

### ✅ DO: Use connection pooling
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## 9. Logging

### ✅ DO: Use structured logging
```go
log.Printf("[%s] %s %s %d %v",
    requestID,
    method,
    path,
    statusCode,
    latency,
)
```

### ✅ DO: Include request ID
```go
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

## 10. Configuration

### ✅ DO: Use typed config struct
```go
type Config struct {
    AppEnv   string
    AppPort  string
    Database DatabaseConfig
    JWT      JWTConfig
}

func Load() *Config {
    return &Config{
        AppEnv:  getEnv("APP_ENV", "development"),
        AppPort: getEnv("APP_PORT", "8080"),
    }
}
```
