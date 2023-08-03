package dto

type ShopSelectReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ShopBuyReq struct {
	ShopID int `json:"shopID"`
	Num    int `json:"num"`
}
