package constants

// User roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// Pagination defaults
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Response messages
const (
	MsgSuccess       = "Success"
	MsgCreated       = "Resource created successfully"
	MsgUpdated       = "Resource updated successfully"
	MsgDeleted       = "Resource deleted successfully"
	MsgUnauthorized  = "Unauthorized"
	MsgForbidden     = "Access forbidden"
	MsgNotFound      = "Resource not found"
	MsgBadRequest    = "Bad request"
	MsgInternalError = "Internal server error"
)

// Context keys
const (
	ContextKeyUserID    = "user_id"
	ContextKeyUserRole  = "user_role"
	ContextKeyRequestID = "request_id"
)

// HTTP headers
const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderRequestID     = "X-Request-ID"
)

// Date/time formats
const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
	ISO8601Format  = "2006-01-02T15:04:05Z07:00"
)

// Validation constraints
const (
	MinPasswordLength = 8
	MaxPasswordLength = 72 // bcrypt limit
	MinNameLength     = 2
	MaxNameLength     = 100
	MaxEmailLength    = 254
)
