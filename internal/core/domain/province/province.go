package province

import "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"

type Province struct {
	domain.Base

	Name string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code *string `gorm:"type:varchar(50)" json:"code,omitempty"`
}
