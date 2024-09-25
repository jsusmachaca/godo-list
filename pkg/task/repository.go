package task

import (
	"database/sql"
	"errors"

	"github.com/jsusmachaca/godo/pkg/model"
	uuid "github.com/satori/go.uuid"
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

func (taskRepository *TaskRepository) Insert(body *model.Task) error {
	body.ID = uuid.NewV4().String()

	query := "INSERT INTO tasks VALUES(?, ?, ?);"

	stmt, err := taskRepository.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.ID, body.Name, body.Done)
	if err != nil {
		return err
	}

	i, err := result.RowsAffected()
	if err != nil || i != 1 {
		return err
	}
	if i != 1 {
		return errors.New("1 row was expected to be affected")
	}

	return nil
}

func (taskRepository *TaskRepository) Delete(id string) error {
	query := "DELETE FROM tasks WHERE id=?"

	stmt, err := taskRepository.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	i, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if i != 1 {
		return errors.New("1 row was expected to be affected")
	}

	return nil
}

func (taskRepository *TaskRepository) Update(id string, body *model.Task) error {
	query := `UPDATE tasks SET done=? WHERE id=?`
	stmt, err := taskRepository.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.Done, id)
	if err != nil {
		return err
	}

	i, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if i != 1 {
		return errors.New("1 row was expected to be affected")
	}

	return nil
}
