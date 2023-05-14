package routes

import (
	"embed"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/frontend/fetch"
	web "github.com/rafli-lutfi/kanban-app-mongodb/src/frontend/handler"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/middleware"
)

func RunClient(mux *http.ServeMux, embed embed.FS) *http.ServeMux {
	userFetch := fetch.NewUserFetch()
	categoryFetch := fetch.NewCategoryFetch()
	taskFetch := fetch.NewTaskFetch()

	authWeb := web.NewAuthWeb(userFetch, embed)
	dashboardWeb := web.NewDashboardWeb(categoryFetch, embed)
	modifyWeb := web.NewModifyWeb(categoryFetch, taskFetch, embed)
	homeWeb := web.NewHomeWeb(embed)

	mux.HandleFunc("/register", authWeb.Register)
	mux.HandleFunc("/register/process", authWeb.RegisterProcess)

	mux.HandleFunc("/login", authWeb.Login)
	mux.HandleFunc("/login/process", authWeb.LoginProcess)

	mux.Handle("/logout", middleware.Auth2(http.HandlerFunc(authWeb.Logout)))

	mux.Handle("/dashboard", middleware.Auth2(http.HandlerFunc(dashboardWeb.Dashboard)))

	mux.Handle("/category/add", middleware.Auth2(http.HandlerFunc(modifyWeb.AddCategory)))
	mux.Handle("/category/create", middleware.Auth2(http.HandlerFunc(modifyWeb.AddCategoryProcess)))
	mux.Handle("/category/delete", middleware.Auth2(http.HandlerFunc(modifyWeb.DeleteCategory)))

	mux.Handle("/task/add", middleware.Auth2(http.HandlerFunc(modifyWeb.AddTask)))
	mux.Handle("/task/create", middleware.Auth2(http.HandlerFunc(modifyWeb.AddTaskProcess)))
	mux.Handle("/task/update", middleware.Auth2(http.HandlerFunc(modifyWeb.UpdateTask)))
	mux.Handle("/task/update/process", middleware.Auth2(http.HandlerFunc(modifyWeb.UpdateTaskProcess)))
	mux.Handle("/task/delete", middleware.Auth2(http.HandlerFunc(modifyWeb.DeleteTask)))

	mux.HandleFunc("/", homeWeb.Index)

	return mux
}
