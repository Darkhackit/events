package service

import (
	"context"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
)

type PermissionService interface {
	CreatePermission(ctx context.Context, permission dto.PermissionRequest) (*dto.PermissionResponse, error)
	GetPermissions(ctx context.Context) ([]dto.PermissionResponse, error)
	AssignPermission(ctx context.Context, permission dto.AssignPermissionRequest) error
}

type DefaultPermissionService struct {
	repo domain.PermissionRepository
}

func (dp *DefaultPermissionService) CreatePermission(ctx context.Context, permission dto.PermissionRequest) (*dto.PermissionResponse, error) {
	arg := domain.Permission{
		Name: permission.Name,
	}
	result, err := dp.repo.CreatePermission(ctx, arg)
	if err != nil {
		return nil, err
	}
	perm := dto.PermissionResponse{
		Name: result.Name,
		ID:   uint(result.ID),
	}
	return &perm, nil
}
func (dp *DefaultPermissionService) GetPermissions(ctx context.Context) ([]dto.PermissionResponse, error) {
	result, err := dp.repo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}
	perms := make([]dto.PermissionResponse, len(result))
	for i, perm := range result {
		perms[i] = dto.PermissionResponse{
			Name: perm.Name,
			ID:   uint(perm.ID),
		}
	}
	return perms, nil
}
func (dp *DefaultPermissionService) AssignPermission(ctx context.Context, permission dto.AssignPermissionRequest) error {
	err := dp.repo.AssignRoles(ctx, permission)
	if err != nil {
		return err
	}
	return nil
}

func NewPermissionService(repo domain.PermissionRepository) *DefaultPermissionService {
	return &DefaultPermissionService{repo: repo}
}
