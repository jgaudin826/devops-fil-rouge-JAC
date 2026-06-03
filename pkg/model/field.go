package model

import (
	"errors"
	"net/http"
)

type FieldRequest struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
}

type FieldUpdateRequest struct {
	UserId *uint   `json:"user_id"`
	Name   *string `json:"name"`
}

func (field *FieldRequest) Bind(r *http.Request) error {
	if field.UserId == 0 {
		return errors.New("user_id must not be null")
	}
	if field.Name == "" {
		return errors.New("name must not be null")
	}

	return nil
}

type FieldResponse struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
}
