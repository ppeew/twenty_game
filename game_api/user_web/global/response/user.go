package response

type UserResponse struct {
	Id       uint32 `json:"id"`
	Nickname string `json:"nickname"`
	Gender   bool   `json:"gender"`
	Username string `json:"username"`
}
