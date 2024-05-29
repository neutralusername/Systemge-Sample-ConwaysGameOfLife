package appWebsocket

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"Systemge/Websocket"
)

type App struct {
	messageBrokerClient *MessageBrokerClient.Client
	websocketServer     *Websocket.Server
	name                string
	randomizer          *Utilities.Randomizer
	logger              *Utilities.Logger
}

func New(name string, logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client, websocketServer *Websocket.Server) *App {
	return &App{
		messageBrokerClient: messageBrokerClient,
		websocketServer:     websocketServer,
		name:                name,
		logger:              logger,
	}
}
