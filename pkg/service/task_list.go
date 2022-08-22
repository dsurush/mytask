package service

import (
	"mytasks"
	"mytasks/pkg/repository"
)

type TaskListService struct {
	repo repository.TaskList
}

func NewTodoListService(repo repository.TaskList) *TaskListService {
	return &TaskListService{repo: repo}
}

func (s *TaskListService) GetAll() ([]mytasks.TaskList, error) {
	return s.repo.GetAll()
}

func (s *TaskListService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TaskListService) Update(id int, input mytasks.UpdateTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(id, input)
}
