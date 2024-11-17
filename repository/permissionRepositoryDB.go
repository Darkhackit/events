package repository

import (
	"context"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionRepositoryDB struct {
	q *pgxpool.Pool
}

func (pr *PermissionRepositoryDB) CreatePermission(ctx context.Context, permission domain.Permission) (*domain.Permission, error) {
	arg := permission.Name
	q := db.New(pr.q)

	result, err := q.CreatePermission(ctx, arg)
	if err != nil {
		return nil, err
	}
	dPermission := domain.Permission{
		Name: result.Name,
		ID:   int(result.ID),
	}
	return &dPermission, nil
}

func (pr *PermissionRepositoryDB) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	q := db.New(pr.q)
	result, err := q.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make([]domain.Permission, len(result))
	for i, permission := range result {
		permissions[i] = domain.Permission{
			Name: permission.Name,
			ID:   int(permission.ID),
		}
	}
	return permissions, nil
}

func (pr *PermissionRepositoryDB) AssignRoles(ctx context.Context, permissions dto.AssignPermissionRequest) error {
	tx, err := pr.q.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	q := db.New(tx)
	roleID := int32(permissions.RoleID)
	permissionIDs := make([]int32, len(permissions.PermissionID))
	for i, permission := range permissions.PermissionID {
		permissionIDs[i] = int32(permission.ID)
	}

	err = q.RemoveAllPermissionsFromRole(ctx, pgtype.Int4{Int32: roleID, Valid: roleID != 0})
	if err != nil {
		return err
	}

	args := db.AssignPermissionsToRoleBatchParams{
		RoleID:  pgtype.Int4{Int32: roleID, Valid: roleID != 0},
		Column2: permissionIDs,
	}
	// Execute batch insert query
	err = q.AssignPermissionsToRoleBatch(ctx, args)
	if err != nil {
		return err
	}
	return nil
}

func NewPermissionRepositoryDB(pool *pgxpool.Pool) *PermissionRepositoryDB {
	return &PermissionRepositoryDB{pool}
}
