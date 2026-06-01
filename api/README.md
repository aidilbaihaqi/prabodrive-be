# api

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains API specifications and technical documentation.

## Files

| File | Description |
|------|-------------|
| `openapi.yaml` | OpenAPI 3.0 specification for REST API |

## Usage

### OpenAPI/Swagger

The `openapi.yaml` file defines:
- API Endpoints
- Request/Response Schemas
- Authentication Requirements
- Error Responses

### Generate Swagger UI

```bash
# Install swaggo
go install github.com/swaggo/swag/cmd/swag@latest

# Generate from code annotations
swag init -g cmd/api/main.go -o api/
```

### Useful Tools

- [Swagger Editor](https://editor.swagger.io/) - Edit and preview OpenAPI
- [Redoc](https://github.com/Redocly/redoc) - Generate documentation
- [OpenAPI Generator](https://openapi-generator.tech/) - Generate client SDKs

## Best Practices

- ✅ Always update spec when endpoints change
- ✅ Use `$ref` for reusable schemas
- ✅ Document all error responses
- ✅ Add examples for each endpoint
