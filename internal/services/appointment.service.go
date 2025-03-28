package services

import (
	"context"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"med/internal/db/psql/sqlc"
	"med/pkg/serialize"
	"med/pkg/utils"
	"time"

	"med/internal/schemas"
)

func (appointmentService *Service) CreateAppointment(ctx context.Context, appointmentParams *schemas.AppointmentPayload) (*schemas.IDResponse, error) {
	AppointmentTime, err := utils.DateValidate(appointmentParams.AppointmentTime)
	if err != nil {
		return nil, err
	}
	appointment, err := appointmentService.storage.InsertOneAppointment(ctx, sqlc.InsertOneAppointmentParams{
		UserID:          zero.StringFrom(appointmentParams.UserID),
		DoctorID:        zero.StringFrom(appointmentParams.DoctorID),
		AppointmentTime: AppointmentTime,
	})
	if err != nil {
		return nil, err
	}
	scheduledAt := AppointmentTime.Add(-1 * time.Hour)
	message := fmt.Sprintf("Hurmatli bemor, sizning uchrashuvingiz %s da.",
		AppointmentTime.Format("2006-01-02 15:04"))

	_, err = appointmentService.storage.InsertNotification(ctx, sqlc.InsertNotificationParams{
		AppointmentID: zero.StringFrom(appointment),
		UserID:        zero.StringFrom(appointmentParams.UserID),
		DoctorID:      zero.StringFrom(appointmentParams.DoctorID),
		ScheduledAt:   scheduledAt,
		Message:       message,
	})

	return &schemas.IDResponse{ID: appointment}, nil
}

func (appointmentService *Service) UpdateAppointment(ctx context.Context, appointmentParams *schemas.AppointmentPayload) (*schemas.IDResponse, error) {
	AppointmentTime, err := utils.DateValidate(appointmentParams.AppointmentTime)
	if err != nil {
		return nil, err
	}
	appointment, err := appointmentService.storage.EditOneAppointment(ctx, sqlc.EditOneAppointmentParams{
		UserID:          zero.StringFrom(appointmentParams.UserID),
		DoctorID:        zero.StringFrom(appointmentParams.DoctorID),
		AppointmentTime: AppointmentTime,
		UpdatedAt:       zero.TimeFrom(time.Now()),
		ID:              appointmentParams.ID,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{ID: appointment}, nil
}

func (appointmentService *Service) UpdateAppointmentStatus(ctx context.Context, appointmentParams *schemas.Appointment) (*schemas.IDResponse, error) {

	appointment, err := appointmentService.storage.EditOneStatusAppointment(ctx, sqlc.EditOneStatusAppointmentParams{
		Status:    sqlc.NullAppointmentStatus{Valid: true, AppointmentStatus: sqlc.AppointmentStatus(appointmentParams.Status)},
		UpdatedAt: zero.TimeFrom(time.Now()),
		ID:        appointmentParams.ID,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{ID: appointment}, nil
}
func (appointmentService *Service) DeleteAppointment(ctx context.Context, appointmentParams *schemas.IDRequest) (*schemas.ResponseSuccess, error) {

	err := appointmentService.storage.DeleteOneAppointment(ctx, appointmentParams.ID)
	if err != nil {
		return nil, err
	}
	return &schemas.ResponseSuccess{
		Data:   "success deleted on appointment",
		Status: 200,
	}, nil
}

func (appointmentService *Service) GetUserAppointment(ctx context.Context, appointmentParams *schemas.GetListRequestOfUserPayload) (*schemas.ManyUserAppointment, error) {
	tx, err := appointmentService.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	appointments, err := appointmentService.storage.SelectUserAppointments(ctx, sqlc.SelectUserAppointmentsParams{
		UserID:    zero.StringFrom(appointmentParams.ID),
		FirstName: appointmentParams.Search,
		Limit:     int32(appointmentParams.Limit),
		Offset:    int32(appointmentParams.Page),
	})

	if err != nil {
		return nil, err
	}
	userAppointmentsCount, err := appointmentService.storage.CountUserAppointments(ctx, sqlc.CountUserAppointmentsParams{
		UserID:    zero.StringFrom(appointmentParams.ID),
		FirstName: appointmentParams.Search,
	})
	if err != nil {
		return nil, err
	}
	var resp []schemas.UserAppointment
	err = serialize.MarshalUnMarshal(appointments, &resp)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)

	return &schemas.ManyUserAppointment{
		UserAppointments: resp,
		Count:            int(userAppointmentsCount),
	}, nil
}
func (appointmentService *Service) GetDoctorAppointment(ctx context.Context, appointmentParams *schemas.GetListRequestOfUserPayload) (*schemas.ManyDoctorAppointment, error) {
	tx, err := appointmentService.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return
			}
		}
	}()
	appointments, err := appointmentService.storage.SelectDoctorAppointments(ctx, sqlc.SelectDoctorAppointmentsParams{
		DoctorID: zero.StringFrom(appointmentParams.ID),
		FullName: zero.StringFrom(appointmentParams.Search),
		Limit:    int32(appointmentParams.Limit),
		Offset:   int32(appointmentParams.Page),
	})
	if err != nil {
		return nil, err
	}

	var resp []schemas.DoctorAppointment
	err = serialize.MarshalUnMarshal(appointments, &resp)
	if err != nil {
		return nil, err
	}
	appointmentCount, err := appointmentService.storage.CountDoctorAppointment(ctx, sqlc.CountDoctorAppointmentParams{
		DoctorID: zero.StringFrom(appointmentParams.ID),
		FullName: zero.StringFrom(appointmentParams.Search),
	})
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)

	return &schemas.ManyDoctorAppointment{
		DoctorAppointments: resp,
		Count:              int(appointmentCount),
	}, nil
}
