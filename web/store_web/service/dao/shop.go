package dao

import (
	"errors"
	"store_web/global"
	"store_web/service/domains"
	"store_web/service/dto"

	"gorm.io/gorm"
)

type ShopDao struct {
}

func (s ShopDao) SelectGoods(req dto.ShopSelectReq) (rsp []domains.Good, err error) {
	err = global.MysqlDB.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&rsp).Error
	return
}

func (s ShopDao) BuyGood(req dto.ShopBuyReq, userID int) error {
	begin := global.MysqlDB.Begin()
	//1.查询出该商品的价格及用什么支付
	good := domains.Good{}
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
	//2.查询该用户货币数量
	user := domains.User{}
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
	//4.给用户添加物品（用户与物品是多对多关系，用户可以对应多个物品记录，物品记录可以有多个用户持有）
	m := make(map[string]interface{})
	tx := begin.Raw("select a.item_nums from user_and_item a inner join users b on a.user_id = b.id where a.item_id = ? AND b.id=?", req.ID, userID).Scan(&m)
	if tx.RowsAffected == 0 {
		//如果找不到该记录，则添加
		exec := begin.Exec("insert into user_and_item (`user_id`,`item_id`,`item_nums`,`lock_nums`) values (?,?,?,0)", userID, req.ID, req.Num)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			begin.Rollback()
			return errors.New("服务器错误")
		}
	} else {
		//否则就是找到了记录，在原来数量添加
		exec := begin.Exec("update user_and_item set item_nums = ? where user_id=? and item_id=?", gorm.Expr("item_nums+?", req.Num), userID, req.ID)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			begin.Rollback()
			return errors.New("服务器错误")
		}
	}
	begin.Commit()
	return nil
}
