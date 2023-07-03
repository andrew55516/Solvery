-- name: CreateEntry :one
    INSERT INTO entries (user_email, amount, comment) VALUES ($1, $2, $3) RETURNING *;

-- name: ListUserEntries :many
    SELECT * FROM entries WHERE user_email = $1 LIMIT $2 OFFSET $3;

-- name: ListAllEntries :many
    SELECT * FROM entries LIMIT $1 OFFSET $2;