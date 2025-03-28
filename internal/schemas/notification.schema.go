package schemas

import "time"

type Notification struct {
	ID               string    `json:"id"`
	AppointmentID    string    `json:"appointment_id"`
	UserID           string    `json:"user_id"`
	DoctorID         string    `json:"doctor_id"`
	NotificationType string    `json:"notification_type"`
	SentAt           time.Time `json:"sent_at"`
	Message          string    `json:"message"`
}
