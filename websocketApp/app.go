package WebsocketApp

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"Systemge/Websocket"
)

type App struct {
	messageBrokerClient *MessageBrokerClient.Client
	websocketServer     *Websocket.Server
	name                string
	logger              *Utilities.Logger
}

func New(name string, logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client, websocketServer *Websocket.Server) *App {
	return &App{
		name:                name,
		messageBrokerClient: messageBrokerClient,
		logger:              logger,
		websocketServer:     websocketServer,
	}
}
