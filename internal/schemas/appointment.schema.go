package schemas

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type UserAppointment struct {
	ID       string `json:"id"`
	UserInfo struct {
		FullName    string `json:"full_name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user_info"`
	AppointmentTime time.Time             `json:"appointment_time"`
	Status          NullAppointmentStatus `json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}
type ManyUserAppointment struct {
	UserAppointments []UserAppointment `json:"user_appointments"`
	Count            int               `json:"count"`
}

type DoctorAppointment struct {
	ID         string `json:"id"`
	DoctorInfo struct {
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Image          string `json:"image"`
		Experience     int    `json:"experience"`
		Specialization string `json:"specialization"`
	} `json:"doctor_info"`
	AppointmentTime time.Time             `json:"appointment_time"`
	Status          NullAppointmentStatus `json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"s"`
}
type ManyDoctorAppointment struct {
	DoctorAppointments []DoctorAppointment `json:"doctor_appointments"`
	Count              int                 `json:"count"`
}

type AppointmentPayload struct {
	ID              string `json:"-"`
	UserID          string `json:"user_id" validate:"required" example:"UUID"`
	DoctorID        string `json:"doctor_id" validate:"required" example:"UUID"`
	AppointmentTime string `json:"appointment_time" validate:"required" example:"YYYY-MM-DD H:M:S"`
}
type AppointmentStatusPayload struct {
	ID     string            `json:"-"`
	Status AppointmentStatus `json:"status" validate:"required"`
}

type Appointment struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	DoctorID        string            `json:"doctor_id"`
	AppointmentTime string            `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCanceled  AppointmentStatus = "canceled"
)

func (e *AppointmentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AppointmentStatus(s)
	case string:
		*e = AppointmentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AppointmentStatus: %T", src)
	}
	return nil
}

type NullAppointmentStatus struct {
	AppointmentStatus AppointmentStatus `json:"appointment_status"`
	Valid             bool              `json:"valid"`
}

func (ns *NullAppointmentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.AppointmentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AppointmentStatus.Scan(value)
}

func (ns NullAppointmentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AppointmentStatus), nil
}
