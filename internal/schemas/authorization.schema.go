package schemas

import "time"

type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type SignUpPayload struct {
	UserName     string `json:"user_name" validate:"required" example:"UNIQUE"`
	Email        string `json:"email" validate:"required" example:"@gmail.com"`
	PasswordHash string `json:"password" validate:"required" example:"**********"`
}
type SignInPayload struct {
	UserName     string `json:"user_name" validate:"required" example:"UNIQUE"`
	PasswordHash string `json:"password" validate:"required" example:"**********"`
}

type ForgetPassPayload struct {
	Email        string `json:"email" validate:"required" example:"@gmail.com"`
	PasswordHash string `json:"password" validate:"required" example:"**********"`
}
type TokenResponse struct {
	AccessToken        string  `json:"access_token"`
	RefreshToken       string  `json:"refresh_token"`
	AccessExpiredTime  float64 `json:"access_expired_time"`
	RefreshExpiresTime float64 `json:"refresh_expires_time"`
	Success            bool    `json:"success"`
}
type TokenPayload struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
