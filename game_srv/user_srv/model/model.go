package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Name   string `gorm:"unique;type:varchar(11);not null"`
	OpenID string `gorm:"index;not null"`
	Gender string `gorm:"not null;"`
}
