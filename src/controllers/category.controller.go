package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTask(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *categoryHandler {
	return &categoryHandler{categoryService}
}

func (h *categoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	id := r.URL.Query().Get("category_id")

	if id == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrQueryParamEmpty.Error())
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, "invalid id")
		return
	}

	category, err := h.categoryService.GetCategoryByID(ctx, categoryID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success get category", category)
}

func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var body models.CategoryRequest

	id := fmt.Sprintf("%s", r.Context().Value("id"))

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if body.Type == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	category := models.Category{
		Id:     primitive.NewObjectID(),
		Type:   body.Type,
		UserId: userID,
	}

	categoryID, err := h.categoryService.StoreCategory(ctx, category)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success created category", map[string]interface{}{
		"id": categoryID,
	})
}

func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {}

func (h *categoryHandler) GetCategoryWithTask(w http.ResponseWriter, r *http.Request) {}