package tests

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	//Mysql
	dsn := "root:518888@tcp(139.159.234.134:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("[InitDB]打开mysql错误:%s", err.Error())
	}
	//mysqlDB.AutoMigrate(&domains.User{})
	//mysqlDB.AutoMigrate(&domains.Good{})
	//mysqlDB.AutoMigrate(&domains.TradeItem{})
	//mysqlDB.AutoMigrate(&domains.UserAndItem{})
	//mysqlDB.AutoMigrate(&domains.Item{})
	return mysqlDB
}

func TestRawCreateUserItem(t *testing.T) {
	db := InitDB()
	itemID := 1
	userID := 58
	//查找用户对应的物品记录是否存在，如果不存在，如果存在，则添加(共同字段id，会冲突)
	m := make(map[string]interface{})
	q1 := db.Raw("select a.*,b.* from user_and_item a inner join users b on a.user_id = b.id where a.item_id = ? AND b.id=?", itemID, userID).Scan(&m)
	if q1.RowsAffected == 0 {
		fmt.Printf("%v", q1.Error)
	}
	//for row.Next() {
	//	row.Scan(&data)
	//}
}

func TestRawCreateUserItem2(t *testing.T) {
	db := InitDB()
	itemID := 1
	userID := 58
	m := make(map[string]interface{})
	tx := db.Raw("select a.item_nums from user_and_item a inner join users b on a.user_id = b.id where a.item_id = ? AND b.id=?", itemID, userID).Scan(&m)
	if tx.RowsAffected == 0 {
		//如果找不到该记录，则添加
		exec := db.Exec("insert into user_and_item (`user_id`,`item_id`,`item_nums`,`lock_nums`) values (?,?,?,0)", userID, itemID, 6)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			db.Rollback()
			return
		}
	} else {
		//否则就是找到了记录，在原来数量添加
		exec := db.Exec("update user_and_item set item_nums = ? where user_id=? and item_id=?", gorm.Expr("item_nums+?", 6), userID, itemID)
		if exec.RowsAffected == 0 {
			//出现问题回滚一切操作
			db.Rollback()
			return
		}
	}
	db.Commit()

}

func TestRaw3(t *testing.T) {
	db := InitDB()
	s1 := 7623
	s2 := 3
	s3 := 3
	db.Debug().Table("users").Where("id=?", s1).Updates(map[string]interface{}{
		"good":    s2 - 1,
		"diamond": s3 - 1,
	})
}
