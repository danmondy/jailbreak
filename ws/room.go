package ws

import "github.com/danmondy/jailbreak/data"

type Room struct {
	ID        string             `json:"id"`
	Clients   map[string]*Client `json:"clients"`
	Host      *Client            `json:"host"` //this is a special player for the shared monitor
	GameState string             `json:"game-state"`
}

func (r *Room) GetPlayerList() map[string]*data.Player {
	players := make(map[string]*data.Player)
	for _, c := range r.Clients {
		players[c.Player.ID] = c.Player
	}
	return players
}

type GameState struct {
	Stage  int
	Scores map[string]int
}

const (
	GameStateLobby  string = "lobby"
	GameStateInGame string = "in-game"
)
