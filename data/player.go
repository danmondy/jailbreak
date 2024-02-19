package data

type Player struct {
	ID       string `json:"id"`
	Username string `json:"roomId"`
	Score    int    `json:"score"`
}

func NewPlayer(username string) *Player {
	return &Player{ID: NewUniqueID(), Username: username, Score: 0}
}
