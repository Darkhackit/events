package dto

type PermissionResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PermissionRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AssignPermissionRequest struct {
	PermissionID []PermissionRequest `json:"permission_id"`
	RoleID       uint                `json:"role_id"`
}
