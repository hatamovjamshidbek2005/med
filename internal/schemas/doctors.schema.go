package schemas

import "time"

type Doctor struct {
	ID               string              `json:"id"`
	FirstName        string              `json:"first_name" `
	LastName         string              `json:"last_name"`
	Slug             string              `json:"slug"`
	Image            string              `json:"image"`
	Experience       int                 `json:"experience"`
	Specialization   string              `json:"specialization"`
	TreatmentProfile string              `json:"treatment_profile"`
	WorkingHours     []map[string]string `json:"working_hours"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}
type DoctorPayload struct {
	ID                   string              `json:"-"`
	FirstName            string              `json:"first_name" validate:"required" example:"STRING"`
	LastName             string              `json:"last_name" validate:"required" example:"STRING"`
	Slug                 string              `json:"-"`
	Image                string              `json:"-"`
	Experience           int                 `json:"experience" validate:"required"`
	Specialization       string              `json:"specialization" validate:"required"`
	TreatmentProfile     string              `json:"treatment_profile" validate:"required"`
	ProfessionalActivity string              `json:"professional_activity" validate:"required"`
	WorkingHours         []map[string]string `json:"working_hours" validate:"required"`
}
type ManyDoctors struct {
	Doctors []Doctor `json:"doctors"`
	Count   int      `json:"count"`
}
