package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/jsusmachaca/godo/api/response"
	"github.com/jsusmachaca/godo/internal/validation"
	"github.com/jsusmachaca/godo/pkg/model"
	"github.com/jsusmachaca/godo/pkg/task"
)

func Index(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	task := task.TaskRepository{DB: db}

	tasksList, err := task.GetAll()
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to obtain data"}`))
		return
	}

	tmpl, err := template.ParseFiles("web/template/index.html")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, map[string][]model.Task{
		"Tasks": tasksList,
	})
}

func GetAll(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	task := task.TaskRepository{DB: db}

	tasksList, err := task.GetAll()
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to obtain data"}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksList)
}

func AddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body model.Task
	task := task.TaskRepository{DB: db}

	err := validation.RequestValidator(r.Body, &body)
	if err != nil {
		var resp map[string]string
		if errors.Is(err, validation.ErrInvalidDataType) {
			resp = map[string]string{"error": "Invalid type data"}
		} else {
			resp = map[string]string{"error": "Error to parsing data"}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = task.Insert(&body)
	if err != nil {
		var resp map[string]string
		if err.Error() == "1 row was expected to be affected" {
			resp = map[string]string{"error": "1 row was expected to be affected"}
		} else {
			resp = map[string]string{"error": "Error to parsing data"}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	response := response.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.PathValue("id")
	task := task.TaskRepository{DB: db}

	err := task.Delete(id)
	if err != nil {
		var resp map[string]string
		if err.Error() == "1 row was expected to be affected" {
			resp = map[string]string{"error": "1 row was expected to be affected"}
		} else {
			resp = map[string]string{"error": "Error to parsing data"}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func UpdateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body model.Task
	id := r.PathValue("id")
	task := task.TaskRepository{DB: db}

	err := validation.RequestValidator(r.Body, &body)
	if err != nil {
		var resp map[string]string
		if errors.Is(err, validation.ErrInvalidDataType) {
			resp = map[string]string{"error": "Invalid type data"}
		} else {
			resp = map[string]string{"error": "Error to parsing data"}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = task.Update(id, &body)
	if err != nil {
		var resp map[string]string
		if err.Error() == "1 row was expected to be affected" {
			resp = map[string]string{"error": "1 row was expected to be affected"}
		} else {
			resp = map[string]string{"error": "Error to parsing data"}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	response := response.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
