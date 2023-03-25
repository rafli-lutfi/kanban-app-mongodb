package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
	"github.com/rafli-lutfi/kanban-app-mongodb/routes"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
}

func main() {
	var mux = http.NewServeMux()
	var db = config.GetDBConnection()
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routes.RunServer(mux, db)

	fmt.Println("Server Running On Port", port)
	http.ListenAndServe(":"+port, mux)
}
