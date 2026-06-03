package tag

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

type TagConfig struct {
	*config.Config
}

func New(config *config.Config) *TagConfig {
	return &TagConfig{config}
}

// PostHandler godoc
// @Summary      Create a new Tag
// @Description  Creates a new Tag entry in the database
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Param        tag   body      model.TagRequest  true  "Tag creation payload"
// @Security     BearerAuth
// @Success      200   {object}  model.TagResponse
// @Failure      400   {object}  map[string]string  "Invalid Tag POST request payload !"
// @Failure      500   {object}  map[string]string  "Failed to create Tag !"
// @Router       /tags [post]
func (config *TagConfig) PostHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.TagRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid tag payload: " + err.Error()})
		return
	}

	tag := &dbmodel.Tag{
		UserId: req.UserId,
		Name:   req.Name,
		Color:  req.Color,
	}

	savedTag, err := config.TagRepository.Create(tag)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create tag: " + err.Error()})
		return
	}

	res := &model.TagResponse{
		ID:     savedTag.ID,
		UserId: savedTag.UserId,
		Name:   savedTag.Name,
		Color:  savedTag.Color,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get tag by ID
// @Description  Retrieves a specific tag from the database by its ID
// @Tags         Tag
// @Produce      json
// @Param        id   path      string  true  "tag ID"
// @Security     BearerAuth
// @Success      200  {object}  model.TagResponse
// @Failure      404  {object}  map[string]string  "Tag not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific Tag !"
// @Router       /tags/{id} [get]
func (config *TagConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	tag, err := config.TagRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find tag: " + err.Error()})
		return
	}

	res := &model.TagResponse{
		ID:     tag.ID,
		UserId: tag.UserId,
		Name:   tag.Name,
		Color:  tag.Color,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all tags
// @Description  Retrieve all tags
// @Tags         Tag
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.TagResponse
// @Failure      500  {object}  map[string]string
// @Router       /tags [get]
func (config *TagConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := config.TagRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch tags: " + err.Error()})
		return
	}

	res := make([]model.TagResponse, len(tags))
	for i, tag := range tags {
		res[i] = model.TagResponse{
			ID:     tag.ID,
			UserId: tag.UserId,
			Name:   tag.Name,
			Color:  tag.Color,
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetMediaByTagHandler godoc
// @Summary      Get media by tag
// @Description  Retrieves all media associated with a specific tag
// @Tags         Tag
// @Produce      json
// @Param        id   path      string  true  "tag ID"
// @Security     BearerAuth
// @Success      200  {array}   model.MediaResponse
// @Failure      404  {object}  map[string]string  "Tag not found"
// @Failure      500  {object}  map[string]string  "Failed to fetch Media for Tag !"
// @Router       /tags/{id}/media [get]
func (config *TagConfig) GetMediaByTagHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	mediaItems, err := config.TagRepository.FindMediaByTag(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch media for tag: " + err.Error()})
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
// @Summary      Update a tag
// @Description  Update an existing tag
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Param        id     path     string      true  "Tag ID"
// @Param        tag    body     model.TagUpdateRequest  true  "Updated tag payload"
// @Security     BearerAuth
// @Success      200   {object}  model.TagResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tags/{id} [patch]
func (config *TagConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	req := &model.TagUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid tag payload: " + err.Error()})
		return
	}

	existing, err := config.TagRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "tag not found: " + err.Error()})
		return
	}

	if req.UserId != nil && *req.UserId > 0 {
		existing.UserId = *req.UserId
	}
	if req.Name != nil && *req.Name != "" {
		existing.Name = *req.Name
	}
	if req.Color != nil && *req.Color != "" {
		existing.Color = *req.Color
	}

	updatedTag, err := config.TagRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update tag: " + err.Error()})
		return
	}

	res := model.TagResponse{
		ID:     updatedTag.ID,
		UserId: updatedTag.UserId,
		Name:   updatedTag.Name,
		Color:  updatedTag.Color,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a tag
// @Description  Deletes a tag from the database by its ID
// @Tags         Tag
// @Produce      json
// @Param        id   path      string  true  "Tag ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Tag deleted successfully"
// @Failure      404  {object}  map[string]string  "Tag not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete Tag !"
// @Router       /tags/{id} [delete]
func (config *TagConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	err = config.TagRepository.Delete(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete tag: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Tag deleted successfully."})
}
