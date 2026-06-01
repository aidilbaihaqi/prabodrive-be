# migrations

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains database migration files.

## Files

| File | Description |
|------|-------------|
| `000001_init.up.sql` | Migration to create initial schema |
| `000001_init.down.sql` | Rollback migration |

## Naming Format

```
{version}_{description}.{direction}.sql

Example:
000001_init.up.sql
000001_init.down.sql
000002_add_products_table.up.sql
000002_add_products_table.down.sql
```

## Usage

### With Makefile

```bash
# Run all pending migrations
make migrate-up

# Rollback one migration
make migrate-down

# Create new migration
make migrate-create name=add_products_table
```

### With golang-migrate

```bash
# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback
migrate -path migrations -database "$DATABASE_URL" down 1

# Force version (if stuck)
migrate -path migrations -database "$DATABASE_URL" force 1
```

### With script

```bash
./scripts/migrate.sh up
./scripts/migrate.sh down
./scripts/migrate.sh create add_products
```

## Migration Structure

### Up Migration (000001_init.up.sql)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(254) UNIQUE NOT NULL,
    -- ...
);

CREATE INDEX idx_users_email ON users(email);
```

### Down Migration (000001_init.down.sql)
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## Best Practices

- ✅ Always create up and down migrations
- ✅ Migrations must be idempotent
- ✅ Test rollback before deploying
- ✅ Backup database before production migration
- ❌ Never edit already-applied migrations
- ❌ Don't delete data in up migrations (use soft delete)
