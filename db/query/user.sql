-- name: CreateUser :one
INSERT INTO users (
  id, appid, openid, session_key, unionid, appid_from, unionid_from, openid_from, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, now(), now()
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET session_key=$1, updated_at=now()
WHERE id = $2;

-- name: GetUserOpenDataByID :one
SELECT openid, session_key FROM users
WHERE id=$1;

-- name: GetUserIDByAppidAndOpenid :one
SELECT id FROM users WHERE appid = $1 AND openid = $2;