package model

type UserItem struct {
	ID      int
	Gold    int
	Diamond int
	Apple   int
	Banana  int
	Ready   bool
	*WSConn
}
