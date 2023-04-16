package response

type RoomResponse struct {
	RoomID     int `json:"roomID"`
	RoomNumber int `json:"roomNumber"`
	GameCount  int `json:"gameCount"`
	RoomOwner  int `json:"roomOwner"`
}
