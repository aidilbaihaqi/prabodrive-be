package request

type CreateFolderRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id"`
}

type UpdateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

// ParentID: nil = all folders; "root" = root-level only; uuid = children of that folder
type ListFoldersQuery struct {
	ParentID *string `form:"parent_id"`
}
