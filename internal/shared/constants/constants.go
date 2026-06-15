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

	MaxFileSize = 100 * 1024 * 1024 // 100 MB in bytes

	DefaultQuotaMax = 3 * 1024 * 1024 * 1024 // 3 GB in bytes

	ActionUpload       = "upload"
	ActionDownload     = "download"
	ActionRename       = "rename"
	ActionDelete       = "delete"
	ActionCreateFolder = "create_folder"
	ActionRenameFolder = "rename_folder"
	ActionDeleteFolder = "delete_folder"
	ActionShareCreate  = "share_create"
	ActionShareAccess  = "share_access"
	ActionShareDelete  = "share_delete"
	ActionLogin        = "login"
	ActionLogout       = "logout"
)
