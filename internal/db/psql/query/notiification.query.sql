-- name: InsertNotification :one
INSERT INTO public.notifications (
    appointment_id, user_id, doctor_id, scheduled_at, message
) VALUES (
             $1, $2, $3, $4, $5
         ) RETURNING id;

-- name: SelectPendingNotifications :many
SELECT
    n.id, n.appointment_id, n.user_id, n.doctor_id, n.scheduled_at, n.message,
    u.full_name, u.email,
    d.first_name AS doctor_first_name, d.last_name AS doctor_last_name
FROM public.notifications n
         JOIN public.users u ON n.user_id = u.id
         JOIN public.doctors d ON n.doctor_id = d.id
WHERE n.status = 'pending' AND n.scheduled_at <= NOW();

-- name: UpdateNotificationStatus :exec
UPDATE public.notifications
SET status = $1, sent_at = NOW()
WHERE id = $2;