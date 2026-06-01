# infrastructure

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains connections to external services.

## Structure

```
infrastructure/
└── database/
    └── postgres.go    # PostgreSQL connection
```

## Files

### database/postgres.go

Database connection with GORM:

```go
func NewPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Connection pool settings
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}
```

## Usage

```go
// main.go
db, err := database.NewPostgres(cfg.Database)
if err != nil {
    log.Fatal(err)
}

// Pass to repositories
userRepo := repository.NewUserRepository(db)
```

## Extended Structure

For more complex applications:

```
infrastructure/
├── database/
│   ├── postgres.go      # PostgreSQL
│   └── migrations.go    # Auto migrations
├── cache/
│   └── redis.go         # Redis client
├── email/
│   └── smtp.go          # SMTP client
├── storage/
│   └── s3.go            # S3/Cloud storage
└── queue/
    └── rabbitmq.go      # Message queue
```

## Best Practices

- ✅ Connection pooling
- ✅ Health checks
- ✅ Graceful shutdown
- ✅ Retry logic
- ❌ Don't hardcode credentials
