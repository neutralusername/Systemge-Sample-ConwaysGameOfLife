package appWebsocket

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"Systemge/Websocket"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
	websocketServer     *Websocket.Server
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client, websocketServer *Websocket.Server) MessageBrokerClient.WebsocketApplication {
	return &App{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
		websocketServer:     websocketServer,
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
