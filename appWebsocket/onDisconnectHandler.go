package appWebsocket

import (
	"Systemge/MessageBrokerClient"
)

func (app *App) GetOnDisconnectHandler() MessageBrokerClient.OnDisconnectHandler {
	return func(connection *MessageBrokerClient.WebsocketConnection) {
		app.logger.Log("Connection closed")
	}
}
