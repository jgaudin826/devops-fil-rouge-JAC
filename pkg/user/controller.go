package user

import (
	"encoding/json"
	"mediadex/config"
	"mediadex/database/dbmodel"
	"mediadex/pkg/model"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type UserConfig struct {
	*config.Config
}

func New(config *config.Config) *UserConfig {
	return &UserConfig{config}
}

// PostHandler godoc
// @Summary      Create a new User
// @Description  Creates a new User entry in the database
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.UserRequest  true  "User creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.UserResponse
// @Failure      400    {object}  map[string]string  "Invalid User POST request payload !"
// @Failure      500    {object}  map[string]string  "Failed to create User !"
// @Router       /users [post]
func (config *UserConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request.
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid user payload: " + err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to hash password"})
		return
	}

	// Convert the requested data into dbmodel.User type for the "Create" function.
	user := &dbmodel.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Request the DB to Create the User.
	savedUser, err := config.UserRepository.Create(user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create user: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response.
	res := &model.UserResponse{
		ID:       savedUser.ID,
		Username: savedUser.Username,
		Email:    savedUser.Email,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get User by ID
// @Description  Retrieves a specific User from the database by its ID
// @Tags         User
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  model.UserResponse
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific User !"
// @Router       /users/{id} [get]
func (config *UserConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Get the needed informations
	user, err := config.UserRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "failed to find user: " + err.Error()})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all user
// @Description  Retrieve all user
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.UserResponse
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (config *UserConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	users, err := config.UserRepository.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to fetch users: " + err.Error()})
		return
	}

	var result []model.UserResponse

	for _, user := range users {
		res := model.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}
		result = append(result, res)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}

// UpdateHandler godoc
// @Summary      Update a user
// @Description  Update an existing user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id     path     string        true  "User ID"
// @Param        user  body     model.UserUpdateRequest  true  "Updated user payload"
// @Security     BearerAuth
// @Success      200   {object}  model.UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [patch]
func (config *UserConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Get the request
	req := &model.UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid user payload: " + err.Error()})
		return
	}

	// Request the DB to get the User data
	existing, err := config.UserRepository.FindById(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "user not found: " + err.Error()})
		return
	}

	if req.Username != nil && *req.Username != "" {
		existing.Username = *req.Username
	}
	if req.Email != nil && *req.Email != "" {
		existing.Email = *req.Email
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to hash password"})
			return
		}
		existing.Password = string(hashedPassword)
	}

	updatedUser, err := config.UserRepository.Update(existing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update user: " + err.Error()})
		return
	}

	res := model.UserResponse{
		ID:       updatedUser.ID,
		Email:    updatedUser.Email,
		Username: updatedUser.Username,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a user
// @Description  Deletes a user from the database by its ID
// @Tags         User
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "User deleted successfully"
// @Failure      404  {object}  map[string]string  "User not found !"
// @Failure      500  {object}  map[string]string  "Failed to delete User !"
// @Router       /users/{id} [delete]
func (config *UserConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	// Request the DB to Delete the informations
	err = config.UserRepository.Delete(uint(id))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete user: " + err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User deleted successfully."})
}
