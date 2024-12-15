package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
	Username string`json:"username" validate:"alphanum,required,min=3,max=20"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginUser struct {
	Email string`json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	ID string `json:"id"`
	Username string`json:"username"`
	Email string `json:"email"`
	Role string	`json:"role"`
} 

// Hashes user password for storage
func (u RegisterUser) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash provided password: %w", err)
	}
	return string(bytes), nil
}

// Compares user provided password with the user password hash stored in the db
func (u User) ValidatePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}