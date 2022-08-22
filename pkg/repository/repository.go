package repository

import (
	"github.com/jmoiron/sqlx"
	"mytasks"
)

type TaskList interface {
	GetAll() ([]mytasks.TaskList, error)
	Delete(id int) error
	Update(id int, input mytasks.UpdateTaskInput) error
}

type Repository struct {
	TaskList
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{TaskList: NewTaskListPostgres(db)}
}
