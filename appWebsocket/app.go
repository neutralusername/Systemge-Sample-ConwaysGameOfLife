package appWebsocket

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) MessageBrokerClient.WebsocketApplication {
	return &App{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
}
