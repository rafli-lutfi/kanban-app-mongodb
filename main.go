package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
	"github.com/rafli-lutfi/kanban-app-mongodb/routes"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
}

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		var mux = http.NewServeMux()
		var db = config.GetDBConnection()

		routes.RunServer(mux, db)

		fmt.Println("Server Running On Port", port)
		http.ListenAndServe(":"+port, mux)
	}()

	wg.Wait()
}
