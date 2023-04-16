package models

type Room struct {
	RoomID       int `json:"roomID"`
	PeopleNumber int `json:"roomNumber"`
	GameCount    int `json:"gameCount"`
	RoomOwner    int `json:"roomOwner"`
}
