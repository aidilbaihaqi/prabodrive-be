# config

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains application configuration.

## Files

| File | Description |
|------|-------------|
| `config.go` | Configuration loader from environment variables |

## config.go

Struct-based configuration that loads:
- **App Config** - name, environment, port, debug mode
- **Database Config** - host, port, name, credentials
- **JWT Config** - secrets, expiry times
- **Redis Config** - connection settings
- **SMTP Config** - email settings

## Usage

```go
// Load configuration
cfg := config.Load()

// Access config values
port := cfg.AppPort
dbHost := cfg.Database.Host
jwtSecret := cfg.JWT.AccessSecret
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | development | Environment mode |
| `APP_PORT` | 8080 | Server port |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_NAME` | myapp_db | Database name |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `JWT_ACCESS_SECRET` | secret | JWT access token secret |
| `JWT_REFRESH_SECRET` | secret | JWT refresh token secret |

## Best Practices

- ✅ Use typed struct for config
- ✅ Provide sensible defaults for development
- ✅ Validate required config at startup
- ❌ Don't hardcode secrets
- ❌ Don't commit .env file
