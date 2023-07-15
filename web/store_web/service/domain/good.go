package domain

// 商城商品表(商城中的商品一定在物品表中)
type Good struct {
	ID           int `gorm:"primaryKey;autoIncrement;comment:主键"`
	ItemID       int `gorm:"not null;comment:外键(查询物品表)"`
	Inventory    int `gorm:"not null;comment:商品库存"`
	PriceGood    int `gorm:"not null;comment:商品金币价格"`
	PriceDiamond int `gorm:"not null;comment:商品钻石价格"`
}

func (g Good) TableName() string {
	return "goods"
}
