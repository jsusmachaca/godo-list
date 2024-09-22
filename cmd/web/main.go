package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jsusmachaca/godo/api/handler"
	"github.com/jsusmachaca/godo/internal/config"
)

var db *sql.DB

func init() {
	var err error
	db, err = config.GetConnection()
	if err != nil {
		log.Fatal("failed to connect to database, please check if file exists")
	}

	err = config.Migrate(db)
	if err != nil {
		log.Fatal("failed to migrate database, please check if file exists")
	}
}

func routes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /", handler.Index)
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAll(db, w, r)
	})
	mux.HandleFunc("POST /api/add-task", handler.AddTask)
	mux.HandleFunc("DELETE /api/delete-task", handler.DeleteTask)
	mux.HandleFunc("PUT /api/update-task/{id}", handler.UpdateTask)
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	routes(mux, db)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Printf("Server listen on http://localhost%s\n", server.Addr)
	server.ListenAndServe()
}
