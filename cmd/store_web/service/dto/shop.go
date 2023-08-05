package dto

type ShopSelectReq struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type ShopBuyReq struct {
	ShopID int `uri:"id"`
	Num    int `json:"num"`
}
