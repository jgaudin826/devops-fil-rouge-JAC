package collection

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

type CollectionConfig struct {
	*config.Config
}

func New(config *config.Config) *CollectionConfig {
	return &CollectionConfig{config}
}

// PostHandler godoc
// @Summary      Create a new Collection
// @Description  Creates a new Collection entry in the database
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        collection  body      model.CollectionRequest  true  "Collection creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.CollectionResponse
// @Failure      400    {object}  map[string]string  "Invalid Collection POST request payload !"
// @Failure      500    {object}  map[string]string  "Failed to create Collection !"
// @Router       /collections [post]
func (config *CollectionConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request.
	req := &model.CollectionRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid collection payload: " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.Collection type for the "Create" function.
	collection := &dbmodel.Collection{
		UserId:  req.UserId,
		Name:    req.Name,
		Filters: req.Filters,
	}

	// Request the DB to Create the Collection.
	savedCollection, err := config.CollectionRepository.Create(collection)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create collection: " + err.Error()})
		return
	}

	res := &model.CollectionResponse{
		UserId:  savedCollection.UserId,
		Name:    savedCollection.Name,
		Filters: savedCollection.Filters,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get collection by ID
// @Description  Retrieves a specific collection from the database by its ID
// @Tags         Collections
// @Produce      json
// @Param        id   path      string  true  "collection ID"
// @Security     BearerAuth
// @Success      200  {object}  model.CollectionResponse
// @Failure      404  {object}  map[string]string  "Collection not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific Collection !"
// @Router       /collections/{id} [get]
func (config *CollectionConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Get the needed informations
	collection, err := config.CollectionRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find collection: " + err.Error()})
		return
	}

	res := &model.CollectionResponse{
		UserId:  collection.UserId,
		Name:    collection.Name,
		Filters: collection.Filters,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all collection
// @Description  Retrieve all collection
// @Tags         Collections
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.CollectionResponse
// @Failure      500  {object}  map[string]string
// @Router       /collections [get]
func (config *CollectionConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	collections, err := config.CollectionRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch collections: " + err.Error()})
		return
	}

	var result []model.CollectionResponse

	for _, collection := range collections {
		result = append(result, model.CollectionResponse{
			UserId:  collection.UserId,
			Name:    collection.Name,
			Filters: collection.Filters,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}

// UpdateHandler godoc
// @Summary      Update a collection
// @Description  Update an existing collection
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id     path     string        true  "Collection ID"
// @Param        collection  body     model.CollectionUpdateRequest  true  "Updated collection payload"
// @Security     BearerAuth
// @Success      200   {object}  model.CollectionResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /collections/{id} [patch]
func (config *CollectionConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Get the request
	req := &model.CollectionUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid collection payload: " + err.Error()})
		return
	}

	// Request the DB to get the Collection data
	existing, err := config.CollectionRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "collection not found: " + err.Error()})
		return
	}

	if req.UserId != nil && *req.UserId > 0 {
		existing.UserId = *req.UserId
	}
	if req.Name != nil && *req.Name != "" {
		existing.Name = *req.Name
	}
	if req.Filters != nil {
		existing.Filters = *req.Filters
	}

	updatedCollection, err := config.CollectionRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update collection: " + err.Error()})
		return
	}

	res := &model.CollectionResponse{
		UserId:  updatedCollection.UserId,
		Name:    updatedCollection.Name,
		Filters: updatedCollection.Filters,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a collection
// @Description  Deletes a collection from the database by its ID
// @Tags         Collections
// @Produce      json
// @Param        id   path      string  true  "Collection ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Collection deleted successfully"
// @Failure      404  {object}  map[string]string  "Collection not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete Collection !"
// @Router       /collections/{id} [delete]
func (config *CollectionConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Delete the informations
	err = config.CollectionRepository.Delete(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete collection: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Collection deleted successfully."})
}
