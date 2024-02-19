package ws

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/danmondy/jailbreak/data"
	"github.com/danmondy/jailbreak/templates"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // true for now accepts all origins, letter we will specify origin for CORS
		},
	}
)

type Handler struct {
	hub *Hub
}

func (h *Hub) RoomCanBeJoined(roomID string) bool {
	h.RLock()
	defer h.RUnlock()
	r, ok := h.Rooms[roomID]
	if !ok || r.GameState != GameStateLobby {
		return false
	}
	return true
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

// this is used by ever handler to render a view with Echo and Templ
func render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (h *Handler) WelcomeHandler(c echo.Context) error {
	return render(c, 200, templates.WelcomePage())
}

func (h *Handler) HostLobbyHandler(c echo.Context) error {
	code := data.GetRandGameCode(6)
	return render(c, 200, templates.HostLobbyPage(code))
}

func (h *Handler) HostRoom(c echo.Context) (err error) { //this player becomes the "host" player (no longer a participant, intended to be cast to a TV)
	conn, err := websocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	roomID := c.Param("roomId")

	fmt.Printf("Room Code: %s - HostRoom", roomID)

	p := data.NewPlayer("host")

	r := &Room{
		ID:        roomID,
		Clients:   make(map[string]*Client),
		GameState: GameStateLobby,
		Host: &Client{
			Conn:        conn,
			MessageChan: make(chan *Message, 10),
			RoomID:      roomID,
			ID:          data.NewUniqueID(),
			Player:      p,
		},
	}

	h.hub.Lock()
	if _, ok := h.hub.Rooms[r.ID]; !ok { //room does not exist
		h.hub.Rooms[r.ID] = r
		fmt.Println("room created")
	}
	h.hub.Unlock()

	r.Host.MessageChan <- &Message{Content: "Test", MsgType: TypeJSON}

	go r.Host.WriteMessage()

	r.Host.ReadMessage(h.hub)

	return nil
}

func (h *Handler) PlayerJoin(c echo.Context) error {
	roomID := c.FormValue("room_code")
	if h.hub.RoomCanBeJoined(roomID) == false {
		fmt.Println("room not found")
		return c.String(http.StatusNotFound, "Room not found.")
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/room/%s", roomID))
}

func (h *Handler) PlayerLobby(c echo.Context) error {
	roomID := c.Param("roomId")
	return render(c, 200, templates.PlayerLobby(roomID))
}

func (h *Handler) PlayerSocket(c echo.Context) error {
	roomID := c.Param("roomId")
	conn, err := websocketUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		conn.Close()
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !h.hub.RoomCanBeJoined(roomID) {
		fmt.Println("room not found")
		return c.String(http.StatusNotFound, "Room not found.")
	}

	cl := &Client{
		Conn:        conn,
		MessageChan: make(chan *Message, 10),
		RoomID:      roomID,
		ID:          data.NewUniqueID(),
		Player:      data.NewPlayer(""),
	}

	//register a new client through the register channel
	h.hub.Register <- cl

	go cl.WriteMessage()

	cl.ReadMessage(h.hub)

	return nil
}
