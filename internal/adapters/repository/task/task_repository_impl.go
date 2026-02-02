package task

import (
	taskDomain "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/task"
	repoPort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func (g RepositoryImpl) FindById(id uint) (*taskDomain.Task, error) {
	// TODO: implement repository logic
	return nil, nil
}

func (g RepositoryImpl) Save(t *taskDomain.Task) (string, error) {
	// TODO: implement repository logic
	return "", nil
}

func NewGormTaskRepository(db *gorm.DB) repoPort.TaskRepository {
	return &RepositoryImpl{db: db}
}
