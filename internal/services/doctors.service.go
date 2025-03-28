package services

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"gopkg.in/guregu/null.v4/zero"
	"med/internal/db/psql/sqlc"
	"med/internal/schemas"
	"time"

	"github.com/google/uuid"
)

func (doctorsService *Service) CreateDoctor(ctx context.Context, doctorParams *schemas.DoctorPayload) (*schemas.IDResponse, error) {
	jsonData, err := json.Marshal(doctorParams.WorkingHours)
	if err != nil {
		return nil, err
	}
	doctor, err := doctorsService.storage.InsertOneDoctor(ctx, sqlc.InsertOneDoctorParams{
		ID:                   uuid.NewString(),
		FirstName:            doctorParams.FirstName,
		LastName:             doctorParams.LastName,
		Slug:                 doctorParams.Slug,
		Image:                doctorParams.Image,
		Experience:           int32(doctorParams.Experience),
		Specialization:       doctorParams.Specialization,
		TreatmentProfile:     doctorParams.TreatmentProfile,
		ProfessionalActivity: doctorParams.ProfessionalActivity,
		WorkingHours:         jsonData,
		CreatedAt:            zero.TimeFrom(time.Now()),
		UpdatedAt:            zero.TimeFrom(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{doctor}, nil
}
func (doctorsService *Service) UpdateDoctor(ctx context.Context, doctorParams *schemas.DoctorPayload) (*schemas.IDResponse, error) {
	jsonData, err := json.Marshal(doctorParams.WorkingHours)
	if err != nil {
		return nil, err
	}
	doctor, err := doctorsService.storage.EditOneDoctor(ctx, sqlc.EditOneDoctorParams{
		ID:                   doctorParams.ID,
		FirstName:            doctorParams.FirstName,
		LastName:             doctorParams.LastName,
		Slug:                 doctorParams.Slug,
		Image:                doctorParams.Image,
		Experience:           int32(doctorParams.Experience),
		Specialization:       doctorParams.Specialization,
		TreatmentProfile:     doctorParams.TreatmentProfile,
		ProfessionalActivity: doctorParams.ProfessionalActivity,
		WorkingHours:         jsonData,
		UpdatedAt:            zero.TimeFrom(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{doctor}, nil
}

func (doctorsService *Service) GetDoctor(ctx context.Context, doctorParams *schemas.IDRequest) (*schemas.Doctor, error) {

	doctor, err := doctorsService.storage.SelectOneDoctor(ctx, doctorParams.ID)
	if err != nil {
		return nil, err
	}
	var workingHours []map[string]string
	err = json.Unmarshal(doctor.WorkingHours, &workingHours)
	if err != nil {
		return nil, err
	}

	return &schemas.Doctor{
		ID:               doctorParams.ID,
		FirstName:        doctor.FirstName,
		LastName:         doctor.LastName,
		Slug:             doctor.Slug,
		Image:            doctor.Image,
		Experience:       int(doctor.Experience),
		Specialization:   doctor.Specialization,
		TreatmentProfile: doctor.TreatmentProfile,
		WorkingHours:     workingHours,
		CreatedAt:        doctor.CreatedAt.Time,
		UpdatedAt:        doctor.UpdatedAt.Time,
	}, nil
}

func (doctorsService *Service) DeleteDoctor(ctx context.Context, doctorParams *schemas.IDRequest) (*schemas.ResponseSuccess, error) {

	err := doctorsService.storage.DeleteOneDoctor(ctx, doctorParams.ID)
	if err != nil {
		return nil, err
	}
	return &schemas.ResponseSuccess{
		Status: 200,
		Data:   "success deleted on doctor",
	}, nil
}
func (doctorsService *Service) GetAllDoctor(ctx context.Context, doctorParams *schemas.GetListRequest) (*schemas.ManyDoctors, error) {
	tx, err := doctorsService.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	doctors, err := doctorsService.storage.SelectManyDoctors(ctx, sqlc.SelectManyDoctorsParams{
		FirstName: doctorParams.Search,
		Limit:     int32(doctorParams.Limit),
		Offset:    int32((doctorParams.Limit - 1) * doctorParams.Page),
	})
	if err != nil {
		fmt.Println("----------------------", err.Error())
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return
			}
		}
	}()

	var doctorsList []schemas.Doctor
	for _, doctor := range doctors {
		var workingHours []map[string]string
		err = json.Unmarshal(doctor.WorkingHours, &workingHours)
		if err != nil {
			return nil, err
		}

		doctorsList = append(doctorsList, schemas.Doctor{
			ID:               doctor.ID,
			FirstName:        doctor.FirstName,
			LastName:         doctor.LastName,
			Slug:             doctor.Slug,
			Image:            doctor.Image,
			Experience:       int(doctor.Experience),
			Specialization:   doctor.Specialization,
			TreatmentProfile: doctor.TreatmentProfile,
			WorkingHours:     workingHours,
			CreatedAt:        doctor.CreatedAt.Time,
			UpdatedAt:        doctor.UpdatedAt.Time,
		})
	}
	tx.Commit(ctx)
	manyDoctorsCount, err := doctorsService.storage.CountManyDoctors(ctx, doctorParams.Search)
	if err != nil {
		fmt.Println("----------------------+++++++++++", err.Error())

		return nil, err
	}
	return &schemas.ManyDoctors{
		Doctors: doctorsList,
		Count:   int(manyDoctorsCount),
	}, nil
}
