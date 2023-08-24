package contract

import (
	"errors"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (user *User) Bind(r *http.Request) error {
	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

type UserResponse struct {
	ID uint `json:"id"`
}
