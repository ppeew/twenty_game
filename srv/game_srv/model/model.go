package model

import (
	"time"

	"gorm.io/gorm"
)

// mysql存储
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
	UserID  uint32 `gorm:"index;unique"` //外键，对应用户id
	Gold    uint32
	Diamond uint32
	Apple   uint32
	Banana  uint32
}

// redis存储
type User struct {
	ID    uint32 `json:"ShopID"`
	Ready bool   `json:"Ready"`
}

type Room struct {
	RoomID        uint32  `json:"roomID"`
	MaxUserNumber uint32  `json:"maxUserNumber"`
	GameCount     uint32  `json:"gameCount"`
	UserNumber    uint32  `json:"userNumber"`
	RoomOwner     uint32  `json:"roomOwner"`
	RoomWait      bool    `json:"roomWait"`
	Users         []*User `json:"users"`
	RoomName      string  `json:"roomName"`
}
