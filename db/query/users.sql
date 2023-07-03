-- name: CreateUser :one
INSERT INTO users (name, class, email) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
select * from users limit $1 offset $2;

-- name: AddUserCredit :one
UPDATE users SET credit = credit + sqlc.arg(amount) WHERE email = sqlc.arg(email) RETURNING *;