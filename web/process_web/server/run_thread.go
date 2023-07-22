package server

type Data struct {
	roomID        uint32
	maxUserNumber uint32
	gameCount     uint32
	userNumber    uint32
	roomOwner     uint32
	roomName      string
	users         []uint32
	gameMode      int
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
		gameMode:      0,
	}
}

func Run(data *Data) {
	for true {
		room := NewRoomStruct(data)
		d, isQuit := room.RunRoom()
		if isQuit {
			return
		}
		switch d.gameMode {
		case 0:
			game := NewGameData(d)
			game.RunGame()
		}
	}
}
