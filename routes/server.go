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
	categoryService := services.NewCategoryService(categoryRepository, taskRepository)
	taskService := services.NewTaskService(taskRepository, categoryRepository)

	userAPIHandler := controllers.NewUserHandler(userService)
	categoryAPIHandler := controllers.NewCategoryHandler(categoryService)
	taskAPIHandler := controllers.NewTaskHandler(taskService)

	APIHandler := APIHandler{
		userHandler:     userAPIHandler,
		categoryHandler: categoryAPIHandler,
		taskHandler:     taskAPIHandler,
	}

	muxRoute(mux, "POST", "/api/v1/users/register", middleware.Post(http.HandlerFunc(APIHandler.userHandler.Register))) // User REGISTER
	muxRoute(mux, "POST", "/api/v1/users/login", middleware.Post(http.HandlerFunc(APIHandler.userHandler.Login)))       // User Login
	muxRoute(mux, "GET", "/api/v1/users/logout", middleware.Get(http.HandlerFunc(APIHandler.userHandler.Logout)))       // user logout

	muxRoute(mux, "GET", "/api/v1/categories/dashboard", middleware.Get(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.GetCategory))))                        // Show Dashboard with categories with their tasks
	muxRoute(mux, "POST", "/api/v1/categories/create", middleware.Post(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.CreateCategory))))                      // Create new category
	muxRoute(mux, "DELETE", "/api/v1/categories/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(APIHandler.categoryHandler.DeleteCategory))), "?category_id=") // Delete category

	muxRoute(mux, "GET", "/api/v1/tasks/get", middleware.Get(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.GetTaskByID))), "?task_id=")                    // Get Task
	muxRoute(mux, "POST", "/api/v1/tasks/create", middleware.Post(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.StoreTask))))                              // Create Task
	muxRoute(mux, "PUT", "/api/v1/tasks/update", middleware.Put(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.UpdateTask))))                               // Update Task
	muxRoute(mux, "PUT", "/api/v1/tasks/category/update", middleware.Put(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.UpdateTaskCategory))), "?task_id=") // Update Task's Category
	muxRoute(mux, "DELETE", "/api/v1/tasks/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(APIHandler.taskHandler.DeleteTask))), "?task_id=")            // Delete Task

	return mux
}
