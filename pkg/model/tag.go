package model

import (
	"errors"
	"net/http"
)

type TagRequest struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type TagUpdateRequest struct {
	UserId *uint   `json:"user_id"`
	Name   *string `json:"name"`
	Color  *string `json:"color"`
}

func (tag *TagRequest) Bind(r *http.Request) error {
	if tag.UserId == 0 {
		return errors.New("user_id must not be null")
	}
	if tag.Name == "" {
		return errors.New("name must not be null")
	}
	if tag.Color == "" {
		return errors.New("color must not be null")
	}
	return nil
}

type TagResponse struct {
	ID     uint   `json:"tag_id"`
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
