# shared

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains shared code used across the application.

## Structure

```
shared/
├── utils/
│   └── utils.go        # Utility functions
├── token/
│   └── jwt.go          # JWT utilities
└── constants/
    └── constants.go    # Application constants
```

## Files

### utils/utils.go

Common utility functions:

| Function | Description |
|----------|-------------|
| `HashPassword(password)` | Hash password with bcrypt |
| `CheckPassword(hash, password)` | Verify password |
| `GenerateRandomString(length)` | Generate random string |
| `GenerateOTP(length)` | Generate numeric OTP |
| `IsValidEmail(email)` | Validate email format |
| `IsValidPassword(password)` | Check password strength |
| `Slugify(s)` | Convert to URL-friendly slug |
| `SanitizeString(s)` | Remove dangerous characters |

### token/jwt.go

JWT token utilities:

| Function | Description |
|----------|-------------|
| `GenerateAccessToken(...)` | Create access token |
| `GenerateRefreshToken(...)` | Create refresh token |
| `GenerateTokenPair(...)` | Create both tokens |
| `ValidateAccessToken(token, secret)` | Validate & parse access token |
| `ValidateRefreshToken(token, secret)` | Validate & parse refresh token |

### constants/constants.go

Application-wide constants:

```go
// User roles
const (
    RoleUser  = "user"
    RoleAdmin = "admin"
)

// Pagination defaults
const (
    DefaultPage  = 1
    DefaultLimit = 10
    MaxLimit     = 100
)

// Date formats
const (
    DateFormat = "2006-01-02"
    TimeFormat = "15:04:05"
)
```

## Usage

```go
import "github.com/yourname/yourapp/internal/shared/utils"
import "github.com/yourname/yourapp/internal/shared/token"
import "github.com/yourname/yourapp/internal/shared/constants"

// Hash password
hashed, _ := utils.HashPassword("password123")

// Generate token
accessToken, _, _ := token.GenerateAccessToken(userID, email, role, secret, expiry)

// Use constants
if role == constants.RoleAdmin { ... }
```

## Best Practices

- ✅ Pure functions without side effects
- ✅ Well-documented functions
- ✅ Unit tests for each function
- ❌ No state or global variables
- ❌ No dependency on other layers
