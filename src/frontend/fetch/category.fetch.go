package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
)

type CategoryFetch interface {
	GetCategories(userID string) ([]category, int, error)
	AddCategory(categoryType string, userID string) (int, error)
	DeleteCategory(categoryID string, userID string) (int, error)
}

type categoryFetch struct {
}

type task struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  string `json:"category_id"`
	UserId      string `json:"user_id"`
}

type category struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	UserId string `json:"user_id"`
	Tasks  []task `json:"tasks"`
}

func NewCategoryFetch() *categoryFetch {
	return &categoryFetch{}
}

func (c *categoryFetch) GetCategories(userID string) ([]category, int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return nil, -1, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/categories/dashboard"), nil)
	if err != nil {
		return nil, -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, errors.New("status code not 200")
	}

	result := map[string]any{}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, -1, err
	}

	if result["data"] == nil {
		return nil, -1, err
	}

	data := result["data"].([]interface{})

	categoriesJSON, _ := json.Marshal(data)

	categories := []category{}

	json.Unmarshal(categoriesJSON, &categories)

	return categories, http.StatusOK, nil
}

func (c *categoryFetch) AddCategory(categoryType string, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	jsonData := map[string]string{
		"type": categoryType,
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/categories/create"), bytes.NewBuffer(data))
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

func (c *categoryFetch) DeleteCategory(categoryID string, userID string) (int, error) {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("DELETE", config.SetUrl("/api/v1/categories/delete?category_id="+categoryID), nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
