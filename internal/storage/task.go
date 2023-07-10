package storage

import (
	"ticket-pimp/internal/domain"
)

type TaskRepository interface {
	ListTasks() ([]*domain.TaskEntity, error)
	SaveTask(task *domain.TaskEntity) (*domain.TaskEntity, error)
	GetTaskByID(id int) (*domain.TaskEntity, error)
	UpdateTask(task *domain.TaskEntity) (*domain.TaskEntity, error)
}

func (s *Storage) SaveTask(task *domain.TaskEntity) (*domain.TaskEntity, error) {
	tx := s.db.Create(task)
	return task, tx.Error
}

func (s *Storage) ListTasks() ([]*domain.TaskEntity, error) {
	var task []*domain.TaskEntity
	tx := s.db.Find(&task)
	return task, tx.Error
}
