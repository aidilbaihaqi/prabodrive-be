# user

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains usecases for User feature.

## Files

| File | Description |
|------|-------------|
| `create.go` | Create new user |
| `get.go` | Get user by ID |
| `list.go` | List users with pagination |
| `update.go` | Update user |
| `delete.go` | Delete user (soft delete) |

## create.go

**CreateUserUsecase** - Creates a new user

Input:
- `Email` (required) - User email
- `Password` (required) - Plain text password
- `Name` (required) - User name
- `Role` (optional) - Default: "user"

Output:
- `ID`, `Email`, `Name`, `Role`, `CreatedAt`

Flow:
1. Validate input
2. Check email exists
3. Hash password
4. Create user entity
5. Save via repository
6. Return output

## get.go

**GetUserUsecase** - Get user by ID

Input: `ID` (uuid)  
Output: User details

## list.go

**ListUsersUsecase** - List users with pagination

Input:
- `Page` - Page number (default: 1)
- `Limit` - Items per page (default: 10, max: 100)
- `Name` - Filter by name
- `Role` - Filter by role

Output:
- `Users` - Array of users
- `Total` - Total count
- `Page`, `Limit`, `TotalPages` - Pagination info

## update.go

**UpdateUserUsecase** - Update user

Input: All fields optional (partial update)
- `Name`
- `Email`
- `Password`
- `Role`
- `IsActive`

## delete.go

**DeleteUserUsecase** - Soft delete user

Input: `ID` (uuid)  
Output: None (success/error)

## Dependency

All usecases depend on `domain.UserRepository` interface.

```go
type CreateUserUsecase struct {
    userRepo domain.UserRepository
}
```
