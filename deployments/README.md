# deployments

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains deployment and orchestration configurations.

## Files

| File | Description |
|------|-------------|
| `docker-compose.yml` | Local development stack |

## docker-compose.yml

Runs all required services:
- **app** - Go API server
- **postgres** - PostgreSQL database
- **redis** - Redis cache
- **adminer** - Database management UI (optional)

## Usage

```bash
# Start all services
docker-compose -f deployments/docker-compose.yml up -d

# View logs
docker-compose -f deployments/docker-compose.yml logs -f

# Stop all services
docker-compose -f deployments/docker-compose.yml down

# Stop and remove volumes
docker-compose -f deployments/docker-compose.yml down -v
```

Or use Makefile:
```bash
make docker-up
make docker-down
make docker-logs
```

## Service Access

| Service | URL | Credentials |
|---------|-----|-------------|
| API | http://localhost:8080 | - |
| Adminer | http://localhost:8081 | postgres/postgres |
| PostgreSQL | localhost:5432 | postgres/postgres |
| Redis | localhost:6379 | - |

## Extended Structure (Optional)

```
deployments/
├── docker-compose.yml          # Local development
├── docker-compose.prod.yml     # Production
├── kubernetes/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
└── terraform/
    ├── main.tf
    └── variables.tf
```

## Best Practices

- ✅ Use environment variables for secrets
- ✅ Health checks for all services
- ✅ Named volumes for data persistence
- ❌ Don't commit production secrets
