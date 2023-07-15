package dao

import (
	"errors"
	"store_web/global"
	"store_web/service/domain"
)

type TradeDao struct {
}

// 查询所有交易商品
func (d TradeDao) SelectPage(page int, size int) (rsp []domain.TradeItem, err error) {
	off := (page - 1) * size
	err = global.MysqlDB.Offset(off).Limit(size).Find(&rsp).Error
	return
}

// 上架物品
func (d TradeDao) Create(data domain.TradeItem) (domain.TradeItem, error) {
	tx := global.MysqlDB.Create(&data)
	return data, tx.Error
}

// 下架物品
func (d TradeDao) DeleteByUserID(id int, userID uint32) error {
	tx := global.MysqlDB.Where("id=?", id).Where("userID=?", userID).Delete(&domain.TradeItem{})
	return tx.Error
}

// 买交易物品
func (d TradeDao) BuyTradeItem(TradeID int, userID int) error {
	begin := global.MysqlDB.Begin()
	//1.查询出该商品的价格及用什么支付
	tradeItem := domain.TradeItem{}
	err := begin.Table("trade_items").Where("id=?", TradeID).Where("status=?", domain.NotBuy).First(&tradeItem).Error
	if err != nil {
		begin.Rollback()
		return errors.New("该物品已售出")
	}
	//2.查询该用户货币数量（用户与货币是一对一关系）
	user := domain.User{}
	begin.Table("users").Where("id=?", userID).First(&user)
	if user.Good >= tradeItem.PriceGood && user.Diamond >= tradeItem.PriceDiamond {
		begin.Table("users").Updates(map[string]interface{}{
			"good":    user.Good - tradeItem.PriceGood,
			"diamond": user.Diamond - tradeItem.PriceDiamond,
		})
	} else {
		//买不起
		begin.Rollback()
		return errors.New("货币不够")
	}
	//3.将交易记录表更新
	err = begin.Table("trade_items").Where("id=?", TradeID).Update("status", domain.HasBuy).Error
	if err != nil {
		begin.Rollback()
		return errors.New("服务器错误")
	}
	begin.Commit()
	return nil
}
