package field

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

type FieldConfig struct {
	*config.Config
}

func New(config *config.Config) *FieldConfig {
	return &FieldConfig{config}
}

// PostHandler godoc
// @Summary      Create a new Field
// @Description  Creates a new Field entry in the database
// @Tags         Field
// @Accept       json
// @Produce      json
// @Param        field  body      model.FieldRequest  true  "Field creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.FieldResponse
// @Failure      400    {object}  map[string]string  "Invalid Field POST request payload !"
// @Failure      500    {object}  map[string]string  "Failed to create Field !"
// @Router       /fields [post]
func (config *FieldConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request.
	req := &model.FieldRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid field payload: " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.Field type for the "Create" function.
	field := &dbmodel.Field{
		UserId: req.UserId,
		Name:   req.Name}

	// Request the DB to Create the Field.
	savedField, err := config.FieldRepository.Create(field)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create field: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response.
	res := &model.FieldResponse{
		UserId: savedField.UserId,
		Name:   savedField.Name}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get field by ID
// @Description  Retrieves a specific field from the database by its ID
// @Tags         Field
// @Produce      json
// @Param        id   path      string  true  "field ID"
// @Security     BearerAuth
// @Success      200  {object}  model.FieldResponse
// @Failure      404  {object}  map[string]string  "Field not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific Field !"
// @Router       /fields/{id} [get]
func (config *FieldConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Get the needed informations
	field, err := config.FieldRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find field: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.FieldResponse{
		UserId: field.UserId,
		Name:   field.Name}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all field
// @Description  Retrieve all field
// @Tags         Field
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.FieldResponse
// @Failure      500  {object}  map[string]string
// @Router       /fields [get]
func (config *FieldConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	fields, err := config.FieldRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch fields: " + err.Error()})
		return
	}

	res := make([]model.FieldResponse, len(fields))
	for i, field := range fields {
		res[i] = model.FieldResponse{
			UserId: field.UserId,
			Name:   field.Name,
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetMediaByFieldHandler godoc
// @Summary      Get media by field
// @Description  Retrieves all media associated with a specific field
// @Tags         Field
// @Produce      json
// @Param        id   path      string  true  "field ID"
// @Security     BearerAuth
// @Success      200  {array}   model.MediaResponse
// @Failure      404  {object}  map[string]string  "Field not found"
// @Failure      500  {object}  map[string]string  "Failed to fetch Media for Field !"
// @Router       /fields/{id}/media [get]
func (config *FieldConfig) GetMediaByFieldHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	mediaItems, err := config.FieldRepository.FindMediaByField(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch media for field: " + err.Error()})
		return
	}

	res := make([]model.MediaResponse, len(mediaItems))
	for i, media := range mediaItems {
		res[i] = model.MediaResponse{
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

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// UpdateHandler godoc
// @Summary      Update a field
// @Description  Update an existing field
// @Tags         Field
// @Accept       json
// @Produce      json
// @Param        id     path     string        true  "Field ID"
// @Param        field  body     model.FieldUpdateRequest  true  "Updated field payload"
// @Security     BearerAuth
// @Success      200   {object}  model.FieldResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /fields/{id} [patch]
func (config *FieldConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Get the request
	req := &model.FieldUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid field payload: " + err.Error()})
		return
	}

	// Request the DB to get the Field data
	existing, err := config.FieldRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "field not found: " + err.Error()})
		return
	}

	if req.UserId != nil && *req.UserId > 0 {
		existing.UserId = *req.UserId
	}
	if req.Name != nil && *req.Name != "" {
		existing.Name = *req.Name
	}

	updatedField, err := config.FieldRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update field: " + err.Error()})
		return
	}

	res := model.FieldResponse{
		UserId: updatedField.UserId,
		Name:   updatedField.Name,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a field
// @Description  Deletes a field from the database by its ID
// @Tags         Field
// @Produce      json
// @Param        id   path      string  true  "Field ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Field deleted successfully"
// @Failure      404  {object}  map[string]string  "Field not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete Field !"
// @Router       /fields/{id} [delete]
func (config *FieldConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Delete the informations
	err = config.FieldRepository.Delete(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete field: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Field deleted successfully."})
}
