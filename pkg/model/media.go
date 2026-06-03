package model

import (
	"errors"
	"net/http"
	"slices"
	t "time"
)

type MediaRequest struct {
	UserId         uint    `json:"user_id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	MediaType      string  `json:"media_type"`
	ImgURL         *string `json:"img_url"`
	Rating         *int    `json:"rating"`
	Notes          *string `json:"notes"`
	Description    *string `json:"description"`
	Genre          *string `json:"genre"`
	StartDate      *t.Time `json:"start_date"`
	CompletionDate *t.Time `json:"completion_date"`
}

type MediaUpdateRequest struct {
	UserId         *uint   `json:"user_id"`
	Name           *string `json:"name"`
	Status         *string `json:"status"`
	MediaType      *string `json:"media_type"`
	ImgURL         *string `json:"img_url"`
	Rating         *int    `json:"rating"`
	Notes          *string `json:"notes"`
	Description    *string `json:"description"`
	Genre          *string `json:"genre"`
	StartDate      *t.Time `json:"start_date"`
	CompletionDate *t.Time `json:"completion_date"`
}

var validMediaStatuses = []string{"Planned", "In Progress", "Paused", "Completed", "Abandoned", "For Later"}
var validMediaTypes = []string{"Film", "Shows", "Games", "Books"}

func validateMediaFields(status *string, mediaType *string) error {
	if status != nil && !slices.Contains(validMediaStatuses, *status) {
		return errors.New("status is invalid")
	}
	if mediaType != nil && !slices.Contains(validMediaTypes, *mediaType) {
		return errors.New("media_type is invalid")
	}

	return nil
}

func (media *MediaRequest) Bind(r *http.Request) error {
	if media.UserId == 0 {
		return errors.New("user_id must not be null")
	}
	if media.Name == "" {
		return errors.New("name must not be null")
	}
	if err := validateMediaFields(&media.Status, &media.MediaType); err != nil {
		return err
	}

	return nil
}

func (media *MediaUpdateRequest) Validate() error {
	return validateMediaFields(media.Status, media.MediaType)
}

type MediaResponse struct {
	ID             uint    `json:"media_id"`
	UserId         uint    `json:"user_id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	MediaType      string  `json:"media_type"`
	ImgURL         *string `json:"img_url"`
	Rating         *int    `json:"rating"`
	Notes          *string `json:"notes"`
	Description    *string `json:"description"`
	Genre          *string `json:"genre"`
	StartDate      *t.Time `json:"start_date"`
	CompletionDate *t.Time `json:"completion_date"`
}
