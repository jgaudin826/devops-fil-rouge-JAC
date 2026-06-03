package model

import (
	"errors"
	"net/http"
)

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (user *UserRequest) Bind(r *http.Request) error {
	if user.Email == "" {
		return errors.New("email must not be null")
	}
	if user.Username == "" {
		return errors.New("username must not be null")
	}
	if user.Password == "" {
		return errors.New("password must not be null")
	}

	return nil
}

type UserResponse struct {
	ID       uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
