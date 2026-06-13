package request

type PresignUploadRequest struct {
	Name     string  `json:"name" binding:"required"`
	Size     int64   `json:"size" binding:"required,min=1"`
	MIMEType string  `json:"mime_type" binding:"required"`
	FolderID *string `json:"folder_id"`
}

type ConfirmUploadRequest struct {
	S3Key    string  `json:"s3_key" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Size     int64   `json:"size" binding:"required,min=1"`
	MIMEType string  `json:"mime_type" binding:"required"`
	FolderID *string `json:"folder_id"`
}

type RenameDocumentRequest struct {
	Name string `json:"name" binding:"required"`
}

type ListDocumentsQuery struct {
	FolderID *string `form:"folder_id"`
	Search   string  `form:"search"`
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=20"`
}
