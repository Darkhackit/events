-- name: CreatePermission :one
INSERT INTO permissions(name)
VALUES ($1) RETURNING *;

-- name: GetPermission :one
SELECT * FROM permissions WHERE id = $1 LIMIT 1;

-- name: UpdatePermission :one
UPDATE permissions SET name = COALESCE(sqlc.narg(name),name) WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: GetRolePermissions :many
SELECT p.id , p.name FROM permissions p JOIN role_permissions rp ON rp.permission_id = p.id WHERE rp.role_id = $1;

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions(role_id, permission_id) VALUES ($1 , $2);

-- name: GetPermissions :many
SELECT * FROM permissions;

-- -- name: AssignPermissionsToRoleBatch :exec
-- INSERT INTO role_permissions (role_id, permission_id)
-- VALUES ($1, unnest($2))
-- ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: AssignPermissionsToRoleBatch :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, UNNEST($2::int[]));
-- name: RemoveAllPermissionsFromRole :exec
DELETE FROM role_permissions WHERE role_id = $1;