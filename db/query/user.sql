-- name: CreateUser :one
INSERT INTO users(username , email , password) VALUES ($1 , $2 , $3) RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserRolesPermissions :one
SELECT
    u.id AS user_id,
    u.username AS user_name,
    u.email AS user_email,
    COALESCE(json_agg(DISTINCT r.*) FILTER (WHERE r.id IS NOT NULL), '[]') AS roles,
    COALESCE(json_agg(DISTINCT p.*) FILTER (WHERE p.id IS NOT NULL), '[]') AS permissions
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
         LEFT JOIN role_permissions rp ON r.id = rp.role_id
         LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = $1
GROUP BY u.id;