package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
	"github.com/rafli-lutfi/kanban-app-mongodb/routes"
)

//go:embed src/views/*
var Resource embed.FS

func main() {
	config.LoadEnv()
	config.ConnectDB()

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
		routes.RunClient(mux, Resource)

		fmt.Println("Server Running On Port", port)
		http.ListenAndServe(":"+port, mux)
	}()

	wg.Wait()
}
