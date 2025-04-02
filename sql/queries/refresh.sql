-- name: CreateRToken :one
INSERT INTO refresh_token (token, user_id)
VALUES (
	$1, $2
)
RETURNING *;

-- name: GetRToken :one
SELECT * FROM refresh_token
WHERE token=$1;

-- name: RevokeToken :one
UPDATE refresh_token SET revoked_at=NOW()
WHERE token=$1
RETURNING *;

