package web

import (
	"embed"
	"log"
	"net/http"
	"path"
	"text/template"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/frontend/fetch"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryFetch fetch.CategoryFetch
	embed         embed.FS
}

func NewDashboardWeb(categoryFetch fetch.CategoryFetch, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{categoryFetch, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")

	categories, respCode, err := d.categoryFetch.GetCategories(userID.(string))
	if err != nil {
		log.Println("error get cat: ", err.Error())
		http.Error(w, err.Error(), respCode)
		return
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
	}

	var filepath = path.Join("src", "views", "main", "dashboard.html")
	var header = path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, filepath, header))

	data := map[string]any{
		"title":      "dashboard Kanban App",
		"categories": categories,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
