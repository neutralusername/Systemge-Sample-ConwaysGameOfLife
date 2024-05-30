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

func (app *App) GetAsyncMessageHandlers() map[string]MessageBrokerClient.AsyncMessageHandler {
	return map[string]MessageBrokerClient.AsyncMessageHandler{
		"getGrid":       app.WebsocketPropagate,
		"getGridChange": app.WebsocketPropagate,
	}
}

func (app *App) GetSyncMessageHandlers() map[string]MessageBrokerClient.SyncMessageHandler {
	return map[string]MessageBrokerClient.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]func() error {
	return map[string]func() error{}
}
