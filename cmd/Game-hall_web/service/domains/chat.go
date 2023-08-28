package domains

type Message struct {
	Id      int         `json:"id"`
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

type Chat struct {
	Id       int    `form:"id" json:"id"`
	Userid   int    `form:"user_id" json:"user_id"`
	Nickname string `form:"nickName" json:"nickName"`
	Image    string `form:"image" json:"image"`
	Time     string `form:"time" json:"time"`
	Content  string `form:"content" json:"content"`
}
