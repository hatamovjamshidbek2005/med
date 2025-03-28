package sqlc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4/zero"
)

// TestCountDoctors
func TestCountDoctors(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countDoctors).Return(int64(10), nil)

		count, err := queries.CountDoctors(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countDoctors).Return(nil, errors.New("db error"))

		count, err := queries.CountDoctors(context.Background())
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestCountManyDoctors(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countManyDoctors, "%test%").Return(int64(5), nil)

		count, err := queries.CountManyDoctors(context.Background(), "%test%")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestDeleteOneDoctor(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("Exec", mock.Anything, deleteOneDoctor, "123").Return(pgx.CommandTag("DELETE 1"), nil)

		err := queries.DeleteOneDoctor(context.Background(), "123")
		assert.NoError(t, err)
	})
}

func TestEditOneDoctor(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		params := EditOneDoctorParams{
			FirstName:            "John",
			LastName:             "Doe",
			Slug:                 "john-doe",
			Image:                "image.jpg",
			Experience:           10,
			Specialization:       "Cardiology",
			TreatmentProfile:     "Profile",
			ProfessionalActivity: "Activity",
			WorkingHours:         json.RawMessage(`{"mon": "9-5"}`),
			UpdatedAt:            zero.TimeFrom(now),
			ID:                   "123",
		}
		mockDB.On("QueryRow", mock.Anything, editOneDoctor,
			params.FirstName, params.LastName, params.Slug, params.Image, params.Experience,
			params.Specialization, params.TreatmentProfile, params.ProfessionalActivity,
			params.WorkingHours, params.UpdatedAt, params.ID).Return("123", nil)

		id, err := queries.EditOneDoctor(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

func TestInsertOneDoctor(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		params := InsertOneDoctorParams{
			ID:                   "123",
			FirstName:            "John",
			LastName:             "Doe",
			Slug:                 "john-doe",
			Image:                "image.jpg",
			Experience:           10,
			Specialization:       "Cardiology",
			TreatmentProfile:     "Profile",
			ProfessionalActivity: "Activity",
			WorkingHours:         json.RawMessage(`{"mon": "9-5"}`),
			CreatedAt:            zero.TimeFrom(now),
			UpdatedAt:            zero.TimeFrom(now),
		}
		mockDB.On("QueryRow", mock.Anything, insertOneDoctor,
			params.ID, params.FirstName, params.LastName, params.Slug, params.Image,
			params.Experience, params.Specialization, params.TreatmentProfile,
			params.ProfessionalActivity, params.WorkingHours, params.CreatedAt,
			params.UpdatedAt).Return("123", nil)

		id, err := queries.InsertOneDoctor(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

func TestSelectManyDoctors(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := &MockRows{
			rows: []interface{}{
				Doctor{
					ID:                   "123",
					FirstName:            "John",
					LastName:             "Doe",
					Slug:                 "john-doe",
					Image:                "image.jpg",
					Experience:           10,
					Specialization:       "Cardiology",
					TreatmentProfile:     "Profile",
					ProfessionalActivity: "Activity",
					WorkingHours:         json.RawMessage(`{"mon": "9-5"}`),
					CreatedAt:            zero.TimeFrom(now),
					UpdatedAt:            zero.TimeFrom(now),
				},
			},
		}
		mockDB.On("Query", mock.Anything, selectManyDoctors, "%test%", int32(10), int32(0)).
			Return(rows, nil)

		result, err := queries.SelectManyDoctors(context.Background(), SelectManyDoctorsParams{
			FirstName: "%test%",
			Limit:     10,
			Offset:    0,
		})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].ID)
	})
}

func TestSelectOneDoctor(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		expected := SelectOneDoctorRow{
			FirstName:            "John",
			LastName:             "Doe",
			Slug:                 "john-doe",
			Image:                "image.jpg",
			Experience:           10,
			Specialization:       "Cardiology",
			TreatmentProfile:     "Profile",
			ProfessionalActivity: "Activity",
			WorkingHours:         json.RawMessage(`{"mon": "9-5"}`),
			CreatedAt:            zero.TimeFrom(now),
			UpdatedAt:            zero.TimeFrom(now),
		}
		mockDB.On("QueryRow", mock.Anything, selectOneDoctor, "123").Return(expected, nil)

		result, err := queries.SelectOneDoctor(context.Background(), "123")
		assert.NoError(t, err)
		assert.Equal(t, "John", result.FirstName)
		assert.Equal(t, "Doe", result.LastName)
	})
}
