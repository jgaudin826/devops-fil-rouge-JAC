package model

import (
	"errors"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (login *LoginRequest) Bind(r *http.Request) error {
	if login.Password == "" {
		return errors.New("password must not be null")
	}
	if login.Email == "" && login.Username == "" {
		return errors.New("email or username must not be null")
	}

	return nil
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (token *TokenRequest) Bind(r *http.Request) error {
	if token.RefreshToken == "" {
		return errors.New("refresh token must not be null")
	}
	return nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
