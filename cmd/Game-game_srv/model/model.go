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

// redis存储
type RoomUser struct {
	ID    uint32 `json:"ID"`
	Ready bool   `json:"Ready"`
}

type Room struct {
	RoomID        uint32      `json:"roomID"`
	MaxUserNumber uint32      `json:"maxUserNumber"`
	GameCount     uint32      `json:"gameCount"`
	UserNumber    uint32      `json:"userNumber"`
	RoomOwner     uint32      `json:"roomOwner"`
	RoomWait      bool        `json:"roomWait"`
	Users         []*RoomUser `json:"users"`
	RoomName      string      `json:"roomName"`
}

type User struct {
	BaseModel
	UserName string `gorm:"index;unique;not null"`
	Password string
	Nickname string `gorm:"index;not null"`
	Gender   bool
	Image    string `gorm:"not null;default: /"` //路径存储
	Good     int    `gorm:"not null;comment:金币"`
	Diamond  int    `gorm:"not null;comment:钻石"`
}
