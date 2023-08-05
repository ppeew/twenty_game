package utils

type UserInfo struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Gender   bool   `json:"gender"`
	Username string `json:"username"`
	Image    string `json:"image"`
}
