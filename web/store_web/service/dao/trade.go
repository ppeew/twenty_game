package dao

import (
	"errors"
	"store_web/global"
	"store_web/service/domains"

	"gorm.io/gorm"
)

type TradeDao struct {
}

// 查询所有交易商品
func (d TradeDao) SelectPage(page int, size int) (rsp []domains.TradeItem, err error) {
	off := (page - 1) * size
	err = global.MysqlDB.Offset(off).Limit(size).Find(&rsp).Error
	return
}

// 上架物品
func (d TradeDao) Create(data domains.TradeItem) (domains.TradeItem, error) {
	begin := global.MysqlDB.Begin()
	//1.创建商家物品记录
	if tx := begin.Create(&data); tx.RowsAffected == 0 {
		begin.Rollback()
		return data, errors.New("服务器错误")
	}
	//2.看是否有足够物品
	var m domains.UserAndItem
	tx := begin.Raw("select * from user_and_item where user_id=? and item_id=?", data.UserID, data.ItemID).Scan(&m)
	if tx.Error != nil || tx.RowsAffected == 0 {
		//没查询到或者其他错误
		begin.Rollback()
		return data, errors.New("服务器错误")
	}
	if m.ItemNums-m.LockNums-data.Count < 0 {
		begin.Rollback()
		return data, errors.New("商品数量不足上架")
	}
	//3.锁定用户物品
	exec := begin.Exec("update user_and_item set lock_nums = ? where user_id = ? and item_id = ?", gorm.Expr("lock_nums+?", data.Count), data.UserID, data.ItemID)
	if exec.Error != nil {
		begin.Rollback()
		return data, errors.New("服务器错误")
	}
	begin.Commit()
	return data, nil
}

// 下架物品
func (d TradeDao) DeleteByID(id int, userID uint32) error {
	begin := global.MysqlDB.Begin()
	// 1.得查询原来物品被锁定的数量
	tradeItem := domains.TradeItem{}
	begin.Where("id = ?", id).Where("user_id = ?", userID).First(&tradeItem)
	// 2.删除上架表记录
	if tx := begin.Where("id = ?", id).Where("user_id = ?", userID).Delete(&domains.TradeItem{}); tx.RowsAffected == 0 {
		begin.Rollback()
		return errors.New("找不到该上架物品")
	}
	// 3.解锁物品
	exec := begin.Exec("update user_and_item set lock_nums = ? where user_id = ? and item_id = ?", gorm.Expr("lock_nums - ?", tradeItem.Count), userID, tradeItem.ItemID)
	if exec.Error != nil {
		//找不到记录,也即是用户没有该物品
		begin.Rollback()
		return errors.New("没找到用户在user_and_item信息")
	}
	begin.Commit()
	return nil
}

// 买交易物品
func (d TradeDao) BuyTradeItem(TradeID int, userID int) error {
	begin := global.MysqlDB.Begin()
	//1.查询出该商品的价格及用什么支付
	tradeItem := domains.TradeItem{}
	err := begin.Table("trade_items").Where("id=?", TradeID).Where("status=?", domains.NotBuy).First(&tradeItem).Error
	if err != nil {
		begin.Rollback()
		return errors.New("该物品已售出")
	}
	//2.查询该用户货币数量（用户与货币是一对一关系）
	user := domains.User{}
	begin.Table("users").Where("id=?", userID).First(&user)
	if user.Good >= tradeItem.PriceGood && user.Diamond >= tradeItem.PriceDiamond {
		begin.Table("users").Where("id=?", userID).Updates(map[string]interface{}{
			"good":    user.Good - tradeItem.PriceGood,
			"diamond": user.Diamond - tradeItem.PriceDiamond,
		})
	} else {
		//买不起
		begin.Rollback()
		return errors.New("货币不够")
	}
	//3.将交易记录表更新
	err = begin.Table("trade_items").Where("id=?", TradeID).Update("status", domains.HasBuy).Error
	if err != nil {
		begin.Rollback()
		return errors.New("服务器错误")
	}
	//4.给用户添加物品（用户与物品是多对多关系，用户可以对应多个物品记录，物品记录可以有多个用户持有）要考虑用户与物品连接表记录不存在的情况，不存在则增加记录，存在则给物品数量添加
	m := make(map[string]interface{})
	tx := begin.Raw("select a.item_nums from user_and_item a inner join users b on a.user_id = b.id where a.item_id = ? AND b.id=?", tradeItem.ItemID, userID).Scan(&m)
	if tx.RowsAffected == 0 {
		//如果找不到该记录，则添加
		exec := begin.Exec("insert into user_and_item (`user_id`,`item_id`,`item_nums`,`lock_nums`) values (?,?,?,0)", userID, tradeItem.ItemID, tradeItem.Count)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			begin.Rollback()
			return errors.New("服务器错误")
		}
	} else {
		//否则就是找到了记录，在原来数量添加
		exec := begin.Exec("update user_and_item set item_nums = ? where user_id=? and item_id=?", gorm.Expr("item_nums + ?", tradeItem.Count), userID, tradeItem.ItemID)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			begin.Rollback()
			return errors.New("服务器错误")
		}
	}
	//5.将售卖人的商品数量减少
	begin.Exec("update user_and_item set item_nums=? , lock_nums=? where user_id=? and item_id=?", gorm.Expr("item_nums - ?", tradeItem.Count), gorm.Expr("lock_nums - ?", tradeItem.Count), tradeItem.UserID, tradeItem.ItemID)
	begin.Commit()
	return nil
}
