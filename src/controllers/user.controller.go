package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/services"
	"golang.org/x/net/context"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user models.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	userID, err := h.userService.Register(ctx, user)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusCreated, "success registered", map[string]interface{}{
		"id": userID,
	})
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var creds models.UserLogin

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if creds.Email == "" || creds.Password == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	user, err := h.userService.Login(ctx, creds)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenString, err := generateJWT(user.Id.Hex())
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600 * 24,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	models.ResponeWithJson(w, http.StatusOK, "success logged in", map[string]interface{}{
		"fullname": user.Fullname,
	})
}

func generateJWT(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
