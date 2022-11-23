package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	UUID       string    `gorm:"index:idx_uuid,unique;not null"`
	Email      string    `gorm:"index:idx_email,unique;size:100;not null"`
	Password   string    `gorm:"size:256"`
	DateJoined time.Time `gorm:"autoCreateTime;not null"`
	LastLogin  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;not null"`
}

// to return a user in JSOn response and hide sensitive info like password, ID etc.
type FetchUser struct {
	UUID       string    `json:"UUID"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"date_joined"`
	LastLogin  time.Time `json:"last_login"`
}

// BeforeCreate, run this before creating user
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.UUID = uuid.NewString()
	return
}
