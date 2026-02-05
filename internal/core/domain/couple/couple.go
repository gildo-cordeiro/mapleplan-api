package couple

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/province"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
)

type Couple struct {
	domain.Base
	ProvinceID *string `gorm:"type:uuid;index" json:"provinceId,omitempty"`
	UserAID    *string `gorm:"type:uuid;index;uniqueIndex:idx_couple_users" json:"userAId,omitempty"`
	UserBID    *string `gorm:"type:uuid;index;uniqueIndex:idx_couple_users" json:"userBId,omitempty"`

	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	DateOfBirth *string `gorm:"type:date" json:"dateOfBirth,omitempty"`

	UserA *user.User `gorm:"foreignKey:UserAID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"userA,omitempty"`
	UserB *user.User `gorm:"foreignKey:UserBID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"userB,omitempty"`

	Province *province.Province `gorm:"foreignKey:ProvinceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"province,omitempty"`
}
