

-- name: RegisterOneUser :one
INSERT INTO users(id, email,username, password_hash, created_at)
VALUES ($1, $2, $3,$4,$5)
    RETURNING id;

-- name: LoginOneUser :one
SELECT
    id
    FROM users WHERE username=$1 AND password_hash=$2;

-- name: SelectOnePassword :one
SELECT
    id,
    password_hash
    FROM users WHERE username=$1;

-- name: ForgotPassword :one
UPDATE public.users SET password_hash=$1 WHERE email=$2 RETURNING id;


