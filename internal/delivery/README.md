# delivery

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains transport/presentation layer. Handles HTTP requests and responses.

## Structure

```
delivery/
└── http/
    ├── handler/
    │   └── user_handler.go    # HTTP handlers
    ├── routes/
    │   └── routes.go          # Route registration
    ├── request/
    │   └── user_request.go    # Request DTOs
    └── response/
        └── response.go        # Response utilities
```

## Files

### handler/user_handler.go

HTTP handlers that connect HTTP requests with usecases:

```go
type UserHandler struct {
    createUserUC *user.CreateUserUsecase
    getUserUC    *user.GetUserUsecase
    // ...
}

func (h *UserHandler) Create(c *gin.Context) {
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    result, err := h.createUserUC.Execute(input)
    // Handle response
}
```

### routes/routes.go

Register endpoints:

```go
func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler, jwtSecret string) {
    v1 := r.Group("/api/v1")
    
    users := v1.Group("/users")
    users.Use(middleware.Auth(jwtSecret))
    {
        users.POST("", h.Create)
        users.GET("", h.List)
        users.GET("/:id", h.Get)
        users.PUT("/:id", h.Update)
        users.DELETE("/:id", h.Delete)
    }
}
```

### request/user_request.go

Request DTOs with validation tags:

```go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required,min=2"`
}
```

### response/response.go

Standard response format:

```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

func Success(c *gin.Context, data interface{})
func Error(c *gin.Context, code int, message string)
func NotFound(c *gin.Context, message string)
```

## Adding a New Handler

1. Create `handler/product_handler.go`
2. Create `request/product_request.go`
3. Register in `routes/routes.go`

## Best Practices

- ✅ Validate request in handler
- ✅ Map request to usecase input
- ✅ Handle errors with switch/case
- ✅ Use consistent response format
- ❌ No business logic in handler
- ❌ Don't access repository from handler
