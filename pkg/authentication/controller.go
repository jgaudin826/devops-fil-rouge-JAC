package authentication

import (
	"mediadex/config"
	"mediadex/database/dbmodel"
	"mediadex/pkg/model"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type AuthConfig struct {
	*config.Config
}

func New(configuration *config.Config) *AuthConfig {
	return &AuthConfig{configuration}
}

// @Summary		User login
// @Description	Authenticate a user and return JWT tokens
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		model.LoginRequest	true	"Login credentials"
// @Success		200		{object}	model.TokenResponse
// @Failure		400		{object}	map[string]string
// @Failure		401		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/auth/login [post]
func (config *AuthConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.LoginRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid login payload: " + err.Error()})
		return
	}

	user, err := config.UserRepository.FindByEmail(req.Email)
	if err != nil {
		user, err = config.UserRepository.FindByUsername(req.Username)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "invalid email/username or password"})
			return
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "invalid email/username or password"})
		return
	}

	accessToken, err := GenerateToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate token"})
		return
	}
	refreshToken, err := GenerateRefreshToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate refresh token"})
		return
	}

	tokens := &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
	}

	render.JSON(w, r, tokens)
}

// @Summary		User register
// @Description	Create a new user and return JWT tokens
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		model.UserRequest	true	"Register credentials"
// @Success		200		{object}	model.TokenResponse
// @Failure		400		{object}	map[string]string
// @Failure		409		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/auth/register [post]
func (config *AuthConfig) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid register payload: " + err.Error()})
		return
	}

	_, err := config.UserRepository.FindByEmail(req.Email)
	if err == nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": "email or username already in use"})
		return
	}
	_, err = config.UserRepository.FindByUsername(req.Username)
	if err == nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": "email or username already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to hash password"})
		return
	}
	req.Password = string(hashedPassword)

	userEntry := &dbmodel.User{Email: req.Email, Password: req.Password, Username: req.Username}
	res, err := config.UserRepository.Create(userEntry)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to create user"})
		return
	}
	user := &model.UserResponse{ID: res.ID, Email: res.Email, Username: res.Username}

	accessToken, err := GenerateToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate token"})
		return
	}
	refreshToken, err := GenerateRefreshToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate refresh token"})
		return
	}
	tokens := &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
	}

	render.JSON(w, r, tokens)
}

// @Summary		Refresh token
// @Description	Generate a new JWT token from an existing valid refresh token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		model.TokenRequest	true	"Refresh token"
// @Success		200		{object}	model.TokenResponse
// @Failure		400		{object}	map[string]string
// @Failure		401		{object}	map[string]string
// @Failure		404		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/auth/refresh [post]
func (config *AuthConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.TokenRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid refresh payload: " + err.Error()})
		return
	}

	email, err := ParseToken(config.JWTSecret, req.RefreshToken)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "invalid refresh token"})
		return
	}

	user, err := config.UserRepository.FindByEmail(email)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "user not found"})
		return
	}
	accessToken, err := GenerateToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate token"})
		return
	}
	refreshToken, err := GenerateRefreshToken(config.JWTSecret, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to generate refresh token"})
		return
	}

	tokens := &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
	}

	render.JSON(w, r, tokens)
}
