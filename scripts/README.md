# scripts

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains shell scripts for various operations.

## Files

| File | Description |
|------|-------------|
| `setup.sh` | Setup development environment |
| `migrate.sh` | Helper script for database migrations |

## setup.sh

Script for initial development environment setup:
- Check Go installation
- Copy .env.example to .env
- Download Go dependencies
- Install development tools (migrate, air, golangci-lint)

```bash
# Run setup
chmod +x scripts/setup.sh
./scripts/setup.sh
```

## migrate.sh

Helper script for running migrations easily:

```bash
chmod +x scripts/migrate.sh

# Run all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Rollback all migrations
./scripts/migrate.sh down-all

# Show current version
./scripts/migrate.sh version

# Create new migration
./scripts/migrate.sh create add_products_table

# Force set version
./scripts/migrate.sh force 1
```

## Additional Scripts (Optional)

```
scripts/
├── setup.sh           # Initial setup
├── migrate.sh         # Migration helper
├── test.sh            # Run tests
├── build.sh           # Build for all platforms
├── deploy.sh          # Deployment script
├── backup-db.sh       # Database backup
└── seed.sh            # Seed database with test data
```

## Windows Users

For Windows, you can use:
- Git Bash
- WSL (Windows Subsystem for Linux)
- PowerShell equivalent scripts

Or run commands directly via Makefile:
```bash
make migrate-up
make migrate-down
```
