package domain

const (
	GoldType = iota + 1
	DiamondType
	Interact
)

// 物品表
type Item struct {
	ID    int    `gorm:"primaryKey;autoIncrement;comment:主键"`
	Name  string `gorm:"not null;comment:道具的名称"`
	Image string `gorm:"not null;comment:道具的图片（存路径，部署nginx）"`
	Type  int    `gorm:"not null;comment:道具的类型(1:金币 2:钻石 3:交互道具)"`
	Desc  string `gorm:"not null;comment:道具的描述"`
}

func (i Item) TableName() string {
	return "items"
}

// 用户与物品表是多对多关系，因此需要中间表来进行关联查询
type UserAndItem struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	UserID   int `gorm:"not null;comment:外键（用户ID）"`
	ItemID   int `gorm:"not null;comment:外键（物品ID）"`
	ItemNums int `gorm:"not null;comment:用户拥有道具数量"`
	LockNums int `gorm:"not null;comment:在售卖中，不可使用的数量"`
}

func (UserAndItem) TableName() string {
	return "user_and_item"
}
