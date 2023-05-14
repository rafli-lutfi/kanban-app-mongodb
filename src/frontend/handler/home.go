package web

import (
	"embed"
	"html/template"
	"net/http"
	"path"
)

type HomeWeb interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type homeWeb struct {
	embed embed.FS
}

func NewHomeWeb(embed embed.FS) *homeWeb {
	return &homeWeb{embed}
}

func (h *homeWeb) Index(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("src", "views", "main", "index.html")
	var header = path.Join("src", "views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(h.embed, filepath, header))

	err := tmpl.Execute(w, map[string]any{"title": "Kanban App"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
