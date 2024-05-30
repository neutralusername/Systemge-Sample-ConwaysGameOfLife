package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/Websocket"
)

func (app *App) GetOnConnectHandler() MessageBrokerClient.OnConnectHandler {
	return func(connection *Websocket.Connection) {
		response, err := app.messageBrokerClient.SyncMessage(Message.NewSync("getGridSync", app.messageBrokerClient.GetName(), connection.Id))
		if err != nil {
			app.logger.Log(Error.New(err.Error()).Error())
			return
		}
		connection.Send([]byte(response.Serialize()))
	}
}
