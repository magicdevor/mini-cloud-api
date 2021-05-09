-- name: CreateProfile :exec
INSERT INTO users_profile (
  id, user_id, nickname, avatar_url, gender, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, now(), now()
);

-- name: UpdateProfile :exec
UPDATE users_profile SET nickname=$1, avatar_url=$2, gender=$3, updated_at = now()
WHERE "user_id" = $4;

-- name: GetProfile :one
SELECT nickname, avatar_url FROM users_profile
WHERE "user_id" = $1;

-- name: GetProfileIDByUserId :one
SELECT id FROM users_profile WHERE user_id = $1;