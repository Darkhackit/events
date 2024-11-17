package domain

import (
	"context"
	"github.com/Darkhackit/events/dto"
)

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PermissionRepository interface {
	CreatePermission(ctx context.Context, permission Permission) (*Permission, error)
	GetPermissions(ctx context.Context) ([]Permission, error)
	AssignRoles(ctx context.Context, permissions dto.AssignPermissionRequest) error
}
