package repository

import (
	//"github.com/google/martian/log"
	"github.com/jmoiron/sqlx"

	"mytasks"
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
	_, err := r.db.Exec("UPDATE tasks SET title=$1, description=$2, done=$3 WHERE id=$4", input.Title, input.Description, input.Done, id)
	return err
}
