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

type TaskHandler interface {
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	StoreTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *taskHandler {
	return &taskHandler{taskService}
}

func (h *taskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrQueryParamEmpty.Error())
		return
	}

	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrInvalidID.Error())
		return
	}

	task, err := h.taskService.GetTaskByID(ctx, taskObjectID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success get task", task)
}

func (h *taskHandler) StoreTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var body models.RequestTask

	fmt.Println("hitted")

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if body.Title == "" || body.CategoryID == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrEmptyDataBody.Error())
		return
	}

	userID := fmt.Sprintf("%s", r.Context().Value("id"))

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, models.ErrInvalidID.Error())
		return
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(body.CategoryID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, models.ErrInvalidID.Error())
		return
	}

	fmt.Println(categoryObjectID)

	task := models.Task{
		Id:          primitive.NewObjectID(),
		Title:       body.Title,
		Description: body.Description,
		CategoryId:  categoryObjectID,
		UserId:      userObjectID,
	}

	taskID, err := h.taskService.StoreTask(ctx, task)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("done")

	models.ResponeWithJson(w, http.StatusCreated, "success created task", map[string]interface{}{
		"id": taskID,
	})
}

func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	err = h.taskService.UpdateTask(ctx, task)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success update task", nil)
}

func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrQueryParamEmpty.Error())
		return
	}

	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrInvalidID.Error())
		return
	}

	err = h.taskService.DeleteTask(ctx, taskObjectID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success delete task", nil)
}

func (h *taskHandler) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	task := models.TaskCategoryRequest{}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrFailedDecodeBody.Error())
		return
	}

	if task.CategoryID == primitive.NilObjectID || task.ID == primitive.NilObjectID {
		models.ResponeWithError(w, http.StatusBadRequest, models.ErrQueryParamEmpty.Error())
		return
	}

	err = h.taskService.UpdateTaskCategory(ctx, task.CategoryID, task.ID)
	if err != nil {
		models.ResponeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponeWithJson(w, http.StatusOK, "success update", nil)
}
