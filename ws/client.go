package ws

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/danmondy/jailbreak/data"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	Conn        *websocket.Conn
	MessageChan chan *Message
	ID          string
	Player      *data.Player
	RoomID      string
}

type Message struct {
	Content  string          `json:"content"`
	RoomID   string          `json:"roomId"`
	UserID   string          `json:"userId"`
	Username string          `json:"username"`
	MsgType  MsgType         `json:"msgType"`
	Template templ.Component `json:" - "`
}

type MsgType int

const (
	TypeJSON MsgType = iota
	TypeHTML
	TypeMsg
)

func (c *Client) WriteMessage() { //pass a message to the front end player

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.MessageChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("Write MessageChan was not ok")
				return
			}

			m := message.MsgType
			switch m {
			case TypeHTML:
				buff := &bytes.Buffer{}
				message.Template.Render(context.Background(), buff)
				c.Conn.WriteMessage(websocket.TextMessage, buff.Bytes())
			case TypeJSON:
				c.Conn.WriteJSON(message)
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}

func (c *Client) ReadMessage(hub *Hub) { // read from player
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			UserID:   c.Player.ID,
			Username: c.Player.Username,
		}

		hub.Broadcast <- msg
	}
}
