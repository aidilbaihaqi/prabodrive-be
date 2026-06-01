# build

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains configurations for packaging and CI/CD.

## Structure

```
build/
└── docker/
    └── Dockerfile    # Multi-stage Docker build
```

## Files

### docker/Dockerfile

Multi-stage Dockerfile for production:
- **Stage 1 (Builder)**: Compile Go binary
- **Stage 2 (Final)**: Alpine-based minimal image

Features:
- Build binary without CGO
- Non-root user for security
- Built-in health check
- Minimal image size (~20MB)

## Usage

```bash
# Build image
docker build -t myapp:latest -f build/docker/Dockerfile .

# Run container
docker run -p 8080:8080 myapp:latest
```

## Extended Structure (Optional)

For more complex projects, you can add:

```
build/
├── ci/
│   ├── github-actions.yml
│   └── gitlab-ci.yml
├── docker/
│   ├── Dockerfile
│   └── Dockerfile.dev
└── package/
    ├── deb/
    └── rpm/
```

## Best Practices

- ✅ Use multi-stage builds for small images
- ✅ Pin dependency versions in Dockerfile
- ✅ Don't include secrets in image
- ✅ Use non-root user
