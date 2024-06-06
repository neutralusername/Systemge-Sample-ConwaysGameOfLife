package appWebsocket

import (
	"Systemge/Application"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) Application.WebsocketApplication {
	return &App{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
}
