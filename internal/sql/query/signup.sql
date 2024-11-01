-- name: InsertSignup :one
INSERT INTO signup (email, username, password, verification_code, expire)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetSignupCode :one
SELECT verification_code
FROM signup
WHERE id = $1;

-- name: GetSignUpData :one
SELECT *
FROM signup
WHERE id = $1;

-- name: DeleteSignup :exec
DELETE
FROM signup
WHERE id = $1;