package response

type UserResponse struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}
