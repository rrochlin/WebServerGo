-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
	$1, $2
)
RETURNING *;


-- name: TruncateUsers :exec
TRUNCATE TABLE users CASCADE;


-- name: GetUser :one
SELECT * FROM users
WHERE users.email=$1;
