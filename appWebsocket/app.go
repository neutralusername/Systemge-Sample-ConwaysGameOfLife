package appWebsocket

import (
	"Systemge/Application"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
)

type WebsocketApp struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) Application.WebsocketApplication {
	return &WebsocketApp{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
}

func (app *WebsocketApp) OnStart() error {
	return nil
}

func (app *WebsocketApp) OnStop() error {
	return nil
}
