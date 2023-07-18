package server

type Data struct {
	roomID        uint32
	maxUserNumber uint32
	gameCount     uint32
	userNumber    uint32
	roomOwner     uint32
	roomName      string
	users         []uint32
}

func NewData(roomID, maxUserNumber, gameCount, userNumber, roomOwner uint32, roomName string, users []uint32) *Data {
	return &Data{
		roomID:        roomID,
		maxUserNumber: maxUserNumber,
		gameCount:     gameCount,
		userNumber:    userNumber,
		roomOwner:     roomOwner,
		roomName:      roomName,
		users:         users,
	}
}

func Run(data *Data) {
	for true {
		room := NewRoomStruct(data)
		isQuit := room.RunRoom(data)
		if isQuit {
			return
		}
		game := NewGameData(data)
		game.RunGame()
	}
}
