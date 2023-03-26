package routes

import (
	"fmt"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/controllers"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/middleware"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/repository"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func muxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}

func RunServer(mux *http.ServeMux, db *mongo.Database) *http.ServeMux {
	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userAPIHandler := controllers.NewUserHandler(userService)

	muxRoute(mux, "POST", "/api/v1/user/register", middleware.Post(http.HandlerFunc(userAPIHandler.Register))) // User REGISTER
	muxRoute(mux, "POST", "/api/v1/user/login", middleware.Post(http.HandlerFunc(userAPIHandler.Login)))
	return mux
}
