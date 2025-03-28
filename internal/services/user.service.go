package services

import (
	"context"
	"gopkg.in/guregu/null.v4/zero"
	"med/internal/db/psql/sqlc"
	"med/internal/schemas"
)

func (userService *Service) UpdateUser(ctx context.Context, userParams *schemas.UpdateUserProfilePayload) (*schemas.IDResponse, error) {
	user, err := userService.storage.EditProfileUser(ctx, sqlc.EditProfileUserParams{
		FullName:    zero.StringFrom(userParams.FullName),
		PhoneNumber: zero.StringFrom(userParams.PhoneNumber),
		ID:          userParams.ID,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{user}, nil
}
func (userService *Service) DeleteUser(ctx context.Context, userParams *schemas.IDRequest) (*schemas.ResponseSuccess, error) {
	err := userService.storage.DeleteOneUser(ctx, userParams.ID)
	if err != nil {
		return nil, err
	}
	return &schemas.ResponseSuccess{Data: "success deleted user", Status: 200}, nil
}

func (userService *Service) GetUser(ctx context.Context, userParams *schemas.IDRequest) (*schemas.UserResponse, error) {

	user, err := userService.storage.SelectOneUser(ctx, userParams.ID)
	if err != nil {
		return nil, err
	}

	return &schemas.UserResponse{
		ID:          userParams.ID,
		FullName:    user.FullName.String,
		PhoneNumber: user.PhoneNumber.String,
		UserName:    user.Username,
		Email:       user.Email.String,
		CreatedAt:   user.CreatedAt.Time,
	}, nil
}

func (userService *Service) GetAllUsers(ctx context.Context, userParams *schemas.GetSearchRequest) (*schemas.ManyUsers, error) {
	tx, err := userService.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	users, err := userService.storage.SelectManyUsers(ctx, zero.StringFrom(userParams.Search))
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return
			}
		}
	}()
	var userList []schemas.UserResponse
	for _, user := range users {

		userList = append(userList, schemas.UserResponse{
			ID:          user.ID,
			FullName:    user.FullName.String,
			PhoneNumber: user.PhoneNumber.String,
			UserName:    user.Username,
			Email:       user.Email.String,
			CreatedAt:   user.CreatedAt.Time,
		})
	}
	tx.Commit(ctx)
	countUsers, err := userService.storage.CountUsers(ctx, zero.StringFrom(userParams.Search))
	if err != nil {
		return nil, err
	}
	return &schemas.ManyUsers{
		Users: userList,
		Count: int(countUsers),
	}, nil
}
