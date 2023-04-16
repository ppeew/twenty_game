package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint32 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Nickname string
	Gender   bool
	UserName string `gorm:"index;unique"`
	Password string
}
