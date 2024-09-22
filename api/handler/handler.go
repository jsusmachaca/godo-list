package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"slices"

	"github.com/jsusmachaca/godo/api/response"
	"github.com/jsusmachaca/godo/internal/validation"
	"github.com/jsusmachaca/godo/pkg/file"
	"github.com/jsusmachaca/godo/pkg/model"
	uuid "github.com/satori/go.uuid"
)

var taskFile string

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Please indicate the path of the task file, " +
			"e.g., \033[32m./main $(pwd)/tasks.json\033[0m")
		os.Exit(1)
		return
	}
	taskFile = os.Args[1]

	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(taskFile)
			if err != nil {
				panic(err)
			}
			_, err = file.WriteString("[]\n")
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer file.Close()
}

func Index(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task

	err := file.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
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

func GetAll(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

func AddTask(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task
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

	err = file.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	body.ID = uuid.NewV4().String()

	tasksList = append(tasksList, body)

	err = file.WriteJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	response := response.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task
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

	err = file.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	tasksList = slices.DeleteFunc(tasksList, func(task model.Task) bool {
		return task.ID == body.ID
	})

	err = file.WriteJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task
	var body model.Task
	id := r.PathValue("id")

	err := file.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	err = validation.RequestValidator(r.Body, &body)
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

	for i := 0; i < len(tasksList); i++ {
		if tasksList[i].ID == id {
			tasksList[i].Name = body.Name
			tasksList[i].Done = body.Done
			break
		}
	}

	err = file.WriteJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	response := response.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
