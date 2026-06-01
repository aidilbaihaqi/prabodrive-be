# middleware

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains HTTP middleware for Gin.

## Files

| File | Description |
|------|-------------|
| `auth.go` | JWT authentication middleware |
| `middleware.go` | Common middleware (CORS, RequestID, Logger, etc.) |

## auth.go

### Auth()
Validates JWT token and extracts user info:

```go
users.Use(middleware.Auth(jwtSecret))
```

Will:
- Check Authorization header
- Validate Bearer token
- Parse JWT claims
- Set `user_id` and `user_role` in context

### RequireRole()
Checks if user has specific role:

```go
admin := v1.Group("/admin")
admin.Use(middleware.Auth(jwtSecret))
admin.Use(middleware.RequireRole("admin"))
```

### OptionalAuth()
Tries to authenticate but doesn't require it:

```go
public.Use(middleware.OptionalAuth(jwtSecret))
```

## middleware.go

### CORS()
Handles Cross-Origin Resource Sharing:

```go
router.Use(middleware.CORS())
```

### RequestID()
Generates unique ID for each request:

```go
router.Use(middleware.RequestID())
// Access via c.Get("request_id")
```

### Logger()
Logs request details:

```go
router.Use(middleware.Logger())
// Output: [abc-123] GET /users 200 15ms
```

### SecureHeaders()
Adds security headers:

```go
router.Use(middleware.SecureHeaders())
// X-Content-Type-Options: nosniff
// X-Frame-Options: DENY
// X-XSS-Protection: 1; mode=block
```

## Usage in main.go

```go
router := gin.Default()

// Global middleware
router.Use(middleware.CORS())
router.Use(middleware.RequestID())
router.Use(middleware.Logger())

// Route-specific middleware
protected := router.Group("/api")
protected.Use(middleware.Auth(cfg.JWT.AccessSecret))
```

## Adding New Middleware

```go
func RateLimiter(maxRequests int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Rate limiting logic
        c.Next()
    }
}
```
