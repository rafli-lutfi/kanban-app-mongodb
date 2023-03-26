package controllers

import (
	"encoding/json"
	"net/http"
	"time"

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
		responeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		responeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	userID, err := h.userService.Register(ctx, user)
	if err != nil {
		responeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responeWithJson(w, http.StatusCreated, "success registered", map[string]interface{}{
		"id": userID,
	})
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var creds models.UserLogin

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		responeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if creds.Email == "" || creds.Password == "" {
		responeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	user, err := h.userService.Login(ctx, creds)
	if err != nil {
		responeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responeWithJson(w, http.StatusOK, "success logged in", map[string]interface{}{
		"fullname": user.Fullname,
	})
}

func responeWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": statusCode,
		"error":  msg,
	})
}

func responeWithJson(w http.ResponseWriter, statusCode int, msg string, payload interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  statusCode,
		"message": msg,
		"data":    payload,
	})
}
