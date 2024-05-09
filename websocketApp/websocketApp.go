package websocketApp

import (
	"Systemge/MessageServer"
	"Systemge/Websocket"
)

type WebsocketApp struct {
	websocketServer *Websocket.Server
	messageBroker   MessageServer.Endpoint
}

func New(websocketServer *Websocket.Server, messageBroker MessageServer.Endpoint) Websocket.WebsocketApplication {
	app := &WebsocketApp{
		websocketServer: websocketServer,
		messageBroker:   messageBroker,
	}
	return app
}
