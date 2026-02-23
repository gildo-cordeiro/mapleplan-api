package profile

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type ImmigrationProfile struct {
	models.Base

	UserID     string  `gorm:"type:uuid;not null;index" json:"userId"`
	ProvinceID *string `gorm:"type:uuid;index" json:"provinceId,omitempty"`

	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	ProgramType string `gorm:"type:varchar(100);not null" json:"programType"`
	Status      string `gorm:"type:varchar(100);not null" json:"status"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}
