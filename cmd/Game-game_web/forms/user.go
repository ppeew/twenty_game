package forms

type CreateRoomForm struct {
	RoomID        int    `form:"room_id" binding:"required"`
	MaxUserNumber int    `form:"max_user_number" binding:"required"`
	GameCount     int    `form:"game_count" binding:"required"`
	RoomName      string `form:"room_name" binding:"required"`
}
