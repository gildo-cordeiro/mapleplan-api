package repositories

import (
	taskDomain "github.com/gildo-cordeiro/mapleplan-api/internal/data/models/task"
	repoPort "github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func (g TaskRepositoryImpl) FindById(id uint) (*taskDomain.Task, error) {
	// TODO: implement repository logic
	return nil, nil
}

func (g TaskRepositoryImpl) Save(t *taskDomain.Task) (string, error) {
	// TODO: implement repository logic
	return "", nil
}

func NewGormTaskRepository(db *gorm.DB) repoPort.TaskRepository {
	return &TaskRepositoryImpl{db: db}
}
