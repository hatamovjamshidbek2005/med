package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4/zero"
)

func TestCountUsers(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countUsers, zero.StringFrom("%test%")).Return(int64(5), nil)

		count, err := queries.CountUsers(context.Background(), zero.StringFrom("%test%"))
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, countUsers, zero.StringFrom("%test%")).Return(nil, errors.New("db error"))

		count, err := queries.CountUsers(context.Background(), zero.StringFrom("%test%"))
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// TestDeleteOneUser
func TestDeleteOneUser(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		mockDB.On("Exec", mock.Anything, deleteOneUser, "123").Return(pgx.CommandTag("DELETE 1"), nil)

		err := queries.DeleteOneUser(context.Background(), "123")
		assert.NoError(t, err)
	})
}

func TestEditProfileUser(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		params := EditProfileUserParams{
			FullName:    zero.StringFrom("John Doe"),
			PhoneNumber: zero.StringFrom("+123456789"),
			ID:          "123",
		}
		mockDB.On("QueryRow", mock.Anything, editProfileUser, params.FullName, params.PhoneNumber, params.ID).
			Return("123", nil)

		id, err := queries.EditProfileUser(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})
}

func TestSelectManyUsers(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := &MockRows{
			rows: []interface{}{
				SelectManyUsersRow{
					ID:          "123",
					FullName:    zero.StringFrom("John Doe"),
					PhoneNumber: zero.StringFrom("+123456789"),
					Email:       zero.StringFrom("john@example.com"),
					Username:    "johndoe",
					CreatedAt:   zero.TimeFrom(now),
				},
			},
		}
		mockDB.On("Query", mock.Anything, selectManyUsers, zero.StringFrom("%test%")).Return(rows, nil)

		result, err := queries.SelectManyUsers(context.Background(), zero.StringFrom("%test%"))
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].ID)
		assert.Equal(t, "John Doe", result[0].FullName.String)
	})
}

func TestSelectOneUser(t *testing.T) {
	mockDB := new(MockDB)
	queries := &Queries{db: mockDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		expected := SelectOneUserRow{
			FullName:    zero.StringFrom("John Doe"),
			PhoneNumber: zero.StringFrom("+123456789"),
			Email:       zero.StringFrom("john@example.com"),
			Username:    "johndoe",
			CreatedAt:   zero.TimeFrom(now),
		}
		mockDB.On("QueryRow", mock.Anything, selectOneUser, "123").Return(expected, nil)

		result, err := queries.SelectOneUser(context.Background(), "123")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", result.FullName.String)
		assert.Equal(t, "+123456789", result.PhoneNumber.String)
		assert.Equal(t, "john@example.com", result.Email.String)
		assert.Equal(t, "johndoe", result.Username)
	})

	t.Run("No Rows", func(t *testing.T) {
		mockDB.On("QueryRow", mock.Anything, selectOneUser, "123").Return(nil, sql.ErrNoRows)

		_, err := queries.SelectOneUser(context.Background(), "123")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
