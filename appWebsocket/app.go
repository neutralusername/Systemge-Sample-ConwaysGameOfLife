package appWebsocket

import (
	"Systemge/RequestServer"
	"Systemge/Websocket"
)

type App struct {
	websocketServer *Websocket.Server

	messageBroker RequestServer.Endpoint
}

func New(websocketServer *Websocket.Server, messageBroker RequestServer.Endpoint) Websocket.WebsocketApplication {
	app := &App{
		websocketServer: websocketServer,

		messageBroker: messageBroker,
	}
	return app
}
