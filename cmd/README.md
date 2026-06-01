# cmd

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains application entry points.

## Structure

```
cmd/
└── api/
    └── main.go    # Entry point for API server
```

## Files

### api/main.go

Main function that runs the HTTP server. Responsible for:

1. Load environment variables
2. Load configuration
3. Initialize infrastructure (database, cache)
4. Initialize repositories
5. Initialize usecases
6. Initialize handlers
7. Setup router & middleware
8. Start HTTP server

## Dependency Injection

```go
// main.go pattern
func main() {
    // 1. Config
    cfg := config.Load()
    
    // 2. Infrastructure
    db := database.NewPostgres(cfg.Database)
    
    // 3. Repositories
    userRepo := repository.NewUserRepository(db)
    
    // 4. Usecases
    createUserUC := user.NewCreateUserUsecase(userRepo)
    
    // 5. Handlers
    userHandler := handler.NewUserHandler(createUserUC)
    
    // 6. Router
    router := gin.Default()
    routes.RegisterUserRoutes(router, userHandler)
    
    // 7. Start
    router.Run(":8080")
}
```

## Adding New Entry Points

For applications with multiple binaries:

```
cmd/
├── api/
│   └── main.go      # REST API server
├── worker/
│   └── main.go      # Background worker
├── migrate/
│   └── main.go      # Migration tool
└── cli/
    └── main.go      # CLI application
```

## Best Practices

- ✅ Keep main.go minimal - only wire dependencies
- ✅ All business logic in internal/
- ✅ Use graceful shutdown
- ❌ Don't put logic in main.go
