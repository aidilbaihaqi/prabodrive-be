package constants

const (
	ContextUserID    = "user_id"
	ContextUserEmail = "user_email"
	ContextUserRole  = "user_role"

	RoleAdmin = "admin"
	RoleUser  = "user"

	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 50

	MaxFileSize = 5 * 1024 * 1024 // 5 MB in bytes

	DefaultQuotaMax = 3 * 1024 * 1024 * 1024 // 3 GB in bytes

	ActionUpload      = "upload"
	ActionDownload    = "download"
	ActionShareCreate = "share_create"
	ActionShareAccess = "share_access"
	ActionDelete      = "delete"
	ActionLogin       = "login"
	ActionLogout      = "logout"
)
