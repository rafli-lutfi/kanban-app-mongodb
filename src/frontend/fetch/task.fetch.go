package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskFetch interface {
	GetTaskByID(taskID, userID string) (task, error)
	AddTask(categoryID, title, description, userID string) (int, error)
	UpdateTask(taskID, title, description, userID string) (int, error)
	DeleteTask(taskID, userID string) (int, error)
	UpdateTaskCategory(taskID, categoryID, userID string) (int, error)
}

type taskFetch struct {
}

func NewTaskFetch() *taskFetch {
	return &taskFetch{}
}

func (t *taskFetch) GetTaskByID(taskID, userID string) (task, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return task{}, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/tasks/get?task_id="+taskID), nil)
	if err != nil {
		return task{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return task{}, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return task{}, err
	}

	result := map[string]any{}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return task{}, err
	}

	if result["data"] == nil {
		return task{}, err
	}

	taskJson, _ := json.Marshal(result["data"])

	taskDetail := task{}

	json.Unmarshal(taskJson, &taskDetail)

	return taskDetail, nil
}

func (t *taskFetch) AddTask(categoryID, title, description, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	jsonData := map[string]any{
		"category_id": categoryID,
		"title":       title,
		"description": description,
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/tasks/create"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (t *taskFetch) UpdateTask(taskID, title, description, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return -1, err
	}

	jsonData := map[string]any{
		"id":          taskObjectID,
		"title":       title,
		"description": description,
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", config.SetUrl("/api/v1/tasks/update"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (t *taskFetch) DeleteTask(taskID, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("DELETE", config.SetUrl("/api/v1/tasks/delete?task_id="+taskID), nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (t *taskFetch) UpdateTaskCategory(taskID, categoryID, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return -1, err
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return -1, err
	}

	jsonData := map[string]any{
		"id":          taskObjectID,
		"category_id": categoryObjectID,
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", config.SetUrl("/api/v1/tasks/category/update"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
