package task

import (
	"database/sql"
	"errors"

	"github.com/jsusmachaca/godo/pkg/model"
)

type TaskRepository struct {
	DB *sql.DB
}

func (taskRepository *TaskRepository) GetAll() ([]model.Task, error) {
	var tasksList []model.Task
	var task model.Task

	query := `SELECT * FROM tasks;`

	rows, err := taskRepository.DB.Query(query)
	if err != nil {
		return tasksList, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Done,
		); err != nil {
			return tasksList, err
		}
		tasksList = append(tasksList, task)
	}

	return tasksList, nil
}

func (taskRepository *TaskRepository) Filter(id string) (model.Task, error) {
	var task model.Task

	query := `SELECT * FROM tasks WHERE id=?;`

	rows, err := taskRepository.DB.Query(query, id)
	if err != nil {
		return task, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Done,
		); err != nil {
			return task, err
		}
	} else {
		return task, nil
	}

	return task, nil
}

func (taskRepository *TaskRepository) Insert(id string, name string, done bool) error {
	query := "INSERT INTO tasks VALUES(?, ?, ?);"

	stmt, err := taskRepository.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, name, done)
	if err != nil {
		return err
	}

	if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("1 row was expected to be affected")
	}

	return nil
}
