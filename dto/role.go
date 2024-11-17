package dto

type RoleRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}
type RoleResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserRoleRequest struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}

type RolePermissionResponse struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Permissions []PermissionResponse
}
