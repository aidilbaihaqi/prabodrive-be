package request

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}
