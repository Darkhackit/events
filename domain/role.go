package domain

import (
	"context"
	"github.com/Darkhackit/events/dto"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoleRepository interface {
	CreateRole(ctx context.Context, role Role) (*Role, error)
	UpdateRole(ctx context.Context, role Role) (*Role, error)
	DeleteRole(ctx context.Context, id int) error
	GetRole(ctx context.Context, id int) (*Role, error)
	AssignRoleToUser(ctx context.Context, UserRole dto.UserRoleRequest) error
	GetRoles(ctx context.Context) ([]*dto.RolePermissionResponse, error)
}
