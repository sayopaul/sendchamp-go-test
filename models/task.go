package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UUID        string `gorm:"index:idx_task_uuid,unique;not null"`
	UserID      uint   `gorm:"not null"`
	Name        string `gorm:"size:150;not null"`
	Description string `gorm:"size:500;not null"`
	Status      string `gorm:"default:pending"`
}

// BeforeCreate, run this before creating user
func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	task.UUID = uuid.NewString()
	return
}
