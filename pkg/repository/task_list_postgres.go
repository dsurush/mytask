package repository

import (
	"fmt"
	//"github.com/google/martian/log"
	"github.com/jmoiron/sqlx"

	"mytasks"
	"strings"
)

type TaskListPostgres struct {
	db *sqlx.DB
}

func NewTaskListPostgres(db *sqlx.DB) *TaskListPostgres {
	return &TaskListPostgres{db: db}
}

func (r *TaskListPostgres) GetAll() ([]mytasks.TaskList, error) {
	var lists []mytasks.TaskList
	query := "SELECT * FROM tasks"
	err := r.db.Select(&lists, query)
	return lists, err
}

func (r *TaskListPostgres) Delete(id int) error {
	query := "DELETE FROM tasks where id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TaskListPostgres) Update(id int, input mytasks.UpdateTaskInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	setQuery := strings.Join(setValues, ",")
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}
