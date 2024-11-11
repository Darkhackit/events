-- name: CreateUser :one
INSERT INTO users(username , email , password) VALUES ($1 , $2 , $3) RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1;