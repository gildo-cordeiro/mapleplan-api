package repositories

import "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/task"

type TaskRepository interface {
	FindById(id uint) (*task.Task, error)
	Save(t *task.Task) (string, error)
}
