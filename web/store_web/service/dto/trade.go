package dto

import (
	"store_web/service/domain"
)

type TradeSelectReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type TradePushReq struct {
	UserID       int
	ItemID       int
	GoodPrice    int
	DiamondPrice int
	Count        int
	Desc         string
}

func (r TradePushReq) ToDomain(userID int) domain.TradeItem {
	return domain.TradeItem{
		UserID:       userID,
		ItemID:       r.ItemID,
		PriceGood:    r.GoodPrice,
		PriceDiamond: r.DiamondPrice,
		Count:        r.Count,
		Status:       domain.NotBuy,
		Desc:         r.Desc,
	}
}

type TradeDownReq struct {
	TradeItemID int
}

type TradeBuyReq struct {
	TradeItemID int
}
