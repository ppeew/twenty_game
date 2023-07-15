package dao

import (
	"errors"
	"store_web/global"
	"store_web/service/domain"
	"store_web/service/dto"
)

type ShopDao struct {
}

func (s ShopDao) SelectGoods(req dto.ShopSelectReq) (rsp []domain.Good, err error) {
	err = global.MysqlDB.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&rsp).Error
	return
}

func (s ShopDao) BuyGood(req dto.ShopBuyReq, userID int) error {
	begin := global.MysqlDB.Begin()
	//1.查询出该商品的价格及用什么支付
	good := domain.Good{}
	err := begin.Table("goods").Where("id=?", req.ID).First(&good).Error
	if err != nil {
		begin.Rollback()
		return errors.New("商品不存在")
	}
	newInventory := good.Inventory - req.Num
	if newInventory < 0 {
		//没这么多库存
		begin.Rollback()
		return errors.New("库存不足")
	}
	//2.查询该用户货币数量（用户与货币是一对一关系）
	user := domain.User{}
	begin.Table("users").Where("id=?", userID).First(&user)
	newGood := user.Good - req.Num*good.PriceGood
	newDiamond := user.Good - req.Num*good.PriceDiamond
	if newGood >= 0 && newDiamond >= 0 {
		begin.Table("users").Updates(map[string]interface{}{
			"good":    newGood,
			"diamond": newDiamond,
		})
	} else {
		//买不起
		begin.Rollback()
		return errors.New("货币不够")
	}
	//3.将库存扣减
	err = begin.Table("goods").Where("id=?", req.ID).Update("inventory", newInventory).Error
	if err != nil {
		begin.Rollback()
		return errors.New("服务器错误")
	}
	begin.Commit()
	return nil
}
