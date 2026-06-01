# internal

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains private application code. Code in the `internal` folder cannot be imported by other projects (enforced by Go compiler).

## Structure

```
internal/
├── config/           # Configuration loader
├── domain/           # Business entities & interfaces
├── repository/       # Data access implementation
├── usecase/          # Business logic
├── delivery/         # HTTP/Transport layer
├── middleware/       # HTTP middleware
├── infrastructure/   # External services
└── shared/           # Shared utilities
```

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      delivery/                              │
│               (HTTP Handlers, Routes)                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      usecase/                               │
│                  (Business Logic)                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      domain/                                │
│              (Entities, Interfaces)                         │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │
┌─────────────────────────────────────────────────────────────┐
│                    repository/                              │
│              (Data Access Implementation)                   │
└─────────────────────────────────────────────────────────────┘
```

## Dependency Rule

- **Delivery** depends on **Usecase**
- **Usecase** depends on **Domain** (interfaces only)
- **Repository** implements **Domain** interfaces
- **Domain** has NO external dependencies

## Folder Details

| Folder | Layer | Description |
|--------|-------|-------------|
| `config/` | Support | Load configuration from environment |
| `domain/` | Core | Entities, repository interfaces, errors |
| `repository/` | Data | GORM implementation, DB models |
| `usecase/` | Application | Business logic per-feature |
| `delivery/` | Presentation | HTTP handlers, routes, DTOs |
| `middleware/` | Support | Auth, CORS, logging middleware |
| `infrastructure/` | Support | Database connection, cache, email |
| `shared/` | Support | Utilities, constants, helpers |

## Best Practices

- ✅ Domain layer must be pure (no framework imports)
- ✅ Usecase only depends on domain interfaces
- ✅ Repository maps between domain entity and DB model
- ✅ Handler only handles HTTP, delegates to usecase
- ❌ Don't import usecase from repository
- ❌ Don't import delivery from domain
