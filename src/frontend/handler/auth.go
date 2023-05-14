package web

import (
	"embed"
	"net/http"
	"path"
	"text/template"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/frontend/fetch"
)

type AuthWeb interface {
	Register(w http.ResponseWriter, r *http.Request)
	RegisterProcess(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	LoginProcess(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authWeb struct {
	userFetch fetch.UserFetch
	embed     embed.FS
}

func NewAuthWeb(userFetch fetch.UserFetch, embed embed.FS) *authWeb {
	return &authWeb{userFetch, embed}
}

func (a *authWeb) Register(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("src", "views", "auth", "register.html")
	var header = path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.ParseFS(a.embed, filepath, header))

	data := map[string]any{
		"title": "Register Kanban App",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *authWeb) RegisterProcess(w http.ResponseWriter, r *http.Request) {
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	respCode, err := a.userFetch.Register(fullname, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCode == 201 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

func (a *authWeb) Login(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("src", "views", "auth", "login.html")
	var header = path.Join("src", "views", "general", "header.html")

	tmpl := template.Must(template.ParseFS(a.embed, filepath, header))

	data := map[string]any{
		"title": "Login Kanban App",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *authWeb) LoginProcess(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, respCode, err := a.userFetch.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == 200 {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    id,
			Path:     "/",
			MaxAge:   3600 * 24,
			Domain:   "localhost",
			HttpOnly: true,
			Secure:   false,
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (a *authWeb) Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("token")

	err := a.userFetch.Logout(userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
