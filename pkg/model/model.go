package model

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

type Response struct {
	Success bool `json:"success"`
	Data    Task `json:"data"`
}
