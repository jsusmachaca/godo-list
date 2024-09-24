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
	uuid "github.com/satori/go.uuid"
)

func Index(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var tasksList []model.Task
	var tasks model.Task

	query := `SELECT * FROM tasks;`
	rows, err := db.Query(query)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&tasks.ID,
			&tasks.Name,
			&tasks.Done,
		)
		tasksList = append(tasksList, tasks)
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
	var tasksList []model.Task
	var tasks model.Task

	query := `SELECT * FROM tasks;`

	rows, err := db.Query(query)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&tasks.ID,
			&tasks.Name,
			&tasks.Done,
		)
		tasksList = append(tasksList, tasks)
	}
	if len(tasksList) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksList)
}

func AddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body model.Task

	err := validation.RequestValidator(r.Body, &body)
	if err != nil {
		if errors.Is(err, validation.ErrInvalidDataType) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`{"error": "Invalid type data"}`))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	query := `INSERT INTO tasks VALUES (?, ?, ?);`
	stmt, err := db.Prepare(query)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}
	defer stmt.Close()

	body.ID = uuid.NewV4().String()

	result, err := stmt.Exec(body.ID, body.Name, body.Done)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to insert data"}`))
		return
	}
	if i, err := result.RowsAffected(); err != nil || i != 1 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "1 row was expected to be affected"}`))
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

	query := `DELETE FROM tasks WHERE id=?;`
	stmt, err := db.Prepare(query)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to delete data"}`))
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to delete data"}`))
		return
	}
	if i, err := result.RowsAffected(); err != nil || i != 1 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "1 row was expected to be affected "}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func UpdateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body model.Task
	id := r.PathValue("id")

	err := validation.RequestValidator(r.Body, &body)
	if err != nil {
		if errors.Is(err, validation.ErrInvalidDataType) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`{"error": "Invalid type data"}`))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	query := `UPDATE tasks SET done=? WHERE id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to delete data"}`))
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.Done, id)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to delete data"}`))
		return
	}
	if i, err := result.RowsAffected(); err != nil || i != 1 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "1 row was expected to be affected "}`))
		return
	}

	response := response.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
