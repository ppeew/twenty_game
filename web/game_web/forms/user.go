package forms

type CreateRoomForm struct {
	RoomID        int    `form:"room_id"`
	MaxUserNumber int    `form:"max_user_number"`
	GameCount     int    `form:"game_count"`
	RoomName      string `form:"room_name"`
}
