# usecase

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains application business logic. Each usecase handles one specific operation.

## Structure

```
usecase/
└── user/
    ├── create.go     # Create user logic
    ├── get.go        # Get single user
    ├── list.go       # List users with pagination
    ├── update.go     # Update user
    └── delete.go     # Delete user
```

## Pattern

Each usecase has:
1. **Struct** - stores dependencies
2. **Input DTO** - required data
3. **Output DTO** - returned data
4. **Execute()** - main method

```go
// usecase/user/create.go

type CreateUserUsecase struct {
    userRepo domain.UserRepository
}

type CreateUserInput struct {
    Email    string
    Password string
    Name     string
}

type CreateUserOutput struct {
    ID    uuid.UUID
    Email string
    Name  string
}

func (uc *CreateUserUsecase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Validate input
    // 2. Check business rules
    // 3. Create entity
    // 4. Save via repository
    // 5. Return output
}
```

## Files

| File | Description |
|------|-------------|
| `create.go` | Create new user with validation and password hashing |
| `get.go` | Get user details by ID |
| `list.go` | List users with pagination and filtering |
| `update.go` | Update user with partial update support |
| `delete.go` | Soft delete user |

## Adding a New Feature

For Product feature, create folder `usecase/product/`:

```
usecase/
├── user/
│   └── ...
└── product/
    ├── create.go
    ├── get.go
    ├── list.go
    ├── update.go
    └── delete.go
```

## Best Practices

### ✅ DO
- One usecase = one operation (Single Responsibility)
- Validate input at the start of Execute()
- Use domain entities
- Return Output DTO, not domain entity

### ❌ DON'T
- Don't combine multiple operations
- Don't import HTTP packages (gin, http)
- Don't access database directly (use repository)
- Don't return GORM models
