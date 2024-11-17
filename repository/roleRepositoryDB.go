package repository

import (
	"context"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoleRepositoryDB struct {
	q *db.Queries
}

func (rr *RoleRepositoryDB) CreateRole(ctx context.Context, role domain.Role) (*domain.Role, error) {
	arg := role.Name

	result, err := rr.q.CreateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	role = domain.Role{
		Name: result.Name,
		ID:   int(result.ID),
	}
	return &role, nil
}

func (rr *RoleRepositoryDB) UpdateRole(ctx context.Context, role domain.Role) (*domain.Role, error) {
	arg := db.UpdateRoleParams{
		ID: int32(role.ID),
		Name: pgtype.Text{
			String: role.Name,
			Valid:  role.Name != "",
		},
	}
	result, err := rr.q.UpdateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	role = domain.Role{
		Name: result.Name,
		ID:   int(result.ID),
	}
	return &role, nil
}

func (rr *RoleRepositoryDB) DeleteRole(ctx context.Context, id int) error {
	err := rr.q.DeleteRole(ctx, int32(id))
	if err != nil {
		return err
	}
	return nil
}
func (rr *RoleRepositoryDB) GetRole(ctx context.Context, id int) (*domain.Role, error) {
	result, err := rr.q.GetRole(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	role := domain.Role{
		Name: result.Name,
		ID:   int(result.ID),
	}
	return &role, nil
}
func (rr *RoleRepositoryDB) AssignRoleToUser(ctx context.Context, UserRole dto.UserRoleRequest) error {
	arg := db.AssignRoleToUserParams{
		RoleID: pgtype.Int4{
			Int32: int32(UserRole.RoleID),
			Valid: UserRole.RoleID != 0,
		},
		UserID: pgtype.Int4{
			Int32: int32(UserRole.UserID),
			Valid: UserRole.UserID != 0,
		},
	}
	err := rr.q.AssignRoleToUser(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}

func (rr *RoleRepositoryDB) GetRoles(ctx context.Context) ([]*dto.RolePermissionResponse, error) {
	result, err := rr.q.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	rolesMap := make(map[int]*dto.RolePermissionResponse)
	for _, row := range result {
		roleID := int(row.RoleID)

		// Check if the role already exists in the map
		if _, exists := rolesMap[roleID]; !exists {
			// If not, create a new role entry
			rolesMap[roleID] = &dto.RolePermissionResponse{
				RoleID:      roleID,
				RoleName:    row.RoleName,
				Permissions: []dto.PermissionResponse{},
			}
		}
		if row.PermissionID.Valid {
			rolesMap[roleID].Permissions = append(rolesMap[roleID].Permissions, dto.PermissionResponse{
				ID:   uint(row.PermissionID.Int32),
				Name: row.PermissionName.String,
			})
		}
	}
	roles := make([]*dto.RolePermissionResponse, 0, len(rolesMap))
	for _, role := range rolesMap {
		roles = append(roles, role)
	}

	return roles, nil
}

func NewRoleRepositoryDB(q *db.Queries) *RoleRepositoryDB {
	return &RoleRepositoryDB{q}
}
