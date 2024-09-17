package handler

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"slices"

	"github.com/jsusmachaca/godo/pkg/model"
	"github.com/jsusmachaca/godo/pkg/util"
	uuid "github.com/satori/go.uuid"
)

var taskFile = os.Args[1]

func init() {
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(taskFile)
			if err != nil {
				panic(err)
			}
			_, err = file.WriteString("[]")
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

	err := util.ReadJson(taskFile, &tasksList)
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

func GetAll(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task

	err := util.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksList)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task
	var body model.Task

	err := util.RequestValidator(r.Body, &body)
	if err != nil {
		if errors.Is(err, util.ErrInvalidDataType) {
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

	err = util.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	body.ID = uuid.NewV4().String()

	tasksList = append(tasksList, body)

	err = util.WriteJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	response := model.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var tasksList []model.Task
	var body model.Task

	err := util.RequestValidator(r.Body, &body)
	if err != nil {
		if errors.Is(err, util.ErrInvalidDataType) {
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

	err = util.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	tasksList = slices.DeleteFunc(tasksList, func(task model.Task) bool {
		return task.ID == body.ID
	})

	err = util.WriteJson(taskFile, &tasksList)
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

	err := util.ReadJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	err = util.RequestValidator(r.Body, &body)
	if err != nil {
		if errors.Is(err, util.ErrInvalidDataType) {
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

	err = util.WriteJson(taskFile, &tasksList)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error to parsing data"}`))
		return
	}

	response := model.Response{
		Success: true,
		Data:    body,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
