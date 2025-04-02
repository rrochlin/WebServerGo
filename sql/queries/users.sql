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


-- name: UpdateUser :one
UPDATE users SET email=$1, hashed_password=$2
WHERE users.id=$3
RETURNING *;
