-- name: CreateRole :one
INSERT INTO roles(name) VALUES ($1) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1 LIMIT 1;

-- name: UpdateRole :one
UPDATE roles SET name = COALESCE(sqlc.narg(name),name) WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: GetUserRoles :many
SELECT r.id , r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles(user_id, role_id) VALUES ($1 , $2);

-- name: GetRoles :many
SELECT
    r.id AS role_id,
    r.name AS role_name,
    p.id AS permission_id,
    p.name AS permission_name
FROM
    roles r
        LEFT JOIN
    role_permissions rp ON r.id = rp.role_id
        LEFT JOIN
    permissions p ON rp.permission_id = p.id
ORDER BY
    r.id, p.id;

