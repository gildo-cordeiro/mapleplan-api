package repositories

import "github.com/gildo-cordeiro/mapleplan-api/internal/data/models/task"

type TaskRepository interface {
	FindById(id uint) (*task.Task, error)
	Save(t *task.Task) (string, error)
}
