package main

import (
	"fmt"
	"net/http"

	"github.com/jsusmachaca/godo/api/handler"
)

func routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handler.Index)
	mux.HandleFunc("GET /api/tasks", handler.GetAll)
	mux.HandleFunc("POST /api/add-task", handler.AddTask)
	mux.HandleFunc("DELETE /api/delete-task", handler.DeleteTask)
	mux.HandleFunc("PUT /api/update-task/{id}", handler.UpdateTask)
}

func main() {
	mux := http.NewServeMux()

	routes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Printf("Server listen on http://localhost%s\n", server.Addr)
	server.ListenAndServe()
}
