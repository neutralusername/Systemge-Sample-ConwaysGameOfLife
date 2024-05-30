package appWebsocket

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Websocket"
)

func (app *App) GetOnDisconnectHandler() MessageBrokerClient.OnDisconnectHandler {
	return func(connection *Websocket.Connection) {
		app.logger.Log("Connection closed")
	}
}
