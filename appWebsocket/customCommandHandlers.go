package appWebsocket

import "Systemge/MessageBrokerClient"

func (app *App) GetCustomCommandHandlers() map[string]MessageBrokerClient.CustomCommandHandler {
	return map[string]MessageBrokerClient.CustomCommandHandler{}
}
