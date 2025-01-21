package user

import "context"

type User struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	PhoneNumber    string `json:"phone_number"`
	DOB            string `json:"dob"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
	Version        string `json:"-"`
	CreatedAt      string `json:"-"`
	UpdatedAt      string `json:"-"`
}

type UserStore interface {
	CreateUser(context.Context, *User) error
	GetUserByID(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	UpdateUser(context.Context, *User) error
	ChangePassword(context.Context, *User) error
	DeleteUser(context.Context, string) error
}

type RegisterUserPayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=3,max=100"`
	PhoneNumber string `json:"phone_number"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type UserResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateUserPayload struct {
	FirstName      *string `json:"first_name" validate:"omitempty"`
	LastName       *string `json:"last_name" validate:"omitempty"`
	PhoneNumber    *string `json:"phone_number" validate:"omitempty"`
	DOB            *string `json:"dob" validate:"required"`
	Gender         *string `json:"gender" validate:"required,oneof=Male Female"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
