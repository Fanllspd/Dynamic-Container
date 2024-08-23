package models

import (
	"time"

	"gorm.io/gorm"
)

type ID struct {
	ID uint `gorm:"primary_key" json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

type SoftDeletes struct {
	DeleteAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}
