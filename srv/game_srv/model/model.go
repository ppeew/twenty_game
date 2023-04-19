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

type Room struct {
	RoomID        uint32 `json:"roomID"`
	MaxUserNumber uint32 `json:"maxUserNumber"`
	GameCount     uint32 `json:"gameCount"`
	UserNumber    uint32 `json:"userNumber"`
	RoomOwner     uint32 `json:"roomOwner"`
	RoomWait      bool   `json:"roomWait"`
}
