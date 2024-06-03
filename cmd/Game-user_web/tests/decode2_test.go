package tests

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
	"unicode/utf8"
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
	Image    string `gorm:"not null;default: /"` //路径存储
	Good     int    `gorm:"not null;comment:金币"`
	Diamond  int    `gorm:"not null;comment:钻石"`
}

func TestDecode2(t *testing.T) {
	str := "\xE8\x8E\xAB"
	r, size := utf8.DecodeRuneInString(str)
	fmt.Println(string(r), " ", size)
	fmt.Println(str)

	//Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "518888", "139.159.234.134", 3306, "game")
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("[InitDB]打开mysql错误:%s", err.Error())
	}

	u := &User{
		BaseModel: BaseModel{},
		UserName:  "112222",
		Password:  "1111",
		//Nickname:  "\xE8\x8E\xAB",
		Nickname: string(r),
		Gender:   false,
		Image:    "1111",
		Good:     0,
		Diamond:  0,
	}
	err = db.Model(&User{}).Create(u).Error
	if err != nil {
		fmt.Println(err)
	}
}
