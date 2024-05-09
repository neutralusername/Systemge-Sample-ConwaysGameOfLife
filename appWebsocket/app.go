package appWebsocket

import (
	"Systemge/MessageServer"
	"Systemge/Websocket"
)

type App struct {
	websocketServer *Websocket.Server

	messageBroker MessageServer.Endpoint
}

func New(websocketServer *Websocket.Server, messageBroker MessageServer.Endpoint) Websocket.WebsocketApplication {
	app := &App{
		websocketServer: websocketServer,

		messageBroker: messageBroker,
	}
	return app
}
