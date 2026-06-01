# http

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains HTTP delivery layer.

## Structure

```
http/
├── handler/           # HTTP request handlers
│   └── user_handler.go
├── routes/            # Route registration
│   └── routes.go
├── request/           # Request DTOs
│   └── user_request.go
└── response/          # Response utilities
    └── response.go
```

## Flow

```
Request → Route → Middleware → Handler → Usecase → Response
```

## Subfolders

### handler/
HTTP handlers that:
- Parse request body
- Validate input
- Call usecase
- Format response

### routes/
Route registration:
- API versioning (`/api/v1`)
- Middleware assignment
- Group routes by resource

### request/
Request DTOs with Gin validation:
```go
type CreateUserRequest struct {
    Email string `json:"email" binding:"required,email"`
}
```

### response/
Standardized response format:
```go
{
    "success": true,
    "data": {...},
    "meta": {...}
}
```

## Best Practices

- ✅ One handler per resource
- ✅ Consistent response format
- ✅ Proper HTTP status codes
- ✅ Detailed validation messages
