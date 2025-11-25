-- name: CreateSession :one
INSERT INTO Session (session_token, id_usuario, username, expires)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM Session
WHERE session_token = $1 AND expires > NOW();

-- name: DeleteSession :exec
DELETE FROM Session
WHERE session_token = $1;