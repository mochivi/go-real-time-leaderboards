package models

import "time"

type User struct {
	ID string `json:"id"`
	Username string`json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	Role string	`json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUser struct {
	ID string `json:"id" validate:"uuid,required"`
	Username string`json:"username" validate:"alphanum,required,min=3,max=20"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Role string	`json:"role" validate:"oneof=admin moderator client visitor"`
}

type LoginUser struct {
	Email string`json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}