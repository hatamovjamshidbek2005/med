-- name: InsertOneAppointment :one
INSERT INTO public.appointments (user_id, doctor_id, appointment_time)
VALUES ($1, $2, $3)
    RETURNING id;

-- name: EditOneAppointment :one
UPDATE public.appointments
SET user_id = $1,
    doctor_id = $2,
    appointment_time = $3,
    updated_at = $4
WHERE id = $5
    RETURNING id;

-- name: EditOneStatusAppointment :one
UPDATE public.appointments
SET status = $1,
    updated_at = $2
WHERE id = $3
    RETURNING id;

-- name: DeleteOneAppointment :exec
DELETE FROM public.appointments
WHERE id = $1;

-- name: SelectUserAppointments :many
SELECT
    pa.id,
    jsonb_build_object(
            'first_name', pd.first_name,
            'last_name', pd.last_name,
            'image', pd.image,
            'experience', pd.experience,
            'specialization', pd.specialization
    ) AS doctor_info,
    pa.appointment_time,
    pa.status,
    pa.created_at,
    pa.updated_at
FROM public.appointments pa
         LEFT JOIN public.doctors pd ON pa.doctor_id = pd.id
WHERE pa.user_id = $1 AND
           pd.first_name ILIKE  $2
            OR pd.last_name ILIKE  $2
            OR pd.specialization ILIKE  $2
            OR pd.treatment_profile ILIKE  $2
            OR pd.professional_activity ILIKE  $2

    LIMIT $3 OFFSET $4;

-- name: SelectDoctorAppointments :many
SELECT
    pa.id,
    jsonb_build_object(
            'full_name', pu.full_name,
            'phone_number', pu.phone_number
    ) AS user_info,
    pa.appointment_time,
    pa.status,
    pa.created_at,
    pa.updated_at
FROM public.appointments pa
         LEFT JOIN public.users pu ON pa.user_id = pu.id

WHERE pa.doctor_id = $1
  AND
    pu.full_name ILIKE  $2
            OR pu.phone_number ILIKE  $2
            OR pa.status::text ILIKE  $2
 LIMIT $3 OFFSET $4
;

-- name: CountDoctorAppointment :one
SELECT
      COUNT(*)
FROM public.appointments pa
         LEFT JOIN public.users pu ON pa.user_id = pu.id
WHERE pa.doctor_id = $1
  AND
    pu.full_name ILIKE  $2
            OR pu.phone_number ILIKE $2
            OR pa.status::text ILIKE $2 ;


-- name: CountUserAppointments :one
SELECT
   COUNT(*)
FROM public.appointments pa
         LEFT JOIN public.doctors pd ON pa.doctor_id = pd.id
WHERE pa.user_id = $1 AND
    pd.first_name ILIKE $2
            OR pd.last_name ILIKE $2
            OR pd.specialization ILIKE  $2
            OR pd.treatment_profile ILIKE  $2
            OR pd.professional_activity ILIKE  $2 ;



