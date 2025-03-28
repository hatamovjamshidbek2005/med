package sqlc

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v4"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4/zero"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	argsCalled := m.Called(ctx, sql, args)
	return &MockRow{value: argsCalled.Get(0), err: argsCalled.Error(1)}
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsCalled := m.Called(ctx, sql, args)
	return argsCalled.Get(0).(pgx.Rows), argsCalled.Error(1)
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgx.CommandTag, error) {
	argsCalled := m.Called(ctx, sql, args)
	return argsCalled.Get(0).(pgx.CommandTag), argsCalled.Error(1)
}

type MockRow struct {
	value interface{}
	err   error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	switch v := m.value.(type) {
	case int64:
		*dest[0].(*int64) = v
	case string:
		*dest[0].(*string) = v
	}
	return nil
}

type MockRows struct {
	rows    []interface{}
	current int
	err     error
}

func (m *MockRows) Next() bool {
	return m.current < len(m.rows)
}

func (m *MockRows) Close()     {}
func (m *MockRows) Err() error { return m.err }

func TestCountDoctorAppointment(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countDoctorAppointment, zero.StringFrom("1"), zero.StringFrom("%test%")).
			Return(int64(5), nil)

		count, err := queries.CountDoctorAppointment(context.Background(), CountDoctorAppointmentParams{
			DoctorID: zero.StringFrom("1"),
			FullName: zero.StringFrom("%test%"),
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countDoctorAppointment, zero.StringFrom("1"), zero.StringFrom("%test%")).
			Return(nil, errors.New("db error"))

		count, err := queries.CountDoctorAppointment(context.Background(), CountDoctorAppointmentParams{
			DoctorID: zero.StringFrom("1"),
			FullName: zero.StringFrom("%test%"),
		})
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestCountUserAppointments(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countUserAppointments, zero.StringFrom("1"), "%test%").
			Return(int64(3), nil)

		count, err := queries.CountUserAppointments(context.Background(), CountUserAppointmentsParams{
			UserID:    zero.StringFrom("1"),
			FirstName: "%test%",
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestDeleteOneAppointment(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("Exec", mock.Anything, deleteOneAppointment, "123").
			Return(pgx.CommandTag("DELETE 1"), nil)

		err := queries.DeleteOneAppointment(context.Background(), "123")
		assert.NoError(t, err)
	})
}

// TestEditOneAppointment
func TestEditOneAppointment(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		mockDB.On("QueryRow", mock.Anything, editOneAppointment, zero.StringFrom("1"), zero.StringFrom("2"), now, zero.TimeFrom(now), "123").
			Return("123", nil)

		id, err := queries.EditOneAppointment(context.Background(), EditOneAppointmentParams{
			UserID:          zero.StringFrom("1"),
			DoctorID:        zero.StringFrom("2"),
			AppointmentTime: now,
			UpdatedAt:       zero.TimeFrom(now),
			ID:              "123",
		})
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

// TestEditOneStatusAppointment
func TestEditOneStatusAppointment(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		status := NullAppointmentStatus{AppointmentStatus: "confirmed", Valid: true}
		mockDB.On("QueryRow", mock.Anything, editOneStatusAppointment, status, zero.TimeFrom(now), "123").
			Return("123", nil)

		id, err := queries.EditOneStatusAppointment(context.Background(), EditOneStatusAppointmentParams{
			Status:    status,
			UpdatedAt: zero.TimeFrom(now),
			ID:        "123",
		})
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

func TestInsertOneAppointment(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		mockDB.On("QueryRow", mock.Anything, insertOneAppointment, zero.StringFrom("1"), zero.StringFrom("2"), now).
			Return("123", nil)

		id, err := queries.InsertOneAppointment(context.Background(), InsertOneAppointmentParams{
			UserID:          zero.StringFrom("1"),
			DoctorID:        zero.StringFrom("2"),
			AppointmentTime: now,
		})
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

// TestSelectDoctorAppointments
func TestSelectDoctorAppointments(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := &MockRows{
			rows: []interface{}{
				SelectDoctorAppointmentsRow{
					ID:              "123",
					UserInfo:        json.RawMessage(`{"full_name": "John Doe"}`),
					AppointmentTime: now,
					Status:          NullAppointmentStatus{AppointmentStatus: "pending", Valid: true},
					CreatedAt:       zero.TimeFrom(now),
					UpdatedAt:       zero.TimeFrom(now),
				},
			},
		}
		mockDB.On("Query", mock.Anything, selectDoctorAppointments, zero.StringFrom("1"), zero.StringFrom("%test%"), int32(10), int32(0)).
			Return(rows, nil)

		result, err := queries.SelectDoctorAppointments(context.Background(), SelectDoctorAppointmentsParams{
			DoctorID: zero.StringFrom("1"),
			FullName: zero.StringFrom("%test%"),
			Limit:    10,
			Offset:   0,
		})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].ID)
	})
}

func TestSelectUserAppointments(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := &MockRows{
			rows: []interface{}{
				SelectUserAppointmentsRow{
					ID:              "123",
					DoctorInfo:      json.RawMessage(`{"first_name": "Jane"}`),
					AppointmentTime: now,
					Status:          NullAppointmentStatus{AppointmentStatus: "pending", Valid: true},
					CreatedAt:       zero.TimeFrom(now),
					UpdatedAt:       zero.TimeFrom(now),
				},
			},
		}
		mockDB.On("Query", mock.Anything, selectUserAppointments, zero.StringFrom("1"), "%test%", int32(10), int32(0)).
			Return(rows, nil)

		result, err := queries.SelectUserAppointments(context.Background(), SelectUserAppointmentsParams{
			UserID:    zero.StringFrom("1"),
			FirstName: "%test%",
			Limit:     10,
			Offset:    0,
		})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].ID)
	})
}
