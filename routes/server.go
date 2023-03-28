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

type APIHandler struct {
	userHandler     controllers.UserHandler
	categoryHandler controllers.CategoryHandler
	taskHandler     controllers.TaskHandler
}

func RunServer(mux *http.ServeMux, db *mongo.Database) *http.ServeMux {
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	userService := services.NewUserService(userRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	taskService := services.NewTaskService(taskRepository, categoryRepository)

	userAPIHandler := controllers.NewUserHandler(userService)
	categoryAPIHandler := controllers.NewCategoryHandler(categoryService)
	taskAPIHandler := controllers.NewTaskHandler(taskService)

	APIHandler := APIHandler{
		userHandler:     userAPIHandler,
		categoryHandler: categoryAPIHandler,
		taskHandler:     taskAPIHandler,
	}

	muxRoute(mux, "POST", "/api/v1/user/register", middleware.Post(http.HandlerFunc(APIHandler.userHandler.Register))) // User REGISTER
	muxRoute(mux, "POST", "/api/v1/user/login", middleware.Post(http.HandlerFunc(APIHandler.userHandler.Login)))       // User Login

	muxRoute(mux, "GET", "/api/v1/category/get", middleware.Get(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.GetCategory))))
	muxRoute(mux, "POST", "/api/v1/category/create", middleware.Post(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.CreateCategory))))
	muxRoute(mux, "DELETE", "/api/v1/category/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.DeleteCategory))), "?category_id=")
	// Show Dashboard with categories with their tasks

	// Get Task
	muxRoute(mux, "GET", "/api/v1/task/get", middleware.Get(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.GetTaskByID))), "?task_id=")

	// Create Task
	muxRoute(mux, "POST", "/api/v1/task/create", middleware.Post(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.StoreTask))))

	// Update Task
	muxRoute(mux, "PUT", "/api/v1/task/update", middleware.Put(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.UpdateTask))))

	// Update Task's Category OR Change Task State On Kanban

	// Delete Task
	muxRoute(mux, "DELETE", "/api/v1/task/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.DeleteTask))), "?task_id=")

	return mux
}
