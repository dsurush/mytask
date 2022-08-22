package service

import (
	"mytasks"
	"mytasks/pkg/repository"
)

type TaskList interface {
	GetAll() ([]mytasks.TaskList, error)
	Delete(id int) error
	Update(id int, input mytasks.UpdateTaskInput) error
}

type Service struct {
	TaskList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{TaskList: NewTodoListService(repos.TaskList)}
}
