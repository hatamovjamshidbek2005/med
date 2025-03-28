

-- name: InsertOneDoctor :one
        INSERT INTO public.doctors(id, first_name, last_name, slug, image, experience, specialization, treatment_profile, professional_activity, working_hours, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id;

-- name: EditOneDoctor :one
  UPDATE public.doctors
     SET  first_name=$1, last_name=$2, slug=$3, image=$4, experience=$5, specialization=$6, treatment_profile=$7, professional_activity=$8, working_hours=$9, updated_at=$10 WHERE id=$11 RETURNING id;

-- name: DeleteOneDoctor :exec
   DELETE FROM public.doctors WHERE id=$1;

-- name: SelectOneDoctor :one
  SELECT
    first_name, last_name, slug, image, experience, specialization, treatment_profile, professional_activity, working_hours, created_at, updated_at
    FROM public.doctors WHERE id=$1;

-- name: SelectManyDoctors :many
        SELECT
            id, first_name, last_name, slug, image, experience, specialization, treatment_profile, professional_activity, working_hours, created_at, updated_at
        FROM public.doctors
        WHERE
            first_name ILIKE  $1
            OR last_name ILIKE  $1 
            OR specialization ILIKE  $1 
            OR CAST(experience AS TEXT) ILIKE  $1 
            OR treatment_profile ILIKE  $1 
            OR professional_activity ILIKE  $1 
            LIMIT $2 OFFSET $3
        ;

-- name: CountManyDoctors :one
    SELECT COUNT(*)
    FROM public.doctors
    WHERE
        first_name ILIKE  $1 
        OR last_name ILIKE  $1 
        OR specialization ILIKE  $1 
        OR CAST(experience AS TEXT) ILIKE  $1 
        OR treatment_profile ILIKE  $1 
        OR professional_activity ILIKE  $1 ;


-- name: CountDoctors :one
    SELECT COUNT(*) FROM public.doctors;