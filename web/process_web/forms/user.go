package forms

type CreateRoomForm struct {
	RoomID        uint32 `form:"room_id" binding:"required"`
	MaxUserNumber uint32 `form:"max_user_number" binding:"required"`
	GameCount     uint32 `form:"game_count" binding:"required"`
	RoomName      string `form:"room_name" binding:"required"`
}
