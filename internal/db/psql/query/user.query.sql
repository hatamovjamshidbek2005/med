

-- name: EditProfileUser :one
UPDATE public.users SET full_name=$1,phone_number=$2  WHERE id=$3 RETURNING id;

-- name: DeleteOneUser :exec
    DELETE FROM  public.users WHERE id=$1;


-- name: SelectOneUser :one
SELECT
    full_name,
    phone_number,
    email,
    username,
    created_at
    FROM public.users WHERE id=$1;


-- name: SelectManyUsers :many
    SELECT
        id, full_name, phone_number, email,username, created_at
        FROM public.users
    WHERE
        full_name ILIKE  $1 
         OR phone_number ILIKE  $1 
         OR email ILIKE  $1 
    ORDER BY id ASC;


-- name: CountUsers :one
SELECT
    COUNT(*)
FROM public.users
WHERE
          full_name ILIKE  $1
         OR phone_number ILIKE  $1 
         OR email ILIKE  $1 ;
