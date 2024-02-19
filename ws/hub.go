package ws

import (
	"fmt"
	"sync"

	"github.com/danmondy/jailbreak/templates"
)

type Hub struct {
	sync.RWMutex
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	CreateRoom chan *Room
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		CreateRoom: make(chan *Room),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			if _, ok := h.Rooms[c.RoomID]; ok {
				r := h.Rooms[c.RoomID]

				if _, ok := r.Clients[c.ID]; !ok { //only add the player if they are not already there
					r.Clients[c.ID] = c

					fmt.Printf("%v", c)

					r.Host.MessageChan <- &Message{
						Template: templates.PlayerList(r.GetPlayerList()),
						RoomID:   r.ID,
						Username: c.Player.Username,
						MsgType:  TypeHTML,
					}
				}
			}
		case client := <-h.Unregister:
			if r, ok := h.Rooms[client.RoomID]; ok { //room exists
				if _, ok := h.Rooms[client.RoomID].Clients[client.ID]; ok { //client exists
					delete(h.Rooms[client.RoomID].Clients, client.ID)
					close(client.MessageChan)
					r.Host.MessageChan <- &Message{
						Content:  "user left the chat",
						RoomID:   client.RoomID,
						Username: client.Player.Username,
						MsgType:  TypeJSON,
					}
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok { //if room exists
				//todo ensure the player is in the room on the server to avoid attacks
				for _, c := range h.Rooms[m.RoomID].Clients {
					c.MessageChan <- m
				}
			}
		}

	}
}
