package domains

import (
	"time"

	"gorm.io/gorm"
)

// 用户表
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

type BaseModel struct {
	ID        uint32 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u User) TableName() string {
	return "users"
}
