package media

import (
	"encoding/json"
	"mediadex/config"
	"mediadex/database/dbmodel"
	"mediadex/pkg/model"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type MediaConfig struct {
	*config.Config
}

func New(config *config.Config) *MediaConfig {
	return &MediaConfig{config}
}

// PostHandler godoc
// @Summary      Create a new Media
// @Description  Creates a new Media entry in the database
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        media  body      model.MediaRequest  true  "Media creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.MediaResponse
// @Failure      400    {object}  map[string]string  "Invalid Media POST request payload !"
// @Failure      500    {object}  map[string]string  "Failed to create Media !"
// @Router       /media [post]
func (config *MediaConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request.
	req := &model.MediaRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid media payload: " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.Media type for the "Create" function.
	media := &dbmodel.Media{
		UserId:         req.UserId,
		Name:           req.Name,
		Status:         req.Status,
		MediaType:      req.MediaType,
		ImgURL:         req.ImgURL,
		Rating:         req.Rating,
		Notes:          req.Notes,
		Description:    req.Description,
		Genre:          req.Genre,
		StartDate:      req.StartDate,
		CompletionDate: req.CompletionDate,
	}

	// Request the DB to Create the Media.
	savedMedia, err := config.MediaRepository.Create(media)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create media: " + err.Error()})
		return
	}

	res := newMediaResponse(savedMedia)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get media by ID
// @Description  Retrieves a specific media from the database by its ID
// @Tags         Media
// @Produce      json
// @Param        id   path      string  true  "media ID"
// @Security     BearerAuth
// @Success      200  {object}  model.MediaResponse
// @Failure      404  {object}  map[string]string  "Media not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific Media !"
// @Router       /media/{id} [get]
func (config *MediaConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Get the needed informations
	media, err := config.MediaRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find media: " + err.Error()})
		return
	}

	res := newMediaResponse(media)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all media
// @Description  Retrieve all media
// @Tags         Media
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.MediaResponse
// @Failure      500  {object}  map[string]string
// @Router       /media [get]
func (config *MediaConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	medias, err := config.MediaRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch media: " + err.Error()})
		return
	}

	res := make([]model.MediaResponse, len(medias))
	for i, media := range medias {
		res[i] = *newMediaResponse(media)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetTagsByMediaHandler godoc
// @Summary      Get tags by Media ID
// @Description  Retrieves tags associated with a specific media
// @Tags         Media
// @Produce      json
// @Param        id   path      string  true  "media ID"
// @Security     BearerAuth
// @Success      200  {array}   model.TagResponse
// @Failure      404  {object}  map[string]string  "Media not found"
// @Failure      500  {object}  map[string]string  "Failed to fetch Tags for Media !"
// @Router       /media/{id}/tags [get]
func (config *MediaConfig) GetTagsByMediaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	tags, err := config.MediaRepository.FindTagsByMedia(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch tags for media: " + err.Error()})
		return
	}

	res := make([]model.TagResponse, len(tags))
	for i, t := range tags {
		res[i] = model.TagResponse{
			ID:     t.ID,
			UserId: t.UserId,
			Name:   t.Name,
			Color:  t.Color,
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetFieldsByMediaHandler godoc
// @Summary      Get fields by Media ID
// @Description  Retrieves fields associated with a specific media
// @Tags         Media
// @Produce      json
// @Param        id   path      string  true  "media ID"
// @Security     BearerAuth
// @Success      200  {array}   model.FieldResponse
// @Failure      404  {object}  map[string]string  "Media not found"
// @Failure      500  {object}  map[string]string  "Failed to fetch Fields for Media !"
// @Router       /media/{id}/fields [get]
func (config *MediaConfig) GetFieldsByMediaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	fields, err := config.MediaRepository.FindFieldsByMedia(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch fields for media: " + err.Error()})
		return
	}

	res := make([]model.FieldResponse, len(fields))
	for i, f := range fields {
		res[i] = model.FieldResponse{
			UserId: f.UserId,
			Name:   f.Name,
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// UpdateHandler godoc
// @Summary      Update a media
// @Description  Update an existing media
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        id     path     string        true  "Media ID"
// @Param        media  body     model.MediaUpdateRequest  true  "Updated media payload"
// @Security     BearerAuth
// @Success      200   {object}  model.MediaResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /media/{id} [patch]
func (config *MediaConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Get the request
	req := &model.MediaUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid media payload: " + err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Request the DB to get the Media data
	existing, err := config.MediaRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "media not found: " + err.Error()})
		return
	}

	if req.UserId != nil && *req.UserId > 0 {
		existing.UserId = *req.UserId
	}
	if req.Name != nil && *req.Name != "" {
		existing.Name = *req.Name
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.MediaType != nil {
		existing.MediaType = *req.MediaType
	}
	if req.ImgURL != nil {
		existing.ImgURL = req.ImgURL
	}
	if req.Rating != nil {
		existing.Rating = req.Rating
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Genre != nil {
		existing.Genre = req.Genre
	}
	if req.StartDate != nil {
		existing.StartDate = req.StartDate
	}
	if req.CompletionDate != nil {
		existing.CompletionDate = req.CompletionDate
	}

	updatedMedia, err := config.MediaRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update media: " + err.Error()})
		return
	}

	res := newMediaResponse(updatedMedia)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a media
// @Description  Deletes a media from the database by its ID
// @Tags         Media
// @Produce      json
// @Param        id   path      string  true  "Media ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Media deleted successfully"
// @Failure      404  {object}  map[string]string  "Media not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete Media !"
// @Router       /media/{id} [delete]
func (config *MediaConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Delete the informations
	err = config.MediaRepository.Delete(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete media: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Media deleted successfully."})
}

func newMediaResponse(media *dbmodel.Media) *model.MediaResponse {
	return &model.MediaResponse{
		ID:             media.ID,
		UserId:         media.UserId,
		Name:           media.Name,
		Status:         media.Status,
		MediaType:      media.MediaType,
		ImgURL:         media.ImgURL,
		Rating:         media.Rating,
		Notes:          media.Notes,
		Description:    media.Description,
		Genre:          media.Genre,
		StartDate:      media.StartDate,
		CompletionDate: media.CompletionDate,
	}
}
