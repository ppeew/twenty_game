package domains

// Comments 留言表
type Comments struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	UserID  int    `json:"userId" gorm:"not null;comment:外键(查询用户表)"`
	Time    string `json:"time" gorm:"not null;comment:评论时间"`
	Like    int    `json:"like" gorm:"not null;comment:点赞数"`
	Content string `json:"content" gorm:"not null;comment:评论内容"`
}

func (c Comments) TableName() string {
	return "comments"
}
