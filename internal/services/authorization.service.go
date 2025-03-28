package services

import (
	"context"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"med/internal/db/psql/sqlc"
	"med/internal/schemas"
	"med/pkg/password_hash"

	"github.com/google/uuid"
)

func (authService *Service) RegisterUser(ctx context.Context, authParams *schemas.SignUpPayload) (*schemas.IDResponse, error) {
	password, err := password_hash.HashingPassword([]byte(authParams.PasswordHash))
	if err != nil {
		return nil, err
	}
	fmt.Println("--------------", password)
	user, err := authService.storage.RegisterOneUser(ctx, sqlc.RegisterOneUserParams{
		ID:           uuid.NewString(),
		Email:        zero.StringFrom(authParams.Email),
		Username:     authParams.UserName,
		PasswordHash: password,
		CreatedAt:    zero.Time{},
	})
	if err != nil {
		return nil, err
	}

	return &schemas.IDResponse{user}, nil
}
func (authService *Service) LoginUser(ctx context.Context, authParams *schemas.SignInPayload) (*schemas.IDResponse, error) {
	password, err := authService.storage.SelectOnePassword(ctx, authParams.UserName)
	if err != nil {
		return nil, err
	}
	err = password_hash.CompareHashedPassword([]byte(password.PasswordHash), []byte(authParams.PasswordHash))
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{ID: password.ID}, nil
}

func (authService *Service) UpdatePassword(ctx context.Context, authParams *schemas.ForgetPassPayload) (*schemas.IDResponse, error) {

	password, err := authService.storage.ForgotPassword(ctx, sqlc.ForgotPasswordParams{
		PasswordHash: authParams.PasswordHash,
		Email:        zero.StringFrom(authParams.Email),
	})
	if err != nil {
		return nil, err
	}
	return &schemas.IDResponse{password}, nil

}
