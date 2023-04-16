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

type Items struct {
	Apple  uint32
	Banana uint32
}

type UserItem struct {
	BaseModel
	Gold    uint32
	Diamond uint32
	Apple   uint32
	Banana  uint32
}
