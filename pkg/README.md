# pkg

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains library code that **can be imported by external projects**.

## Difference between pkg and internal

| Folder | Can be imported externally? | Use case |
|--------|---------------------------|----------|
| `internal/` | ❌ No | Private application code |
| `pkg/` | ✅ Yes | Reusable libraries |

## When to Use pkg/

Use `/pkg` when:
- ✅ Code is intended for use by other projects
- ✅ Building a reusable library
- ✅ Clearly separating internal vs public code

## When NOT to Use pkg/

Don't use `/pkg` when:
- ❌ Project is small
- ❌ All code is internal
- ❌ No one will import from outside

## Example Structure

```
pkg/
├── logger/
│   ├── logger.go        # Custom logger
│   └── logger_test.go
├── validator/
│   ├── validator.go     # Custom validators
│   └── validator_test.go
├── httpclient/
│   ├── client.go        # HTTP client wrapper
│   └── client_test.go
└── pagination/
    ├── pagination.go    # Pagination helpers
    └── pagination_test.go
```

## Import from External Project

```go
import "github.com/yourname/yourapp/pkg/logger"
import "github.com/yourname/yourapp/pkg/validator"
```

## Best Practices

- ✅ Write comprehensive documentation
- ✅ Include unit tests
- ✅ Semantic versioning
- ✅ Backward compatibility
- ❌ Don't import from internal/
- ❌ Don't depend on application-specific code

## Note

The `/pkg` folder is **optional**. Many successful Go projects don't use it. Consider using `/internal` for all code if sharing with other projects isn't needed.
