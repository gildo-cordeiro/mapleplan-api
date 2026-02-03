package couple

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
)

type Couple struct {
	domain.Base

	Name  string      `gorm:"type:varchar(100);not null" json:"name"`
	Goals []goal.Goal `gorm:"foreignKey:UserID" json:"goals,omitempty"`
}
