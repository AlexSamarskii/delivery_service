package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id          uuid.UUID  `gorm:"column:id;" json:"id"`
	Name        string     `gorm:"column:name;" json:"name"`
	Description string     `gorm:"column:description;" json:"description"`
	Icon        string     `gorm:"column:icon;" json:"icon"`
	Status      string     `gorm:"column:status;" json:"status"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at;"`
}
