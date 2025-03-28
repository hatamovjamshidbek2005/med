package schemas

import "time"

type UserResponse struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}
type ManyUsers struct {
	Users []UserResponse `json:"users"`
	Count int            `json:"count"`
}

type UpdateUserProfilePayload struct {
	ID          string `json:"-"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}
