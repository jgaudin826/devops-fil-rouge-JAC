package mediaField

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

type MediaFieldConfig struct {
	*config.Config
}

func New(config *config.Config) *MediaFieldConfig {
	return &MediaFieldConfig{config}
}

// PostHandler godoc
// @Summary      Create a new MediaField
// @Description  Creates a new MediaField entry in the database
// @Tags         MediaField
// @Accept       json
// @Produce      json
// @Param        mediaField  body      model.MediaFieldRequest  true  "MediaField creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.MediaFieldResponse
// @Failure      400    {object}  map[string]string  "Invalid MediaField POST request payload !"
// @Failure      500    {object}  map[string]string  "Failed to create MediaField !"
// @Router       /mediaFields [post]
func (config *MediaFieldConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request.
	req := &model.MediaFieldRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid media field payload: " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.MediaField type for the "Create" function.
	mediaField := &dbmodel.MediaField{
		FieldID: req.FieldID,
		MediaID: req.MediaID,
		Value:   req.Value}

	// Request the DB to Create the MediaField.
	savedMediaField, err := config.MediaFieldRepository.Create(mediaField)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create media field: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response.
	res := &model.MediaFieldResponse{
		FieldID: savedMediaField.FieldID,
		MediaID: savedMediaField.MediaID,
		Value:   savedMediaField.Value}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get mediaField by ID
// @Description  Retrieves a specific mediaField from the database by its ID
// @Tags         MediaField
// @Produce      json
// @Param        fieldID   path      string  true  "Field ID"
// @Param        mediaID   path      string  true  "Media ID"
// @Security     BearerAuth
// @Success      200  {object}  model.MediaFieldResponse
// @Failure      404  {object}  map[string]string  "MediaField not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific MediaField !"
// @Router       /mediaFields/{fieldID}/{mediaID} [get]
func (config *MediaFieldConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	fieldID, err := strconv.Atoi(chi.URLParam(r, "fieldID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid fieldID"})
		return
	}
	mediaID, err := strconv.Atoi(chi.URLParam(r, "mediaID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid mediaID"})
		return
	}

	// Request the DB to Get the needed informations
	mediaField, err := config.MediaFieldRepository.FindById(uint(fieldID), uint(mediaID))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find media field: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.MediaFieldResponse{
		FieldID: mediaField.FieldID,
		MediaID: mediaField.MediaID,
		Value:   mediaField.Value}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all mediaField
// @Description  Retrieve all mediaField
// @Tags         MediaField
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.MediaFieldResponse
// @Failure      500  {object}  map[string]string
// @Router       /mediaFields [get]
func (config *MediaFieldConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	mediaFields, err := config.MediaFieldRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch media fields: " + err.Error()})
		return
	}

	var result []model.MediaFieldResponse

	for _, mediaField := range mediaFields {
		res := model.MediaFieldResponse{
			FieldID: mediaField.FieldID,
			MediaID: mediaField.MediaID,
			Value:   mediaField.Value,
		}
		result = append(result, res)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}

// UpdateHandler godoc
// @Summary      Update a mediaField
// @Description  Update an existing mediaField
// @Tags         MediaField
// @Accept       json
// @Produce      json
// @Param        fieldID     path     string        true  "Field ID"
// @Param        mediaID     path     string        true  "Media ID"
// @Param        mediaField  body     model.MediaFieldUpdateRequest  true  "Updated mediaField payload"
// @Security     BearerAuth
// @Success      200   {object}  model.MediaFieldResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /mediaFields/{fieldID}/{mediaID} [patch]
func (config *MediaFieldConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	fieldID, err := strconv.Atoi(chi.URLParam(r, "fieldID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid fieldID"})
		return
	}
	mediaID, err := strconv.Atoi(chi.URLParam(r, "mediaID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid mediaID"})
		return
	}

	// Get the request
	req := &model.MediaFieldUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid media field payload: " + err.Error()})
		return
	}

	// Request the DB to get the MediaField data
	existing, err := config.MediaFieldRepository.FindById(uint(fieldID), uint(mediaID))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "media field not found: " + err.Error()})
		return
	}

	if req.Value != nil && *req.Value != "" {
		existing.Value = *req.Value
	}

	updatedMediaField, err := config.MediaFieldRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update media field: " + err.Error()})
		return
	}

	res := model.MediaFieldResponse{
		FieldID: updatedMediaField.FieldID,
		MediaID: updatedMediaField.MediaID,
		Value:   updatedMediaField.Value,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a mediaField
// @Description  Deletes a mediaField from the database by its ID
// @Tags         MediaField
// @Produce      json
// @Param        fieldID   path      string  true  "Field ID"
// @Param        mediaID   path      string  true  "Media ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "MediaField deleted successfully"
// @Failure      404  {object}  map[string]string  "MediaField not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete MediaField !"
// @Router       /mediaFields/{fieldID}/{mediaID} [delete]
func (config *MediaFieldConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	fieldID, err := strconv.Atoi(chi.URLParam(r, "fieldID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid fieldID"})
		return
	}
	mediaID, err := strconv.Atoi(chi.URLParam(r, "mediaID"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid mediaID"})
		return
	}

	// Request the DB to Delete the informations
	err = config.MediaFieldRepository.Delete(uint(fieldID), uint(mediaID))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete media field: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "MediaField deleted successfully."})
}
