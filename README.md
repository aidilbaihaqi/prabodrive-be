# Go Backend Template

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-green.svg)](docs/ARCHITECTURE.md)

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

A production-ready Go backend template that combines the organizational excellence of **golang-standards/project-layout** with the architectural rigor of **Clean Architecture**. This template provides a solid foundation for building scalable, maintainable, and testable backend applications.

---

## 🌟 Why This Template?

Starting a new Go backend project often means spending hours setting up folder structures, configurations, middleware, and boilerplate code. This template eliminates that friction by providing a battle-tested structure that enforces best practices from day one.

The architecture follows the **Dependency Rule**: dependencies only point inward, ensuring that your business logic remains independent of frameworks, databases, and external services. This makes your codebase highly testable and adaptable to change.

---

## 📁 Project Structure Overview

```
go-backend-template/
├── api/                    # API specifications (OpenAPI/Swagger)
├── build/                  # Build and packaging configurations
│   └── docker/
├── cmd/                    # Application entry points
│   └── api/
├── configs/                # Configuration file templates
├── deployments/            # Container orchestration configs
├── docs/                   # Project documentation
├── internal/               # Private application code
│   ├── config/
│   ├── domain/
│   ├── repository/
│   ├── usecase/
│   ├── delivery/
│   ├── middleware/
│   ├── infrastructure/
│   └── shared/
├── migrations/             # Database migration files
├── pkg/                    # Public reusable packages
├── scripts/                # Utility scripts
└── test/                   # Additional test resources
```

---

## 📂 Detailed Folder Descriptions

### `api/`
This folder contains your API specifications, primarily OpenAPI (Swagger) documents that define your REST endpoints, request/response schemas, and authentication requirements. Having a centralized API spec ensures that your frontend team, API consumers, and documentation stay synchronized with your actual implementation.

The OpenAPI specification serves as a single source of truth for your API contract. You can use tools like Swagger UI to visualize your endpoints, generate client SDKs automatically, or validate requests against the schema during development and testing.

### `build/`
The build directory houses all configurations related to packaging and continuous integration. This includes Dockerfiles, CI/CD pipeline configurations, and any scripts needed to create distributable artifacts. Keeping build concerns separate from application code makes it easier to maintain and update your deployment pipeline.

Currently, it contains a multi-stage Dockerfile optimized for production: small image size (~20MB), non-root user for security, and built-in health checks. You can extend this folder with additional Dockerfiles for development or configurations for cloud-specific builds.

### `cmd/`
This folder contains the main entry points for your application. Each subdirectory represents an executable binary. The `api/` subdirectory contains the main HTTP server, which is responsible for wiring up all dependencies (configuration, database, repositories, use cases, and handlers) and starting the server.

Following the principle of keeping `main.go` thin, the entry point only handles dependency injection and server initialization. All business logic resides in the `internal/` directory, making the application easier to test and maintain.

### `configs/`
Configuration templates and sample files live here. This includes example configuration files that demonstrate the expected structure and values for different environments (development, staging, production). These templates help new developers understand what configuration options are available.

Unlike `.env` files which contain actual secrets, this folder contains template files (like `config.example.yaml`) that can be safely committed to version control. Real secrets should never be stored in the repository.

### `deployments/`
Container orchestration and deployment configurations reside in this folder. It includes Docker Compose files for local development, Kubernetes manifests for production deployment, and infrastructure-as-code templates (Terraform, CloudFormation).

The included `docker-compose.yml` provides a complete local development stack with PostgreSQL, Redis, and the application itself. This allows developers to spin up the entire environment with a single command, ensuring consistency across development machines.

### `docs/`
Project documentation beyond code comments lives here. This includes architecture decision records (ADRs), API documentation, development guides, and best practices. Well-maintained documentation reduces onboarding time and serves as a reference for architectural decisions.

Currently includes `ARCHITECTURE.md` (explaining Clean Architecture implementation) and `BEST_PRACTICES.md` (coding standards and patterns). As your project grows, consider adding deployment guides, troubleshooting documents, and contribution guidelines.

### `internal/`
The heart of your application, containing all private code that should not be imported by external projects. Go enforces this at the compiler level: code inside `internal/` cannot be imported from outside the module. This is where Clean Architecture layers live.

The Clean Architecture layers enforce a strict dependency rule: **Domain** (entities, interfaces) → **Use Cases** (business logic) → **Delivery** (HTTP handlers) and **Repository** (data access). Each layer only depends on inner layers, never outer ones, making the codebase highly modular and testable.

#### `internal/config/`
Application configuration management. This package loads settings from environment variables and provides typed access to configuration values. Using structured configuration (instead of scattered `os.Getenv` calls) makes it easier to validate settings at startup and prevents runtime errors from missing configuration.

#### `internal/domain/`
The innermost layer containing pure business entities and repository interfaces. This layer has **zero external dependencies** — no frameworks, no databases, just plain Go structs and interfaces. Entities contain business logic methods (validation, calculations), while interfaces define contracts that outer layers must implement.

#### `internal/repository/`
Data access layer implementing the interfaces defined in the domain. Each repository handles database operations (CRUD, queries) and translates between domain entities and database models. Keeping database models separate from domain entities ensures that ORM-specific concerns don't leak into business logic.

#### `internal/usecase/`
Application business logic organized by feature. Each file represents a single operation (Create, Read, Update, Delete) following the Single Responsibility Principle. Use cases orchestrate data flow: they validate input, apply business rules, interact with repositories, and return structured output.

#### `internal/delivery/`
Transport layer handling HTTP requests and responses. This includes Gin handlers, route registration, request DTOs (with validation tags), and standardized response formatting. Handlers are thin — they parse requests, call use cases, and format responses without containing business logic.

#### `internal/middleware/`
HTTP middleware for cross-cutting concerns like authentication (JWT validation), CORS handling, request logging, rate limiting, and security headers. Middleware centralizes these concerns, keeping handlers focused on their primary responsibility.

#### `internal/infrastructure/`
External service connections: database clients, cache clients (Redis), email services, cloud storage, and message queues. This layer abstracts infrastructure details, allowing you to swap implementations (e.g., PostgreSQL to MySQL) without affecting business logic.

#### `internal/shared/`
Shared utilities and constants used across the application. This includes helper functions (password hashing, string manipulation), JWT token utilities, and application-wide constants. Code here should be stateless and have no dependencies on other layers.

### `migrations/`
Database schema migration files using a versioned approach. Each migration has an "up" file (apply changes) and a "down" file (rollback changes). Sequential numbering (000001, 000002) ensures migrations run in the correct order and provides clear history of schema evolution.

Using SQL migrations (rather than ORM auto-migration) gives you full control over schema changes and allows for complex operations like data transformations. Always write and test rollback scripts before deploying to production.

### `pkg/`
Public packages that can be imported by external projects. Use this folder sparingly — only for truly reusable libraries that provide value beyond your application. Unlike `internal/`, code in `pkg/` has no import restrictions.

Consider using `pkg/` for generic utilities like custom validators, logging wrappers, or HTTP client helpers that other projects in your organization might benefit from. If in doubt, start in `internal/` and move to `pkg/` only when reuse becomes necessary.

### `scripts/`
Shell scripts for common development tasks: project setup, database migrations, code generation, deployment automation. These scripts complement the Makefile by handling more complex multi-step operations that benefit from shell scripting.

Includes `setup.sh` (first-time project setup) and `migrate.sh` (database migration helper). For Windows users, most operations are also available through the Makefile which calls Go commands directly.

### `test/`
Additional test resources beyond unit tests. This includes integration tests, end-to-end tests, test fixtures (sample data), and test utilities. Unit tests live next to the code they test (e.g., `create_test.go`), while broader tests live here.

Separating integration tests from unit tests allows different test strategies: unit tests run quickly on every commit, while integration tests run less frequently and may require external services.

---

## 🚀 Key Features

| Feature | Description |
|---------|-------------|
| **Clean Architecture** | Strict layer separation with dependency inversion |
| **Self-Documenting** | 22 README files across all directories |
| **API Versioning** | `/api/v1` prefix with OpenAPI specification |
| **JWT Authentication** | Access & refresh token support with middleware |
| **Database Migrations** | Versioned SQL migrations with rollback support |
| **Docker Ready** | Multi-stage build, Docker Compose for local dev |
| **Comprehensive Makefile** | 20+ commands for building, testing, deploying |
| **Security Best Practices** | Bcrypt hashing, secure headers, CORS middleware |

---

## 🛠️ Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Setup

```bash
# Clone the repository
git clone https://github.com/aidilbaihaqi/go-projects-layout.git
cd go-projects-layout

# Copy environment file
cp .env.example .env

# Option 1: Run with Docker
make docker-up

# Option 2: Run locally
make migrate-up
make run
```

### Verify Installation

```bash
curl http://localhost:8080/health
# Expected: {"status":"healthy"}
```

---

## 📋 Available Commands

```bash
make help           # Show all available commands

# Development
make run            # Run the application
make dev            # Run with hot reload (requires air)
make build          # Build the binary

# Testing
make test           # Run all tests
make test-cover     # Run tests with coverage report

# Database
make migrate-up     # Apply all pending migrations
make migrate-down   # Rollback last migration
make migrate-create name=add_products  # Create new migration

# Docker
make docker-build   # Build Docker image
make docker-up      # Start all containers
make docker-down    # Stop all containers

# Code Quality
make lint           # Run golangci-lint
make fmt            # Format code with gofmt
```

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/ARCHITECTURE.md) | Clean Architecture implementation details |
| [Best Practices](docs/BEST_PRACTICES.md) | Coding standards and patterns |
| [API Specification](api/openapi.yaml) | OpenAPI 3.0 documentation |

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Aidil Baihaqi**  
- GitHub: [@aidilbaihaqi](https://github.com/aidilbaihaqi)
- LinkedIn: [Connect with me](https://linkedin.com/in/aidilbaihaqi)

---

⭐ If you find this template helpful, please give it a star!
