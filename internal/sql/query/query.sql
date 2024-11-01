-- name: CreateProfile :one
INSERT INTO profile (user_id)
VALUES ($1)
RETURNING id;

-- name: GetProfileUserID :one
SELECT user_id
FROM profile
WHERE id = $1;

-- name: GetProfileID :one
SELECT id
FROM profile
WHERE user_id = $1;

-- name: ChangeProfilePic :exec
UPDATE profile
SET profile_pic_address = $1
WHERE id = $2;
