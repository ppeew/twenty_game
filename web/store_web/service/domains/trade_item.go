package domains

// 交易市场物品售卖表
type TradeItem struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	UserID       int    `gorm:"not null;comment:拥有者"`
	ItemID       int    `gorm:"not null;comment:售卖的item"`
	PriceGood    int    `gorm:"not null;comment:售卖金币总价"`
	PriceDiamond int    `gorm:"not null;comment:售卖钻石总价"`
	Count        int    `gorm:"not null;comment:售卖该物品数量"`
	Status       int    `gorm:"not null;comment:交易信息，0:未成交 1:已成交"`
	Desc         string `gorm:"not null;comment:售卖人填写描述"`
}

func (t TradeItem) TableName() string {
	return "trade_items"
}

const (
	NotBuy = iota
	HasBuy
)
