package domains

import "time"

// Comments 留言表
type Comments struct {
	ID      int       `gorm:"primaryKey;autoIncrement;comment:主键"`
	UserID  int       `gorm:"not null;comment:外键(查询用户表)"`
	Time    time.Time `gorm:"not null;comment:评论时间"`
	Like    int       `gorm:"not null;comment:点赞数"`
	Content string    `gorm:"not null;comment:评论内容"`
}

func (c Comments) TableName() string {
	return "comments"
}
