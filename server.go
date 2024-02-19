package main

import (
	"fmt"

	"github.com/labstack/echo/v4" //framework

	"github.com/danmondy/jailbreak/data"
	"github.com/danmondy/jailbreak/ws"
)

const CONFIG_FILE string = "config.json"

var cfg *data.Config

func main() {
	e := echo.New()

	e.Static("/assets", "assets")

	RegisterRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}

func init() {
	cfg = &data.Config{}

	err := data.ReadConfig(CONFIG_FILE, cfg)
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterRoutes(e *echo.Echo) {

	hub := ws.NewHub() //web socket manager
	handler := ws.NewHandler(hub)
	go hub.Run() //starts the routine the reads from the websockets

	e.GET("/ws/host/:roomId", handler.HostRoom)
	e.GET("/ws/player/:roomId", handler.PlayerSocket)

	e.GET("/", handler.WelcomeHandler)
	e.POST("/", handler.PlayerJoin)
	e.GET("/room/:roomId", handler.PlayerLobby)
	e.GET("/host/lobby", handler.HostLobbyHandler)

	//--partials--
	//e.GET("/player_list/:game_id", PlayerListHandler)
}
