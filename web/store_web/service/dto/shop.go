package dto

type ShopSelectReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ShopBuyReq struct {
	ID  int `json:"ID"`
	Num int `json:"Num"`
}
