package service

import (
	"context"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
)

type RoleService interface {
	CreateRole(ctx context.Context, role dto.RoleRequest) (*dto.RoleResponse, error)
	UpdateRole(ctx context.Context, role dto.RoleRequest) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id int) error
	GetRole(ctx context.Context, id int) (*dto.RoleResponse, error)
	AssignRoleToUser(ctx context.Context, UserRole dto.UserRoleRequest) error
	GetRoles(ctx context.Context) (*[]dto.RolePermissionResponse, error)
}

type DefaultRoleService struct {
	repo domain.RoleRepository
}

func (rs *DefaultRoleService) CreateRole(ctx context.Context, role dto.RoleRequest) (*dto.RoleResponse, error) {
	arg := domain.Role{
		Name: role.Name,
	}
	result, err := rs.repo.CreateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	response := dto.RoleResponse{
		Name: role.Name,
		ID:   result.ID,
	}
	return &response, nil
}
func (rs *DefaultRoleService) UpdateRole(ctx context.Context, role dto.RoleRequest) (*dto.RoleResponse, error) {
	arg := domain.Role{
		Name: role.Name,
		ID:   role.ID,
	}
	result, err := rs.repo.UpdateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	response := dto.RoleResponse{
		Name: result.Name,
		ID:   result.ID,
	}
	return &response, nil
}
func (rs *DefaultRoleService) DeleteRole(ctx context.Context, id int) error {
	err := rs.repo.DeleteRole(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (rs *DefaultRoleService) GetRole(ctx context.Context, id int) (*dto.RoleResponse, error) {
	role, err := rs.repo.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}
	response := dto.RoleResponse{
		Name: role.Name,
		ID:   role.ID,
	}
	return &response, nil
}
func (rs *DefaultRoleService) AssignRoleToUser(ctx context.Context, UserRole dto.UserRoleRequest) error {
	err := rs.repo.AssignRoleToUser(ctx, UserRole)
	if err != nil {
		return err
	}
	return nil
}
func (rs *DefaultRoleService) GetRoles(ctx context.Context) (*[]dto.RolePermissionResponse, error) {
	result, err := rs.repo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	roles := make([]dto.RolePermissionResponse, len(result))
	for i, role := range result {
		roles[i] = dto.RolePermissionResponse{
			RoleName:    role.RoleName,
			RoleID:      role.RoleID,
			Permissions: role.Permissions,
		}
	}
	return &roles, nil
}

func NewRoleService(repo domain.RoleRepository) *DefaultRoleService {
	return &DefaultRoleService{repo: repo}
}
