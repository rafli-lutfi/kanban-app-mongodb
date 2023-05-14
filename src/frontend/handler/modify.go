package web

import (
	"embed"
	"html/template"
	"net/http"
	"path"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/frontend/fetch"
)

type ModifyWeb interface {
	AddCategory(w http.ResponseWriter, r *http.Request)
	AddCategoryProcess(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)

	AddTask(w http.ResponseWriter, r *http.Request)
	AddTaskProcess(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskProcess(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)

	MoveTask(w http.ResponseWriter, r *http.Request)
}

type modifyWeb struct {
	categoryFetch fetch.CategoryFetch
	taskFetch     fetch.TaskFetch
	embed         embed.FS
}

func NewModifyWeb(categoryFetch fetch.CategoryFetch, taskFetch fetch.TaskFetch, embed embed.FS) *modifyWeb {
	return &modifyWeb{categoryFetch, taskFetch, embed}
}

func (m *modifyWeb) AddCategory(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("src", "views", "main", "add-category.html")
	header := path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.ParseFS(m.embed, filepath, header))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *modifyWeb) AddCategoryProcess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")

	categoryType := r.FormValue("type")

	respCode, err := m.categoryFetch.AddCategory(categoryType, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCode == 201 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/category/add", http.StatusSeeOther)
	}
}

func (m *modifyWeb) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	userID := r.Context().Value("token")

	_, err := m.categoryFetch.DeleteCategory(categoryID, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (m *modifyWeb) AddTask(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category")

	filepath := path.Join("src", "views", "main", "add-task.html")
	header := path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.ParseFS(m.embed, filepath, header))

	err := tmpl.Execute(w, map[string]interface{}{"id": categoryID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *modifyWeb) AddTaskProcess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")

	categoryID := r.URL.Query().Get("category")
	title := r.FormValue("title")
	description := r.FormValue("description")

	respCode, err := m.taskFetch.AddTask(categoryID, title, description, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCode == 201 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/task/add?category="+categoryID, http.StatusSeeOther)
	}
}

func (m *modifyWeb) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	userID := r.Context().Value("token")

	taskDetail, err := m.taskFetch.GetTaskByID(taskID, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filepath := path.Join("src", "views", "main", "update-task.html")
	header := path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.ParseFS(m.embed, filepath, header))

	err = tmpl.Execute(w, map[string]any{
		"title":       "Update Task",
		"id":          taskID,
		"task_title":  taskDetail.Title,
		"description": taskDetail.Description})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *modifyWeb) UpdateTaskProcess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")

	taskID := r.URL.Query().Get("task_id")
	categoryID := r.URL.Query().Get("category_id")

	if categoryID == "" {
		title := r.FormValue("title")
		description := r.FormValue("description")

		respCode, err := m.taskFetch.UpdateTask(taskID, title, description, userID.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if respCode == 200 {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/task/update?task_id="+taskID, http.StatusSeeOther)
		}
	} else {
		_, err := m.taskFetch.UpdateTaskCategory(taskID, categoryID, userID.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (m *modifyWeb) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")
	taskID := r.URL.Query().Get("task_id")

	_, err := m.taskFetch.DeleteTask(taskID, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (m *modifyWeb) MoveTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")
	taskID := r.URL.Query().Get("task_id")
	categoryID := r.URL.Query().Get("category_id")

	_, err := m.taskFetch.UpdateTaskCategory(taskID, categoryID, userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
