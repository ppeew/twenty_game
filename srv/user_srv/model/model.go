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
	UserName string `gorm:"index;unique;not null"`
	Password string
	Nickname string `gorm:"index;not null"`
	Gender   bool
	//用户状态，经常修改，不设置索引
	UserState uint32
}
