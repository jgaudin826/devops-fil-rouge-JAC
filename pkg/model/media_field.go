package model

import (
	"errors"
	"net/http"
)

type MediaFieldRequest struct {
	FieldID uint   `json:"field_id"`
	MediaID uint   `json:"media_id"`
	Value   string `json:"value"`
}

type MediaFieldUpdateRequest struct {
	Value *string `json:"value"`
}

func (mediaField *MediaFieldRequest) Bind(r *http.Request) error {
	if mediaField.FieldID == 0 {
		return errors.New("field_id must not be null")
	}
	if mediaField.MediaID == 0 {
		return errors.New("media_id must not be null")
	}
	if mediaField.Value == "" {
		return errors.New("value must not be null")
	}

	return nil
}

type MediaFieldResponse struct {
	FieldID uint   `json:"field_id"`
	MediaID uint   `json:"media_id"`
	Value   string `json:"value"`
}
